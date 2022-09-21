package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DatabaseAdapter[M any] interface {
	QuerySingle(query string, args ...interface{}) (*M, error)
	QueryMany(query string, args ...interface{}) ([]*M, error)
	Execute(query string, args ...interface{}) (err error)
}

type DatabaseConnection struct {
	Server   string
	Username string
	Password string
	DbName   string

	conn *bun.DB
}

type databaseAdapter[M any] struct {
	conn *bun.DB
}

type userRow struct {
	Id       int32
	Username string
	Password string
	Nickname string
}

func ConnectDatabase(server string, username string, password string, dbName string) (*DatabaseConnection, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		username, password, server, dbName)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	dbConn := bun.NewDB(sqldb, pgdialect.New())

	return &DatabaseConnection{
		Server:   server,
		Username: username,
		Password: password,
		DbName:   dbName,

		conn: dbConn,
	}, nil
}

func newDatabaseAdapter[M any](conn *DatabaseConnection) DatabaseAdapter[M] {
	return &databaseAdapter[M]{
		conn: conn.conn,
	}
}

func (db *databaseAdapter[M]) QuerySingle(query string, args ...interface{}) (*M, error) {
	var container M
	ctx := context.Background()
	err := bun.NewRawQuery(db.conn, query, args...).Scan(ctx, &container)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &container, nil
}

func (db *databaseAdapter[M]) QueryMany(query string, args ...interface{}) ([]*M, error) {
	results := make([]*M, 0)
	ctx := context.Background()
	err := bun.NewRawQuery(db.conn, query, args...).Scan(ctx, &results)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return results, err
}

func (db *databaseAdapter[M]) Execute(query string, args ...interface{}) error {
	_, err := db.conn.Exec(query, args...)
	return err
}
