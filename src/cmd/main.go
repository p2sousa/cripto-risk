package main

import (
	"time"

	"github.com/p2sousa/cripto-risk/src/core/usecase"
	"github.com/p2sousa/cripto-risk/src/infra/client"
	repository "github.com/p2sousa/cripto-risk/src/infra/repository/textfile"
)

func main() {
	importSummary()
	calcRisk()
}

func importSummary() {
	client := client.NewMercadoBitcoin()
	repo := repository.NewCoinSummaryRepository()
	usecase := usecase.NewImportDaySummary(client, repo)

	t, _ := time.Parse("2006-01-02", "2017-01-01")
	usecase.Execute("BTC", t)
}

func calcRisk() {
	sumRepo := repository.NewCoinSummaryRepository()
	riskRepo := repository.NewRiskRepository()

	usecase := usecase.NewCalculateRisk(riskRepo, sumRepo)
	usecase.Execute()
}
