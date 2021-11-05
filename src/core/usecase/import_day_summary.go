package usecase

import (
	"fmt"
	"log"
	"time"

	"github.com/p2sousa/cripto-risk/src/core/entity"
	"github.com/p2sousa/cripto-risk/src/infra/client"
	"github.com/p2sousa/cripto-risk/src/infra/repository"
)

type IImportDaySummary interface {
	Execute(coin string, startDate time.Time) error
}

type ImportDaySummary struct {
	client     client.IClient
	repository repository.ICoinSummary
}

func NewImportDaySummary(client client.IClient, repository repository.ICoinSummary) *ImportDaySummary {
	return &ImportDaySummary{
		client:     client,
		repository: repository,
	}
}

func (ids *ImportDaySummary) Execute(coin string, startDate time.Time) error {

	dates := ids.rangeDates(startDate)

	for _, date := range dates {

		if ids.repository.Exists(date.Format("2006-01-02")) {
			continue
		}

		fmt.Printf("Fetch Day Summary: %s \n", date.Format("2006-01-02"))
		summary, err := ids.client.FetchDaySymmary(coin, date)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		coin := entity.Coin{
			AvgPrice: summary.AvgPrice,
			Date:     summary.Date,
		}

		if err := ids.repository.Save(coin); err != nil {
			log.Fatalln(err)
			return err
		}
	}

	return nil
}

func (ids *ImportDaySummary) rangeDates(start time.Time) []time.Time {
	now := time.Now().AddDate(0, 0, -1)

	list := make([]time.Time, 0)

	for start.Before(now) {
		list = append(list, start)
		start = start.Add(time.Duration(time.Hour * 24))
	}

	return list
}
