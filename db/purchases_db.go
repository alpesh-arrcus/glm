package db

import (
	"database/sql"
	"fmt"
	"time"
)

func getCustomerPurchasesTableName(custName string) string {
	return custName + "_purchases"
}

func checkAndGetLicenseFromPurchases(db *sql.DB, custName, feature, fp string) (bool, error) {
	purchaseTblName := getCustomerPurchasesTableName(custName)
	licenseTblName := getCustomerLicenseAllocsTableName(custName)

	period, ok := checkLicenseAvailability(db, custName, feature)
	if !ok {
		logger.Error().Caller().Str("customer", custName).
			Str("feature", feature).Str("deviceFp", fp).Msg("No available license")
		return false, nil
	}

	tsNow := time.Now().UTC().Format(time.RFC3339)
	fetchLicStmt := fmt.Sprintf("update %s set licCount = licCount - 1 where featureName = '%s' and licCount > 0",
		purchaseTblName, feature)
	addAllocStmt := fmt.Sprintf("insert into %s (featureName, deviceFp, status, periodLeft, lastUse) values (?, ?, ?, ?, ?)",
		licenseTblName)

	//This has to be an atomic transaction
	//Decrement by 1 in purchaseTbl and add to licenseAllocsTable
	tx, err := db.Begin()
	if err != nil {
		logger.Error().Caller().AnErr("Error", err).
			Str("customer", custName).Str("feature", feature).Str("deviceFp", fp).
			Msg("Cannot create Txn")
		return false, err
	}
	res, err := tx.Exec(fetchLicStmt)
	if err != nil {
		tx.Rollback()
		logger.Error().Caller().Str("Stmt", fetchLicStmt).
			AnErr("Error", err).
			Msg("Executing Stmt failed")
		return false, err
	}
	res, err = tx.Exec(addAllocStmt, feature, fp, "InUse", period, tsNow)
	if err != nil {
		logger.Error().Caller().Str("Stmt", addAllocStmt).
			Str("customer", custName).Str("feature", feature).Str("deviceFp", fp).
			AnErr("Error", err).
			Msg("Executing Stmt failed")
		tx.Rollback()
		return false, err
	}
	_ = res
	err = tx.Commit()
	if err != nil {
		logger.Error().Caller().
			Str("customer", custName).Str("feature", feature).Str("deviceFp", fp).
			AnErr("Error", err).
			Msg("Executing Committing Tx")
		return false, err
	}
	return true, nil
}

func checkAddPurchase(db *sql.DB, custName, feature string, lcount, usage int) error {
	tblName := getCustomerPurchasesTableName(custName)
	sqlStmt := `
        create table if not exists ` + tblName + ` (featureName text not null, ` +
		` licCount int CHECK (licCount >= 0), ` +
		` usagePeriod int CHECK (usagePeriod >= 0), ` +
		` purchaseTime text );
        `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("creating Customer Purchase Table")
		return err
	}

	tsNow := time.Now().UTC().Format(time.RFC3339)
	addStmt := fmt.Sprintf("insert into %s (featureName, licCount, usagePeriod, purchaseTime) values (?, ?, ?, ?)", tblName)
	_, err = db.Exec(addStmt, feature, lcount, usage, tsNow)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("adding to customer Purchases..")
		return err
	}

	tblName = getCustomerLicenseAllocsTableName(custName)
	sqlStmt = `
        create table if not exists ` + tblName + ` (featureName text not null, ` +
		` deviceFp text not null, ` +
		` status text, ` +
		` periodLeft int CHECK (periodLeft >=0), ` +
		` lastUse text);
        `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logger.Error().AnErr("Error", err).Str("Sql Stmt", sqlStmt).Msg("creating Customer Allocs Table")
		return err
	}
	return nil
}

func AddCustomerPurchase(custName, feature string, licenseCount, usagePeriod int) error {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	err := checkAddCustomer(dbSt.db, custName)
	if err != nil {
		return err
	}
	err = checkAddPurchase(dbSt.db, custName, feature, licenseCount, usagePeriod)
	if err != nil {
		return err
	}
	return nil
}
