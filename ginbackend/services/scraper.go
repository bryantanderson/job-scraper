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

type ScrapedJob struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type ScraperService struct {
	clientMap  map[string]*colly.Collector
	jobService *JobService
}

func NewScraperService(jobService *JobService) *ScraperService {
	s := ScraperService{
		clientMap:  make(map[string]*colly.Collector),
		jobService: jobService,
	}
	return &s
}

func (s *ScraperService) ScrapeJobs() []*ScrapedJob {
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

		go s.jobService.CompleteScrapedJob(j)
	})

	jobListCollector.Visit(fmt.Sprintf("%s%s", seekUrlPrefix, seekSweJobsSuffix))
	return jobs
}

func (s *ScraperService) registerCallbacks(client *colly.Collector) {
	client.OnRequest(func(r *colly.Request) {
		log.Infof("Visiting page with URL %s", r.URL.String())
	})
	client.OnResponse(func(r *colly.Response) {
		log.Infof("Successfully visited page with URL %s", r.Request.URL.String())
	})
	client.OnScraped(func(r *colly.Response) {
		log.Infoln("Successfully scraped page with URL %s", r.Request.URL.String())
	})
	client.OnError(func(_ *colly.Response, err error) {
		log.Errorf("Error while scraping: %s", err.Error())
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
