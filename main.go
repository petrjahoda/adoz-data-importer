package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
	"time"
)

const version = "2019.4.1.10"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	SendMail("Program started", "Adoz-data-importer version "+version+" started")
	LogInfo("MAIN", "Program running")
	persons, operations := DownloadDataFromK2()
	orders, products, users, userTypes := DownloadDataFromZapsi()
	LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(persons))+" persons")
	LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(operations))+" operations")
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(orders))+" orders")
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(products))+" products")
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(users))+" users")
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(userTypes))+" user types")

	for _, person := range persons {
		personInZapsi := false
		for _, user := range users {
			if person.ID_CisP == user.Barcode {
				personInZapsi = true
				LogInfo("MAIN", "Person match found: "+person.JMENO+" "+person.PRIJMENI)
				break
			}
		}
		if !personInZapsi {
			LogInfo("MAIN", "Adding person "+person.JMENO+" "+person.PRIJMENI)
		}
	}

	for _, operation := range operations {
		operationInZapsi := false
		for _, order := range orders {
			if operation.BARCODE == order.Barcode {
				operationInZapsi = true
				LogInfo("MAIN", "Operation match found: "+operation.BARCODE+" "+operation.OPCODE)

				break
			}
		}
		if !operationInZapsi {
			LogInfo("MAIN", "Adding operation "+operation.BARCODE+" "+operation.OPCODE)
		}
	}

}

func DownloadDataFromZapsi() ([]order, []product, []user, []user_type) {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil, nil, nil
	}
	defer db.Close()
	LogInfo("MAIN", "Zapsi database connected")
	var orders []order
	db.Table("order").Find(&orders)
	var products []product
	db.Table("product").Find(&products)
	var users []user
	db.Table("user").Find(&users)
	var userTypes []user_type
	db.Table("user_type").Find(&userTypes)
	return orders, products, users, userTypes
}

func DownloadDataFromK2() ([]ZAPSI_PERS, []ZAPSI_OPERACE) {
	connectionString := "sqlserver://zapsi:RuruRavePivo92@sql:1433?database=K2_ADOZ"
	dialect := "mssql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil
	}
	defer db.Close()
	LogInfo("MAIN", "K2 database connected")
	var persons []ZAPSI_PERS
	db.Table("ZAPSI_PERS").Find(&persons)
	var operations []ZAPSI_OPERACE
	db.Table("ZAPSI_OPERACE").Find(&operations)
	return persons, operations
}
