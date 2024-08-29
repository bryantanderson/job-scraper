package services

import (
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type CandidateDto struct {
	Education   []Education  `json:"education"`
	Experiences []Experience `json:"experiences"`
	Skills		[]string 	 `json:"skills"`
	Summary     string       `json:"summary"`
	Location    string       `json:"location"`
}

type Candidate struct {
	Id          string       `json:"id" bson:"_id"`
	Education   []Education  `json:"education"`
	Experiences []Experience `json:"experiences"`
	Skills      []Skill      `json:"skills"`
	Summary     string       `json:"summary"`
	Location    string       `json:"location"`
}

type Skill struct {
	Description string `json:"description"`
}

type Education struct {
	Title       string `json:"title"`
	Institute   string `json:"institute"`
	Description string `json:"description"`
}

type Experience struct {
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	IsCurrent   bool       `json:"isCurrent"`
}

type CandidateStore interface {
	Create(c *Candidate) error
	Get(id string) (*Candidate, error)
	Delete(id string) error
}

type CandidateService struct {
	store CandidateStore
}

func NewCandidateService(store CandidateStore) *CandidateService {
	return &CandidateService{
		store: store,
	}
}

func (s *CandidateService) CreateCandidate(dto *CandidateDto) (*Candidate, error) {
	skills := make([]Skill, 0, len(dto.Skills))

	for _, s := range dto.Skills {
		skill := Skill{
			Description: s,
		}
		skills = append(skills, skill)
	}

	candidate := Candidate{
		Id: uuid.NewString(),
		Education: dto.Education,
		Experiences: dto.Experiences,
		Skills: skills,
		Summary: dto.Summary,
		Location: dto.Location,
	}
	err := s.store.Create(&candidate)

	if err != nil {
		log.Errorf("Failed to create candidate: %s\n", err.Error())
		return nil, err
	}

	return &candidate, nil
}

func (s *CandidateService) GetCandidate(id string) (*Candidate, error) {
	candidate, err := s.store.Get(id)

	if err != nil {
		log.Errorf("Failed to get candidate: %s\n", err.Error())
		return nil, err
	}

	return candidate, nil
}

func (s *CandidateService) DeleteCandidate(id string) error {
	return s.store.Delete(id)
}