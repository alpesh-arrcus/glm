package db

import (
	"database/sql"
	"fmt"
)

func getCustomersTableName() string {
	return "customers"
}

func GetCustomerNames() []string {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	customers := []string{}
	tblName := getCustomersTableName()
	rows, err := dbSt.db.Query("select customerName from " + tblName)
	if err != nil {
		logger.Error().AnErr("Error", err).Msg("Querying list of Customers..")
		return customers
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			logger.Error().AnErr("Error", err).Msg("Scanning rows...")
		}
		customers = append(customers, name)
	}
	rows.Close()
	return customers
}

func IsValidCustomer(custName string) (status string, ok bool) {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	status = ""
	ok = false

	tblName := getCustomersTableName()
	rows, err := dbSt.db.Query("select status from "+tblName+" where customerName = ?", custName)
	if err != nil {
		logger.Error().AnErr("Error", err).Msg("Querying list of Customers..")
		return
	}
	for rows.Next() {
		err := rows.Scan(&status)
		if err == nil {
			if status == "Active" {
				ok = true
				break
			}
		}
	}
	rows.Close()
	logger.Debug().Caller().Str("Customer", custName).
		Str("Status", status).
		Bool("valid", ok).Msg("Checking customer")
	return
}

func IsValidCustomerSecret(custName string, inSecret string) bool {
	dbSt, _ := ctxCache.DbInfo.(*dbState)
	validCustomer := false

	tblName := getCustomersTableName()
	rows, err := dbSt.db.Query("select customerSecret from "+tblName+" where customerName = ?", custName)
	if err != nil {
		logger.Error().AnErr("Error", err).Msg("Querying list of Customers..")
	}
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
	rows.Close()
	logger.Debug().Caller().Str("Customer", custName).Bool("validSecret", validCustomer).Msg("")
	return validCustomer
}

func checkAddCustomer(db *sql.DB, custName string) error {
	st, ok := IsValidCustomer(custName)
	if ok {
		return nil
	}
	if st != "" {
		//TODO - this customer is NOT active - will need to update
		return nil
	}
	tblName := getCustomersTableName()
	addStmt := fmt.Sprintf("insert into %s (customerName, customerSecret, id, status) values (?, ?, ?, ?)", tblName)
	defaultCustSecret := custName + "123"
	_, err := db.Exec(addStmt, custName, defaultCustSecret, 0, "Active")
	if err != nil {
		logger.Debug().Caller().AnErr("Error", err).Str("SqlStmt", addStmt).Msg("Adding Customer")
		return err
	}
	return nil
}
