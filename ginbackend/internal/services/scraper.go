package services

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

const (
	indeedUrlPrefix   string = "https://www.indeed.com"
	seekUrlPrefix     string = "https://www.seek.com.au"
	seekSweJobsSuffix string = "/software-engineer-jobs"
)

type ScrapePayload struct {
	ShouldAssess bool      `json:"shouldAssess"`
	UserId       string    `json:"userId"`
	Url          string    `json:"url"`
	Candidate    Candidate `json:"candidate"`
}

type ScrapeIndeedPayload struct {
	Url string `json:"url"`
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
	mu              sync.Mutex
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

func (s *ScraperService) ScrapeSeekJobPage(payload *ScrapePayload) []*ScrapedJob {
	var jobs []*ScrapedJob

	// Create a collector for the job listings page
	jobListCollector := s.getCollector("jobList")

	// Create a collector for the individual job pages
	jobDetailCollector := s.getCollector("jobDetail")

	// Scrape the main page containing the list of jobs
	jobListCollector.OnHTML("body", func(body *colly.HTMLElement) {
		// Iterate through all job cards
		body.ForEach("[data-automation='job-list-view-job-link']", func(i int, e *colly.HTMLElement) {
			// Fetch link to job page and navigate to page
			jobLinkSuffix := e.Attr("href")
			jobDetailCollector.Visit(fmt.Sprintf("%s%s", seekUrlPrefix, jobLinkSuffix))
		})
	})

	// For each individual job, deconstruct the page
	jobDetailCollector.OnHTML("body", func(body *colly.HTMLElement) {
		// Construct the job
		scrapedJob := &ScrapedJob{
			Title:    body.DOM.Find("[data-automation='job-detail-title']").Text(),
			Company:  body.DOM.Find("[data-automation='advertiser-name']").Text(),
			Location: body.DOM.Find("[data-automation='job-detail-location']").Text(),
		}

		var parts []string
		callback := func(_ int, p *goquery.Selection) {
			text := p.Text()
			parts = append(parts, text)
		}

		// Construct job description
		jobDescriptionDiv := body.DOM.Find("[data-automation='jobAdDetails']")
		jobDescriptionDiv.Find("p").Each(callback)
		jobDescriptionDiv.Find("div").Each(callback)

		scrapedJob.Description = strings.Join(parts, "\n")
		jobs = append(jobs, scrapedJob)

		if !payload.ShouldAssess {
			return
		}

		go func() {
			job := s.jobService.CompleteScrapedJob(scrapedJob)

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

	jobListCollector.Visit(payload.Url)
	return jobs
}

func (s *ScraperService) ScrapeIndeedJobPage(payload *ScrapeIndeedPayload) []*ScrapedJob {
	var jobs []*ScrapedJob

	// Create a collector for the job listings page
	jobListCollector := s.getCollector("jobList")

	// Create a collector for the individual job pages
	jobDetailCollector := s.getCollector("jobDetail")

	jobListCollector.OnHTML(".mosaic-jobResults", func(div *colly.HTMLElement) {
		// Iterate through all job cards
		div.ForEach("a", func(_ int, a *colly.HTMLElement) {
			jobLinkSuffix := a.Attr("href")
			// Only navigate to the link if it is a link to a job
			if strings.Contains("/pagead/", jobLinkSuffix) {
				jobDetailCollector.Visit(fmt.Sprintf("%s%s", indeedUrlPrefix, jobLinkSuffix))
			}
		})
	})

	jobDetailCollector.OnHTML("body", func(body *colly.HTMLElement) {
		log.Infoln("Scraping individual Indeed job posting")
		// Construct the job
		scrapedJob := &ScrapedJob{
			Title:    body.DOM.Find("[data-testid='jobsearch-JobInfoHeader-title']").Find("span").Text(),
			Company:  body.DOM.Find("[data-testid='inlineHeader-companyName']").Find("span").Find("a").Text(),
			Location: body.DOM.Find("[data-testid='jobsearch-JobInfoHeader-companyLocation']").Find("span").Text(),
		}

		var parts []string
		callback := func(_ int, p *goquery.Selection) {
			text := p.Text()
			parts = append(parts, text)
		}

		// Construct job description
		jobDescriptionDiv := body.DOM.Find("[id='jobDescriptionText']")
		jobDescriptionDiv.Find("p").Each(callback)
		jobDescriptionDiv.Find("ul").Find("li").Each(callback)

		scrapedJob.Description = strings.Join(parts, "\n")
		jobs = append(jobs, scrapedJob)
	})

	jobListCollector.Visit(payload.Url)
	return jobs
}

func (s *ScraperService) GetAssessments(userId string) []*ScrapedJobAssessment {
	params := make(map[string]string)
	params["id"] = UserIdToAssessmentId(userId)
	existingAssessments, err := s.assessorService.QueryAssessments(params)

	if err != nil {
		log.Errorf("Failed to fetch assessments: %s\n", err.Error())
	}

	jobAssessments := make([]*ScrapedJobAssessment, 0)
	for _, a := range existingAssessments {
		params = make(map[string]string)
		params["id"] = a.JobId
		jobs, err := s.jobService.QueryJobs(params)

		if err != nil || len(jobs) == 0 {
			continue
		}

		jobAssessments = append(jobAssessments, &ScrapedJobAssessment{
			Assessment: *a,
			Job:        *jobs[0],
		})
	}
	return jobAssessments
}

func (s *ScraperService) registerCallbacks(client *colly.Collector) {
	headers := map[string]string{}
	client.OnRequest(func(r *colly.Request) {
		for k, v := range headers {
			r.Headers.Set(k, v)
		}
		log.Infof("Visiting page with URL %s\n", r.URL.String())
	})
	client.OnResponse(func(r *colly.Response) {
		log.Infof("Successfully visited page with URL %s\n", r.Request.URL.String())
	})
	client.OnScraped(func(r *colly.Response) {
		log.Infof("Successfully scraped page with URL %s\n", r.Request.URL.String())
	})
	client.OnError(func(r *colly.Response, err error) {
		log.Errorf("Request with URL %s failed with Response: \n%v \nError: \n%s", r.Request.URL.String(), r.Headers, err.Error())
	})
}

func (s *ScraperService) getCollector(name string) *colly.Collector {
	c, ok := s.clientMap[name]
	if ok {
		return c
	}
	newCollector := colly.NewCollector()
	s.registerCallbacks(newCollector)

	s.mu.Lock()
	s.clientMap[name] = newCollector
	s.mu.Unlock()

	return newCollector
}
