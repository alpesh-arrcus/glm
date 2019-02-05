package db

import (
	"database/sql"
	"fmt"
	"time"
)

func getCustomerLicenseAllocsTableName(custName string) string {
	return custName + "_licenseAllocs"
}

func lookupLicenseAllocs(db *sql.DB, custName, deviceFp, feature string) (lastUse string, status string, periodLeft int) {
	status = ""
	lastUse = ""
	periodLeft = 0

	tblName := getCustomerLicenseAllocsTableName(custName)
	qryStmt := fmt.Sprintf("select lastUse, status, periodLeft from %s where devicefp = ? and featureName = ?", tblName)
	rows, err := db.Query(qryStmt, deviceFp, feature)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("performing query")
		return
	}
	for rows.Next() {
		err := rows.Scan(&lastUse, &status, &periodLeft)
		if err != nil {
			logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("performing rows scan")
			continue
		}
		if status == "Available" {
			break
		}
		logger.Info().Str("Status", status).
			Int("periodLeft", periodLeft).Str("Query", qryStmt).
			Msg("performing rows scan")
	}
	rows.Close()
	return
}

func updateLicenseAllocs(db *sql.DB, custName, deviceFp, feature, curTime string, periodLeft int) (ok bool) {
	tblName := getCustomerLicenseAllocsTableName(custName)
	status := "Expired"
	if periodLeft > 0 {
		status = "InUse"
	}
	qryStmt := fmt.Sprintf("update %s set lastUse = '%s', status = '%s', periodLeft = %d where devicefp = ? and featureName = ?", tblName,
		curTime, status, periodLeft)
	_, err := db.Exec(qryStmt, deviceFp, feature)
	if err != nil {
		logger.Error().Caller().Str("Stmt", qryStmt).AnErr("Error", err).
			Str("Customer", custName).Str("FP", deviceFp).
			Msg("Update table failed")
		return false
	}
	return true
}

func checkLicenseAvailability(db *sql.DB, custName, feature string) (usagePeriod int, ok bool) {
	count := 0
	usagePeriod = 0
	ok = false

	purchaseTblName := getCustomerPurchasesTableName(custName)
	qryStmt := fmt.Sprintf("select licCount, usagePeriod from %s where featureName = ?", purchaseTblName)

	rows, err := db.Query(qryStmt, feature)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("Check license Availability")
		return
	}
	for rows.Next() {
		err := rows.Scan(&count, &usagePeriod)
		if err != nil {
			logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("performing rows scan")
			continue
		}
		logger.Info().
			Int("Count", count).Str("Query", qryStmt).
			Msg("performing rows scan")
		if count > 0 {
			ok = true
			break
		}
	}
	rows.Close()
	return
}

func AllocateLicense(custName, deviceFp, feature string) bool {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	_, status := getDeviceStatus(dbSt.db, custName, deviceFp)

	if status != "Active" {
		logger.Error().Str("Status", status).
			Str("deviceFp", deviceFp).
			Msg("allocate license for a non-active device...")
		return false
	}
	//Check if was already allocated
	_, _, pl := lookupLicenseAllocs(dbSt.db, custName, deviceFp, feature)
	if pl > 0 {
		//Start the time
		tsNow := time.Now().UTC().Format(time.RFC3339)
		if !updateLicenseAllocs(dbSt.db, custName, deviceFp, feature, tsNow, pl) {
			return false
		}
		return true
	}
	//No allocations found
	ok, _ := checkAndGetLicenseFromPurchases(dbSt.db, custName, feature, deviceFp)
	return ok
}

func updateLicenseUsage(custName, deviceFp string, autoRealloc bool, secsToSub int) (expiredLics []string, err error) {
	expiredLics = []string{}
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	db := dbSt.db

	tblName := getCustomerLicenseAllocsTableName(custName)

	qryStmt := fmt.Sprintf("select featureName from %s "+
		"where devicefp = ? and periodLeft < %d and status = 'InUse'", tblName,
		secsToSub)
	rows, err := db.Query(qryStmt, deviceFp)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("Searching for expiring Lics")
		return
	}
	featureName := ""
	for rows.Next() {
		err := rows.Scan(&featureName)
		if err != nil {
			logger.Error().AnErr("Error", err).Str("Query", qryStmt).Msg("performing rows scan")
			continue
		}
		expiredLics = append(expiredLics, featureName)
	}
	rows.Close()

	//Set expired Licenses
	qryStmt = fmt.Sprintf("update %s set periodLeft = 0, status = 'Expired' "+
		"where devicefp = ? and periodLeft < %d and status = 'InUse'", tblName,
		secsToSub)
	_, err = db.Exec(qryStmt, deviceFp)
	if err != nil {
		logger.Error().Caller().Str("Stmt", qryStmt).AnErr("Error", err).
			Str("Customer", custName).Str("FP", deviceFp).
			Msg("Update table failed")
		return
	}

	//decrement usage for non expiring licenses
	qryStmt = fmt.Sprintf("update %s set periodLeft = periodLeft - %d "+
		"where devicefp = ? and periodLeft >= %d and status = 'InUse'", tblName,
		secsToSub, secsToSub)
	_, err = db.Exec(qryStmt, deviceFp)
	if err != nil {
		logger.Error().Caller().Str("Stmt", qryStmt).AnErr("Error", err).
			Str("Customer", custName).Str("FP", deviceFp).
			Msg("Update table failed")
		return
	}

	if !autoRealloc {
		logger.Debug().Msgf("No Attempt to Renew these licenses: %s", expiredLics)
		err = nil
		return
	}
	logger.Debug().Msgf("Attempting to Renew these licenses:", expiredLics)
	attemptToRealloc := expiredLics
	expiredLics = []string{}

	for _, fn := range attemptToRealloc {
		if !AllocateLicense(custName, deviceFp, fn) {
			expiredLics = append(expiredLics, fn)
			logger.Debug().Str("Customer", custName).Str("FP", deviceFp).
				Str("Feature", fn).
				Msg("License Expired.")
		} else {
			logger.Debug().Str("Customer", custName).Str("FP", deviceFp).
				Str("Feature", fn).
				Msg("License AutoReAlloced.")
		}
	}
	err = nil
	return
}
