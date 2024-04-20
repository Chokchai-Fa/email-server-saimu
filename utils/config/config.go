package config

import (
	"database/sql"
	"emailserver-saimu/utils/dbpostgres"
	"emailserver-saimu/utils/email"
)

type AppConfig struct {
	DBPG        *dbpostgres.DBPG
	DBPGCli     *sql.DB
	EmailServer email.EmailServer
	// SQSURL
}
