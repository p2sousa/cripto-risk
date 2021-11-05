package textfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/p2sousa/cripto-risk/src/core/entity"
)

const (
	DB_SUMMARY = "db_summary.csv"
)

type coinSummaryRepository struct {
	db string
}

func NewCoinSummaryRepository() *coinSummaryRepository {
	repository := coinSummaryRepository{
		db: DB_SUMMARY,
	}

	repository.createDatabase()

	return &repository
}

func (csr *coinSummaryRepository) createDatabase() {
	if _, err := os.Stat(DB_SUMMARY); errors.Is(err, os.ErrNotExist) {
		colunm := [][]string{
			{"coin", "date", "avg_price"},
		}

		csvFile, err := os.Create(DB_SUMMARY)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close()

		csvwriter := csv.NewWriter(csvFile)
		csvwriter.WriteAll(colunm)
	}
}

func (csr *coinSummaryRepository) FetchAll() (map[string]entity.Coin, error) {
	f, err := os.Open(DB_SUMMARY)
	if err != nil {
		log.Fatal("Unable to read input file "+DB_SUMMARY, err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+DB_SUMMARY, err)
	}

	list := make(map[string]entity.Coin)
	for _, record := range records {
		avg, _ := strconv.ParseFloat(record[2], 32)
		list[record[1]] = entity.Coin{
			Date:     record[1],
			AvgPrice: avg,
		}
	}

	return list, nil
}

func (csr *coinSummaryRepository) Exists(date string) bool {
	data, err := csr.FetchAll()
	if err != nil {
		log.Fatalf("failed fetch data: %s", err)
		return false
	}

	if _, ok := data[date]; ok {
		return true
	}

	return false
}

func (csr *coinSummaryRepository) Save(entity entity.Coin) error {

	records := [][]string{
		{"BTC", entity.Date, fmt.Sprintf("%f", entity.AvgPrice)},
	}

	csvFile, err := os.OpenFile(DB_SUMMARY, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed open file: %s", err)
		return err
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.WriteAll(records)
	return nil
}
