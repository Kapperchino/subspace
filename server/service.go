package server

import (
	"database/sql"
	"github.com/Kapperchino/subspace/util"
)

type Service interface {
	getDB() *sql.DB
	getValidation() *util.Validation
	getConfig() *util.Config
}
