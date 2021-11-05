package repository

import "github.com/p2sousa/cripto-risk/src/core/entity"

type ICoinSummary interface {
	FetchAll() (map[string]entity.Coin, error)
	Exists(date string) bool
	Save(entity entity.Coin) error
}

type IRisk interface {
	FetchAll() (map[string]entity.Risk, error)
	Exists(date string) bool
	Save(entity entity.Risk) error
}
