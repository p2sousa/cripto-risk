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
	DB_RISK = "db_risk.csv"
)

type riskRepository struct {
	db string
}

func NewRiskRepository() *riskRepository {
	repository := riskRepository{
		db: DB_RISK,
	}

	repository.createDatabase()
	return &repository
}

func (csr *riskRepository) createDatabase() {
	if _, err := os.Stat(DB_RISK); errors.Is(err, os.ErrNotExist) {
		colunm := [][]string{
			{"coin", "date", "point"},
		}

		csvFile, err := os.Create(DB_RISK)
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
		defer csvFile.Close()

		csvwriter := csv.NewWriter(csvFile)
		csvwriter.WriteAll(colunm)
	}
}

func (risk *riskRepository) FetchAll() (map[string]entity.Risk, error) {
	f, err := os.Open(risk.db)
	if err != nil {
		log.Fatal("Unable to read input file "+risk.db, err)
	}
	defer f.Close()

	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+risk.db, err)
	}

	list := make(map[string]entity.Risk)
	for _, record := range records {
		point, _ := strconv.ParseFloat(record[2], 32)
		list[record[1]] = entity.Risk{
			Date:  record[1],
			Point: point,
		}
	}

	return list, nil
}

func (risk *riskRepository) Exists(date string) bool {
	data, err := risk.FetchAll()
	if err != nil {
		log.Fatalf("failed fetch data: %s", err)
		return false
	}

	if _, ok := data[date]; ok {
		return true
	}

	return false
}

func (risk *riskRepository) Save(entity entity.Risk) error {

	records := [][]string{
		{"BTC", entity.Date, fmt.Sprintf("%.1f", entity.Point)},
	}

	csvFile, err := os.OpenFile(risk.db, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed open file: %s", err)
		return err
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.WriteAll(records)
	return nil
}
