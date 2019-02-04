package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/toravir/glm/config"
	. "github.com/toravir/glm/context"
	"time"
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

	sqlStmt := `
        create table if not exists customers (name text not null primary key, secret text not null, id integer not null, status text not null);
        `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("Creating Customer Table")
	}
	return ctx
}

func GetCustomerNames() []string {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	customers := []string{}
	rows, _ := dbSt.db.Query("select name from customers")
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			logger.Error().AnErr("Error", err).Msg("Scanning rows...")
		}
		customers = append(customers, name)
	}
	return customers
}

func IsValidCustomer(custName string) bool {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	validCustomer := false

	rows, _ := dbSt.db.Query("select status from customers where name = ?", custName)
	defer rows.Close()
	for rows.Next() {
		var status string
		err := rows.Scan(&status)
		if err == nil {
			if status == "Active" {
				validCustomer = true
				break
			}
		}
	}
	logger.Debug().Caller().Str("Customer", custName).Bool("valid", validCustomer).Msg("")
	return validCustomer
}

func IsValidCustomerSecret(custName string, inSecret string) bool {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	validCustomer := false

	rows, _ := dbSt.db.Query("select secret from customers where name = ?", custName)
	defer rows.Close()
	for rows.Next() {
		var expSecret string
		err := rows.Scan(&expSecret)
		if err == nil {
			if expSecret == inSecret {
				validCustomer = true
				break
			}
		}
	}
	logger.Debug().Caller().Str("Customer", custName).Bool("validSecret", validCustomer).Msg("")
	return validCustomer
}

func createCustDevicesDb(db *sql.DB, custName string) (string, error) {
	tblName := custName + "_devices"
	sqlStmt := `create table if not exists ` + tblName + `(fp text not null primary key, lastHB text, status text not null); 
        `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		logger.Debug().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("Creating Customer Table")
		return tblName, err
	}
	return tblName, nil
}

func getDeviceStatus(db *sql.DB, custName, deviceFp string) (lastHb, status string) {
	tblName := custName + "_devices"
	qryStmt := fmt.Sprintf("select lastHB, status from %s where fp = ?", tblName)
	rows, _ := db.Query(qryStmt, deviceFp)
	status = ""
	lastHb = ""
	for rows.Next() {
		err := rows.Scan(&lastHb, &status)
		if err == nil {
			break
		}
	}
	rows.Close()
	return
}

func AddDevice(custName string, deviceFp string) (bool, bool) {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	tblName, _ := createCustDevicesDb(dbSt.db, custName)
	lastHb, status := getDeviceStatus(dbSt.db, custName, deviceFp)

	tsNow := time.Now().UTC().Format(time.RFC3339)
	isnew := false
	if status != "" {
		logger.Debug().Caller().Str("Customer", custName).
			Str("Device", deviceFp).Str("DBStatus", status).
			Str("lastHb", lastHb).
			Msg("Found device..")
		if status == "RMA" {
			return false, false
		}
		updStmt := fmt.Sprintf("update %s set lastHB = '%s' where fp = '%s'", tblName, tsNow, deviceFp)
		_, err := dbSt.db.Exec(updStmt)
		if err != nil {
			logger.Debug().Caller().AnErr("Error", err).Str("SqlStmt", updStmt).Msg("Update device in DB")
			return false, false
		}
	} else {
		isnew = true
		logger.Debug().Caller().Str("Customer", custName).Str("Device", deviceFp).Msg("Adding device..")
		addStmt := fmt.Sprintf("insert into %s (fp, lastHB, status) values ('%s', '%s', '%s')", tblName, deviceFp, tsNow, "Active")
		_, err := dbSt.db.Exec(addStmt)
		if err != nil {
			logger.Debug().Caller().AnErr("Error", err).Str("SqlStmt", addStmt).Msg("Adding device to DB")
			return false, false
		}
	}
	return true, isnew
}

func AllocateLicense(custName, deviceFp, feature string) bool {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	lastHb, status := getDeviceStatus(dbSt.db, custName, deviceFp)
	_ = lastHb
	_ = status

	return true
}
