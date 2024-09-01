package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sincidium/linkd/api/setup"
	"sync"
	"time"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	log "github.com/sirupsen/logrus"
)

const pointExplanation = `YOU ARE THE WORLD'S BEST EXPERT IN EVIDENCE-BASED ANALYSIS, KNOWN FOR PROVIDING PRECISE, WELL-SUPPORTED REASONING.
### INSTRUCTIONS ###
- PROVIDE clear reasoning for why the point is valid.
- SUPPORT your reasoning with specific details from the provided information.
- INCLUDE the rationale behind the answers or procedures.

### Chain of Thoughts ###

1. **State the Point:**
   - Clearly state your point and its relevance.

2. **Explain Reasoning:**
   - Explain why the point is applicable and fits the context.
   - Provide detailed reasoning behind the answers or procedures.

3. **Cite Evidence:**
   - Provide specific details from the provided information.

4. **Summarize:**
   - Summarize how the evidence supports your point and reaffirm its validity.

### What Not To Do ###

- NEVER PROVIDE REASONING WITHOUT EVIDENCE.
- NEVER USE VAGUE OR AMBIGUOUS STATEMENTS.
- NEVER OMIT REFERENCES TO DETAILS OR RATIONALES.
- NEVER FABRICATE INFORMATION.`

type AssessPayload struct {
	Job       Job       `json:"job"`
	UserId    string    `json:"userId"`
	Candidate Candidate `json:"candidate"`
}

type CandidateAssessment struct {
	mu         sync.Mutex
	Assessment Assessment
	client     *instructor.InstructorOpenAI
}

type Assessment struct {
	Id                  string `json:"id" bson:"_id"`
	JobId               string `json:"jobId"`
	Score               uint8  `json:"score"`
	ExperiencePoint     bool   `json:"experiencePoint"`
	LocationMatch       match  `json:"locationMatch"`
	ResponsibilityScore score  `json:"responsibilityScore"`
	SkillsScore         score  `json:"skillsScore"`
	RequirementPoint    point  `json:"requirementPoint"`
	CompatibilityPoint  point  `json:"compatibilityPoint"`
	CreatedAt           string `json:"createdAt"`
	ElasticId           string `json:"elasticId,omitempty" `
}

type Rubric struct {
	Id     string   `json:"id" bson:"_id"`
	Points []string `json:"points" jsonschema:"description=A list of the responsibilities that the job descriptions desires within the ideal candidate"`
}

type rubricInstruct struct {
	Points []string `json:"points" jsonschema:"description=A list of the responsibilities that the job descriptions desires within the ideal candidate"`
}

type point struct {
	Explanation string `json:"explanation" jsonschema:"description=Provide a detailed explanation as per the given instructions."`
	IsValid     bool   `json:"isValid" jsonschema:"description=The point you justified in the explanation. Remember this is the feedback that is provided to the HR Manager so have more information in this section."`
}

type score struct {
	Explanation string `json:"explanation" jsonschema:"description=Assess the user's past job experiences and their skills against the criteria for the job. Carefully evaluate and justify how the candidate meets/does not meet every single criteria."`
	Score       uint8  `json:"score" jsonschema:"description=Count of the number of criterions the candidate met successfully"`
}

type match struct {
	IsMatch bool `json:"isMatch" jsonschema:"description=Whether or not the candidate matches the desired location."`
}

type AssessorStore interface {
	Create(a *Assessment) error
	CreateInternalJobCriteria(jc *Rubric) error
	QueryInternalJobCriteria(id string) (*Rubric, error)
	FindById(userId string) (*Assessment, error)
	Query(params map[string]string) ([]*Assessment, error)
	Delete(userId string) error
}

type AssessorService struct {
	inTopic             string
	inTopicSubscription string
	outTopic            string
	store               AssessorStore
	client              *instructor.InstructorOpenAI
	eventService        *EventService
}

func InitializeAssessorService(
	s *setup.ApplicationSettings,
	c *instructor.InstructorOpenAI,
	e *EventService,
	as AssessorStore,
) *AssessorService {
	service := &AssessorService{
		inTopic:             s.AssessmentTasksTopic,
		inTopicSubscription: topicNameToSubscriptionName(s.AssessmentTasksTopic),
		outTopic:            s.AssessmentResultsTopic,
		client:              c,
		store:               as,
		eventService:        e,
	}
	service.registerSubscribers()
	return service
}

