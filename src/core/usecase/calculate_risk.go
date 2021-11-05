package usecase

import (
	"log"
	"time"

	"github.com/p2sousa/cripto-risk/src/core/entity"
	"github.com/p2sousa/cripto-risk/src/infra/repository"
	"github.com/thoas/go-funk"
)

const (
	DAYS_MA      = 50
	DAYS_WEEK_MA = 350
	FORMAT_DATE  = "2006-01-02"
)

type ICalculateRisk interface {
	Execute() error
}

type summaryDay struct {
	date        string
	avg         float64
	simpleMedia float64
}

type calculateRisk struct {
	riskRepository    repository.IRisk
	summaryRepository repository.ICoinSummary
}

func NewCalculateRisk(riskRepo repository.IRisk, summaryRepo repository.ICoinSummary) *calculateRisk {
	return &calculateRisk{
		riskRepository:    riskRepo,
		summaryRepository: summaryRepo,
	}
}

func (calc *calculateRisk) Execute() error {

	mapSummarySimpleMedia, err := calc.generateMapSimpleMedia()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	for _, sum := range mapSummarySimpleMedia {
		if calc.riskRepository.Exists(sum.date) {
			continue
		}

		r := entity.Risk{
			Date:  sum.date,
			Point: sum.simpleMedia,
		}

		if err := calc.riskRepository.Save(r); err != nil {
			log.Fatalln(err)
			return err
		}
	}
	return nil
}

func (calc *calculateRisk) generateMapSimpleMedia() (map[string]summaryDay, error) {

	summary, err := calc.summaryRepository.FetchAll()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	mapSummaryDay := funk.Map(summary, func(date string, c entity.Coin) (string, summaryDay) {
		sm := calc.calculateRisk(date, summary)
		return date, summaryDay{
			avg:         c.AvgPrice,
			date:        c.Date,
			simpleMedia: sm,
		}
	}).(map[string]summaryDay)

	listSimpleMedia := funk.Map(mapSummaryDay, func(date string, sd summaryDay) float64 {
		return sd.simpleMedia
	}).([]float64)

	minSimpleMedia := funk.MinFloat64(listSimpleMedia)
	maxSimpleMedia := funk.MaxFloat64(listSimpleMedia)

	mapSummarySimpleMedia := funk.Map(mapSummaryDay, func(date string, sd summaryDay) (string, summaryDay) {
		sd.simpleMedia = calc.normalizeRisk(sd.simpleMedia, minSimpleMedia, maxSimpleMedia)
		return date, sd
	}).(map[string]summaryDay)

	return mapSummarySimpleMedia, nil
}

func (calc *calculateRisk) calculateRisk(day string, list map[string]entity.Coin) float64 {
	d, _ := time.Parse(FORMAT_DATE, day)

	smDay := calc.simpleMedia(d, DAYS_MA, list)
	smWeek := calc.simpleMedia(d, DAYS_WEEK_MA, list)

	return smDay / smWeek
}

func (calc *calculateRisk) simpleMedia(date time.Time, n int, list map[string]entity.Coin) float64 {

	values := make([]float64, 0)

	for i := 1; i < n; i++ {
		summary, ok := list[date.Format(FORMAT_DATE)]
		if !ok {
			continue
		}

		values = append(values, summary.AvgPrice)
		date = date.AddDate(0, 0, -1)
	}

	total := funk.SumFloat64(values)

	return total / float64(len(values))
}

func (calc *calculateRisk) normalizeRisk(sma float64, minSma float64, maxSma float64) float64 {
	return 0 + (sma-minSma)*(1-0)/(maxSma-minSma)
}
