package services

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	seekUrlPrefix     string = "https://www.seek.com.au"
	seekSweJobsSuffix string = "/software-engineer-jobs"
)

type ScrapeSeekPayload struct {
	UserId    string    `json:"userId"`
	Candidate Candidate `json:"candidate"`
}

type ScrapedJobAssessment struct {
	Job        Job        `json:"job"`
	Assessment Assessment `json:"assessment"`
}

type ScrapedJob struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type ScraperService struct {
	clientMap       map[string]*colly.Collector
	jobService      *JobService
	assessorService *AssessorService
}

func InitializeScraperService(jobService *JobService, assessorService *AssessorService) *ScraperService {
	s := ScraperService{
		clientMap:       make(map[string]*colly.Collector),
		jobService:      jobService,
		assessorService: assessorService,
	}
	return &s
}

func (s *ScraperService) ScrapeSeek(payload *ScrapeSeekPayload) []*ScrapedJob {
	var jobs []*ScrapedJob

	// Create a collector for the job listings page
	jobListCollector := s.getCollector("jobList")

	// Create a collector for the individual job pages
	jobDetailCollector := s.getCollector("jobDetail")

	jobListCollector.OnHTML("body", func(body *colly.HTMLElement) {
		// Iterate through all job cards
		body.ForEach("[data-automation='job-list-view-job-link']", func(i int, e *colly.HTMLElement) {
			// Fetch link to job page and navigate to page
			jobLinkSuffix := e.Attr("href")
			jobLink := fmt.Sprintf("%s%s", seekUrlPrefix, jobLinkSuffix)
			jobDetailCollector.Visit(jobLink)
		})
	})

	jobDetailCollector.OnHTML("body", func(body *colly.HTMLElement) {
		// Construct the job
		j := &ScrapedJob{}
		j.Title = body.DOM.Find("[data-automation='job-detail-title']").Text()
		j.Company = body.DOM.Find("[data-automation='advertiser-name']").Text()
		j.Location = body.DOM.Find("[data-automation='job-detail-location']").Text()

		var parts []string
		callback := func(i int, p *goquery.Selection) {
			text := p.Text()
			parts = append(parts, text)
		}

		jobDescriptionDiv := body.DOM.Find("[data-automation='jobAdDetails']")
		jobDescriptionDiv.Find("p").Each(callback)
		jobDescriptionDiv.Find("div").Each(callback)

		j.Description = strings.Join(parts, "\n")
		jobs = append(jobs, j)

		go func() {
			job := s.jobService.CompleteScrapedJob(j)

			if job == nil {
				log.Errorln("AI job completion was not completed successfully, returning early")
				return
			}

			payload := AssessPayload{
				Job:       *job,
				UserId:    payload.UserId,
				Candidate: payload.Candidate,
			}
			s.assessorService.AssessCandidate(&payload)
		}()
	})

	jobListCollector.Visit(fmt.Sprintf("%s%s", seekUrlPrefix, seekSweJobsSuffix))
	return jobs
}

func (s *ScraperService) GetScrapedSeekAssessments(userId string) []*ScrapedJobAssessment {
	sja := make([]*ScrapedJobAssessment, 0)
	queryParams := make(map[string]string)
	queryParams["id"] = s.assessorService.UserIdToAssessmentId(userId)

	assessments, err := s.assessorService.QueryAssessments(queryParams)

	if err != nil {
		log.Errorf("Failed to fetch assessments: %s\n", err.Error())
	}

	for _, a := range assessments {
		sj := ScrapedJobAssessment{}

		jobQueryParams := make(map[string]string)
		jobQueryParams["id"] = a.JobId

		jobs, err := s.jobService.QueryJobs(jobQueryParams)

		if err != nil || len(jobs) == 0 {
			continue
		}

		sj.Assessment = *a
		sj.Job = *jobs[0]

		sja = append(sja, &sj)
	}
	return sja
}

func (s *ScraperService) registerCallbacks(client *colly.Collector) {
	client.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting page with URL %s\n", r.URL.String())
	})
	client.OnResponse(func(r *colly.Response) {
		log.Infof("Successfully visited page with URL %s\n", r.Request.URL.String())
	})
	client.OnScraped(func(r *colly.Response) {
		log.Infof("Successfully scraped page with URL %s\n", r.Request.URL.String())
	})
	client.OnError(func(_ *colly.Response, err error) {
		log.Errorf("Error while scraping: %s\n", err.Error())
	})
}

func (s *ScraperService) getCollector(name string) *colly.Collector {
	c, ok := s.clientMap[name]
	if ok {
		return c
	}
	newCollector := colly.NewCollector()
	s.registerCallbacks(newCollector)
	s.clientMap[name] = newCollector
	return newCollector
}