func (a *AssessorService) registerSubscribers() {
	routine := func() {
		numWorkers := 1
		mChan := make(chan []byte, numWorkers)
		defer close(mChan)
		for i := 1; i <= numWorkers; i++ {
			go a.worker(i, mChan)
		}
		a.eventService.Subscribe(a.inTopic, a.inTopicSubscription, mChan)
	}
	a.eventService.Register(routine)
}

func (a *AssessorService) worker(i int, mChan <-chan []byte) {
	log.Infof("Worker number %d for assessment tasks starting...\n", i)
	for {
		messageBody, ok := <-mChan

		if !ok {
			// If the channel is closed, the worker should stop
			break
		}

		payload := &AssessPayload{}
		err := json.Unmarshal(messageBody, payload)

		if err != nil {
			log.Printf("Unable to unmarshal message body: %s \n", err)
			continue
		}

		a.AssessCandidate(payload)
	}
}

func (a *AssessorService) AssessCandidate(payload *AssessPayload) {
	ca := CandidateAssessment{
		Assessment: Assessment{
			Id:    a.UserIdToAssessmentId(payload.UserId),
			JobId: payload.Job.Id,
		},
		client: a.client,
	}
	jobCriteria, err := a.createCriteria(&payload.Job)
	a.store.CreateInternalJobCriteria(jobCriteria)

	if err != nil {
		log.Errorf("Failed to create criteria for job: %s \n", err.Error())
		return
	}

	// Instantiate wait group, error channel, and cancel context
	// The cancel func is run if any errors occur to halt all operations
	var wg sync.WaitGroup

	wgDone := make(chan bool)
	errChan := make(chan error, 1)
	defer close(errChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Infoln("Tasks for assessing candidate beginning")

	tasks := []func(){
		func() { ca.assessExperience(payload, &wg) },
		func() { ca.assessLocation(ctx, payload, &wg, errChan) },
		func() { ca.assessRequirements(ctx, payload, &wg, errChan) },
		func() { ca.assessCompatibility(ctx, payload, &wg, errChan) },
		func() { ca.assessResponsibilities(ctx, &payload.Candidate, jobCriteria, &wg, errChan) },
		func() { ca.assessSkills(ctx, &payload.Candidate, jobCriteria, &wg, errChan) },
	}
	for _, t := range tasks {
		wg.Add(1)
		go t()
	}
	go func() {
		wg.Wait()
		close(wgDone)
	}()
	/*
	 * The select statement will cancel all goroutine operations if an error
	 * is received, a cancel is called, or if the timeout is reached.
	 * It will otherwise proceed to wait for all goroutines to finish.
	 */
	select {
	case err := <-errChan:
		log.Errorf("Received error: %v, cancelling all goroutines\n", err)
		cancel()

	case <-ctx.Done():
		log.Error(ctx.Err().Error())

	case <-time.After(time.Minute):
		log.Errorln("Timeout reached, cancelling all goroutines")
		cancel()

	case <-wgDone:
		log.Infoln("All goroutines have completed successfully, finalizing assessment score")
		ca.finalizeAssessment()
		err = a.store.Create(&ca.Assessment)

		if err != nil {
			log.Infof("Failed to save assessment: %s\n", err.Error())
		}
	}
}

func (a *AssessorService) GetAssessment(userId string) (*Assessment, error) {
	assessment, err := a.store.FindById(a.UserIdToAssessmentId(userId))
	return assessment, err
}

func (a *AssessorService) QueryAssessments(params map[string]string) ([]*Assessment, error) {
	return a.store.Query(params)
}

func (a *AssessorService) createCriteria(job *Job) (*Rubric, error) {
	// Check if a criteria for the job already exists
	existingRubric, err := a.store.QueryInternalJobCriteria(a.jobIdToRubricId(job.Id))

	if err == nil && existingRubric != nil {
		return existingRubric, nil
	}

	responsibilities, _ := convertToJson(job.Responsibilities)
	prompt := fmt.Sprintf(`
	Given the job description, use the job responsibilities as a criteria to craft a criteria that a candidate can be assessed against.
    The resulting rubric should be as close to the responsibilities mentioned within the job description as possible. Use the same wording.

    Job Responsibilities:
	%s
	`,
		responsibilities)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var rubricInstruct rubricInstruct
	resp, err := a.client.CreateChatCompletion(
		ctx,
		makeChatCompletionRequest(prompt, 500),
		&rubricInstruct,
	)
	_ = resp
	rubric := Rubric{
		Id:     a.jobIdToRubricId(job.Id),
		Points: rubricInstruct.Points,
	}
	return &rubric, err
}

func (a *AssessorService) UserIdToAssessmentId(userId string) string {
	return fmt.Sprintf("%s_assessment", userId)
}

func (a *AssessorService) jobIdToRubricId(jobId string) string {
	return fmt.Sprintf("%s_criteria", jobId)
}

/*

Candidate Assessment receivers from this point onward.

*/

func (ca *CandidateAssessment) assessRequirements(
	ctx context.Context, payload *AssessPayload, wg *sync.WaitGroup, errChan chan<- error,
) {
	log.Infoln("Beginning assessRequirements")
	defer wg.Done()
	educations, err := convertToJson(payload.Candidate.Education)

	if err != nil {
		log.Errorln("Error converting educations to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	qualifications, err := convertToJson(payload.Job.Qualifications)

	if err != nil {
		log.Errorln("Error converting qualifications to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	prompt := fmt.Sprintf(`
	Does the candidate fulfill the required qualifications for the job?

	candidate Education History:
	%s

	Job Education Requirements:
	%s

	Explanation you will provide:
	%s
	`,
		educations, qualifications, pointExplanation)

	var p point

	select {
	case <-ctx.Done():
		log.Errorln("Context canceled, aborting assessRequirements")
		return
	default:
		resp, err := ca.client.CreateChatCompletion(
			ctx,
			makeChatCompletionRequest(prompt, 300),
			&p,
		)
		_ = resp // sends back original response so no information loss from original API

		if err != nil {
			log.Errorf("assessRequirements OpenAI call has failed with error %v\n", err)
			handleGoroutineError(err, errChan)
			return
		}
	}

	ca.mu.Lock()
	ca.Assessment.RequirementPoint = p
	ca.mu.Unlock()
	log.Infoln("Ending assessRequirements")
}

func (ca *CandidateAssessment) assessCompatibility(
	ctx context.Context, payload *AssessPayload, wg *sync.WaitGroup, errChan chan<- error,
) {
	log.Infoln("Beginning assessCompatibility")
	defer wg.Done()
	description, err := convertToJson(payload.Job.Description)

	if err != nil {
		log.Errorln("Error converting job description to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	prompt := fmt.Sprintf(`
	Is the candidate compatible with the job?

	candidate Summary:
	%s

	Job Description:
	%s
	`,
		payload.Candidate.Summary, description)

	var p point

	select {
	case <-ctx.Done():
		log.Errorln("Context canceled, aborting assessCompatibility")
		return
	default:
		resp, err := ca.client.CreateChatCompletion(
			ctx,
			makeChatCompletionRequest(prompt, 500),
			&p,
		)
		_ = resp // sends back original response so no information loss from original API

		if err != nil {
			log.Errorf("assessCompatibility OpenAI call has failed with error %v\n", err)
			handleGoroutineError(err, errChan)
			return
		}
	}

	ca.mu.Lock()
	ca.Assessment.CompatibilityPoint = p
	ca.mu.Unlock()
	log.Infoln("Ending assessCompatibility")
}

func (ca *CandidateAssessment) assessLocation(
	ctx context.Context, payload *AssessPayload, wg *sync.WaitGroup, errChan chan<- error,
) {
	log.Infoln("Beginning assessLocation")
	defer wg.Done()
	if payload.Job.Location == "" {
		ca.assignLocationMatch(false)
		return
	} else if payload.Job.Location == payload.Candidate.Location {
		ca.assignLocationMatch(true)
		return
	}

	prompt := fmt.Sprintf(`
	Does the candidates location match with the desired location of the job description?

    candidate location:
	%s

	Job Description location:
	%s
	`,
		payload.Candidate.Location, payload.Job.Description)

	var m match

	select {
	case <-ctx.Done():
		log.Errorln("Context canceled, aborting assessLocation")
		return
	default:
		resp, err := ca.client.CreateChatCompletion(
			ctx,
			makeChatCompletionRequest(prompt, 10),
			&m,
		)
		_ = resp

		if err != nil {
			log.Errorf("assessLocation OpenAI call has failed with error %v\n", err)
			handleGoroutineError(err, errChan)
			return
		}
	}

	ca.mu.Lock()
	ca.Assessment.LocationMatch = m
	ca.mu.Unlock()
	log.Infoln("Ending assessLocation")
}

func (ca *CandidateAssessment) assignLocationMatch(isMatch bool) {
	m := match{
		IsMatch: isMatch,
	}
	ca.mu.Lock()
	ca.Assessment.LocationMatch = m
	ca.mu.Unlock()
}

func (ca *CandidateAssessment) assessResponsibilities(
	ctx context.Context, candidate *Candidate, rubric *Rubric, wg *sync.WaitGroup, errChan chan<- error,
) {
	log.Infoln("Beginning assessResponsibilities")
	defer wg.Done()
	rubricJson, err := convertToJson(rubric)

	if err != nil {
		log.Errorln("Error converting job rubric to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	experience, err := convertToJson(candidate.Experiences)

	if err != nil {
		log.Errorln("Error converting candidate experience to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	prompt := fmt.Sprintf(`
	Job rubric:
	%s

	candidate Past Experiences (provided in json):
	%s
	`,
		rubricJson, experience)

	var s score

	select {
	case <-ctx.Done():
		log.Errorln("Context canceled, aborting assessResponsibilities")
		return
	default:
		resp, err := ca.client.CreateChatCompletion(
			ctx,
			makeChatCompletionRequest(prompt, 500),
			&s,
		)
		_ = resp

		if err != nil {
			log.Errorf("assessResponsibilities OpenAI call has failed with error %v\n", err)
			handleGoroutineError(err, errChan)
			return
		}
	}

	ca.mu.Lock()
	ca.Assessment.ResponsibilityScore = s
	ca.mu.Unlock()
	log.Infoln("Ending assessResponsibilities")
}

func (ca *CandidateAssessment) assessSkills(
	ctx context.Context, candidate *Candidate, rubric *Rubric, wg *sync.WaitGroup, errChan chan<- error,
) {
	log.Infoln("Beginning assessSkills")
	defer wg.Done()
	rubricJson, err := convertToJson(rubric)

	if err != nil {
		log.Errorln("Error converting rubric to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	skillsJson, err := convertToJson(candidate.Skills)

	if err != nil {
		log.Errorln("Error converting candidate skills to JSON:", err)
		handleGoroutineError(err, errChan)
		return
	}

	prompt := fmt.Sprintf(`
	Job rubric:
	%s

	candidate skills (provided in json):
	%s
	`, rubricJson, skillsJson)

	var s score

	select {
	case <-ctx.Done():
		log.Errorln("Context canceled, aborting assessSkills")
		return
	default:
		resp, err := ca.client.CreateChatCompletion(
			ctx,
			makeChatCompletionRequest(prompt, 500),
			&s,
		)
		_ = resp

		if err != nil {
			log.Errorf("assessSkills OpenAI call has failed with error %v\n", err)
			handleGoroutineError(err, errChan)
			return
		}
	}

	ca.mu.Lock()
	ca.Assessment.SkillsScore = s
	ca.mu.Unlock()
}

func (ca *CandidateAssessment) assessExperience(
	payload *AssessPayload, wg *sync.WaitGroup,
) {
	log.Infoln("Beginning assessExperience")
	defer wg.Done()
	candidate_yoe := 0
	for _, experience := range payload.Candidate.Experiences {
		currentDate := time.Now()
		endDate := experience.EndDate
		if experience.IsCurrent || experience.EndDate == nil {
			endDate = &currentDate
		}
		duration := (*endDate).Sub(experience.StartDate)
		yoe := int(math.Round(duration.Hours() / (24 * 365)))
		candidate_yoe += yoe
	}
	yoeDiff := candidate_yoe - int(payload.Job.YearsOfExperience)

	ca.mu.Lock()
	ca.Assessment.ExperiencePoint = yoeDiff >= 0
	ca.mu.Unlock()
	log.Infoln("Ending assessExperience")
}

//

func (ca *CandidateAssessment) finalizeAssessment() {
	log.Infoln("Finalizing assessment score")
	var score uint8
	score = 0
	a := &ca.Assessment
	if a.RequirementPoint.IsValid && a.RequirementPoint.Explanation != "" {
		score += 1
	}
	if a.CompatibilityPoint.IsValid && a.CompatibilityPoint.Explanation != "" {
		score += 1
	}
	if a.LocationMatch.IsMatch {
		score += 1
	}
	if a.ExperiencePoint {
		score += 1
	}
	if a.ResponsibilityScore.Explanation != "" {
		score += a.ResponsibilityScore.Score
	}
	a.Score = score
	a.CreatedAt = time.Now().String()
}
