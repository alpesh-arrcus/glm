package db

import (
	"database/sql"
	"fmt"
	"time"
)

func getCustomerDevicesTableName(custName string) string {
	return custName + "_devices"
}

func createCustDevicesDb(db *sql.DB, custName string) (string, error) {
	tblName := getCustomerDevicesTableName(custName)
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
	status = ""
	lastHb = ""
	tblName := getCustomerDevicesTableName(custName)
	qryStmt := fmt.Sprintf("select lastHB, status from %s where fp = ?", tblName)
	rows, err := db.Query(qryStmt, deviceFp)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("performing query")
		return
	}
	for rows.Next() {
		err := rows.Scan(&lastHb, &status)
		if err == nil {
			break
		}
	}
	rows.Close()
	return
}

func updateDevicesDb(db *sql.DB, custName, deviceFp, tsNow string) bool {
	tblName := getCustomerDevicesTableName(custName)
	updStmt := fmt.Sprintf("update %s set lastHB = '%s' where fp = '%s'", tblName, tsNow, deviceFp)
	_, err := db.Exec(updStmt)
	if err != nil {
		logger.Debug().Caller().AnErr("Error", err).Str("SqlStmt", updStmt).Msg("Update device in DB")
		return false
	}
	return true
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
		if !updateDevicesDb(dbSt.db, custName, deviceFp, tsNow) {
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

func DeviceHeartBeat(custName, deviceFp string, autoRealloc bool) (expiredLics []string, err error) {
	expiredLics = []string{}
	dbSt, _ := ctxCache.DbInfo.(*dbState)

	lastHb, status := getDeviceStatus(dbSt.db, custName, deviceFp)
	if status != "Active" {
		logger.Error().Str("Status", status).
			Str("deviceFp", deviceFp).
			Msg("heartbeat from a non-active device...")

		err = fmt.Errorf("Device not active - but getting HB...")
		return
	}

	timeNow := time.Now().UTC()
	timeLast, err := time.Parse(time.RFC3339, lastHb)
	if err != nil {
		logger.Error().Str("lastHB", lastHb).
			AnErr("ParseError", err).
			Msg("Parsing last HB...")
		timeLast = timeNow
	}
	tsNow := timeNow.Format(time.RFC3339)
	if !updateDevicesDb(dbSt.db, custName, deviceFp, tsNow) {
		err = fmt.Errorf("Error updating Devices Db!")
		return
	}
	elapsed := timeNow.Sub(timeLast)
	secondsToSub := int(elapsed.Seconds())
	expiredLics, err = updateLicenseUsage(custName, deviceFp, autoRealloc, secondsToSub)
	return
}
