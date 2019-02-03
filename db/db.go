package db

import (
        . "github.com/toravir/glm/context"
        "github.com/toravir/glm/config"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var dbDRIVER = "sqlite3" //because we have loaded 'go-sqlite3'
var ctxCache Context

type dbState struct {
    db *sql.DB
}


func InitLicenseDb(ctx Context) Context {
        dbSt := dbState{}
        ctx.DbInfo = &dbSt

        ctxCache = ctx
        logger := config.GetLogger(ctx)
        dataSource := config.GetDBSourceName(ctx)
	db, err := sql.Open(dbDRIVER, dataSource)
        if err != nil {
            logger.Fatal().Str("DB", dataSource).AnErr("Error", err).Msg("Error opening dataSource")
        }
        dbSt.db = db

	sqlStmt := `
        create table if not exists customers (name text not null primary key, secret text not null, id integer not null);
        `
	_, err = db.Exec(sqlStmt)
        if err != nil {
            logger.Fatal().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("Creating Customer Table")
        }
        return ctx
}

func GetCustomerNames() []string {
        dbSt, err := ctxCache.DbInfo.(*dbState)
        logger := config.GetLogger(ctxCache)
        if !err {
            return nil
        }
        customers := []string{}
	rows, _ := dbSt.db.Query("select name from customers")
	defer rows.Close()
	for rows.Next() {
		var name string
                err := rows.Scan(&name)
                if err != nil {
                    logger.Fatal().AnErr("Error", err).Msg("Scanning rows...")
                }
                customers = append(customers, name)
	}
	//err = rows.Err()
        //logger.Fatal().AnErr("Error", err).Msg("Scanning rows...")
        return customers
}
