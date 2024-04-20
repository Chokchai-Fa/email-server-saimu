package dbpostgres

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func (d *DBPG) Connect(user, pass, server, database string, port int) (*sql.DB, error) {
	/* 	username:password@protocol(address)/dbname?param=value */
	//Connect pgsql
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		server, port, user, pass, database)

	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, server, port, database)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DBPG) ConnectWithURL(url string) (*sql.DB, error) {

	pqUrl, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, server, port, database)
	db, err := sql.Open("postgres", pqUrl)
	if err != nil {
		return nil, err
	}

	return db, nil

}

func (d *DBPG) ConnectWithOTEL(user, pass, server, database string, port int) (*sql.DB, error) {
	/* 	username:password@protocol(address)/dbname?param=value */
	//Connect pgsql
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		server, port, user, pass, database)

	// connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, server, port, database)
	db, err := otelsql.Open("postgres", connectionString, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	return db, nil
}
