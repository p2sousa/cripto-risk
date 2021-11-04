package main

import (
	"time"

	"github.com/p2sousa/cripto-risk/src/core/usecase"
	"github.com/p2sousa/cripto-risk/src/infra/client"
	"github.com/p2sousa/cripto-risk/src/infra/repository"
)

func main() {
	client := client.NewMercadoBitcoin()
	repo := repository.NewCoinSummaryRepository()
	usecase := usecase.NewImportDaySummary(client, repo)

	t, _ := time.Parse("2006-01-02", "2021-10-01")
	usecase.Execute("BTC", t)
}
