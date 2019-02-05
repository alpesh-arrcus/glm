package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/toravir/glm/config"
	. "github.com/toravir/glm/context"
)

var dbDRIVER = "sqlite3" //because we have loaded 'go-sqlite3'
var ctxCache Context
var logger *zerolog.Logger

type dbState struct {
	db *sql.DB
}

func InitLicenseDb(ctx Context) Context {
	dbSt := dbState{}
	ctx.DbInfo = &dbSt

	ctxCache = ctx
	logger = config.GetLogger(ctx)
	dataSource := config.GetDBSourceName(ctx)
	db, err := sql.Open(dbDRIVER, dataSource)
	if err != nil {
		logger.Error().Str("DB", dataSource).AnErr("Error", err).Msg("Error opening dataSource")
	}
	//Enable non-blocking read/write
	db.Exec("PRAGMA journal_mode=WAL;")
	dbSt.db = db

	tblName := getCustomersTableName()
	sqlStmt := `
        create table if not exists ` + tblName + ` (customerName text not null primary key, customerSecret text not null, id integer not null, status text not null);
        `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("Creating Customer Table")
	}
	return ctx
}
