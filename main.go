package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	for _, person := range persons {
		LogInfo("MAIN", person.ID_CisP+" "+person.JMENO+" "+person.PRIJMENI)
	}
	for _, operation := range operations {
		LogInfo("MAIN", operation.ID+" "+operation.BARCODE+" "+operation.PRODUKT_NAZ)
	}

	orders, products, users, userTypes := DownloadDataFromZapsi()
	for _, order := range orders {
		LogInfo("MAIN", order.Name)
	}
	for _, product := range products {
		LogInfo("MAIN", product.Name)
	}
	for _, user := range users {
		LogInfo("MAIN", user.Name)
	}
	for _, userType := range userTypes {
		LogInfo("MAIN", userType.Name)
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
	LogInfo("MAIN", "Connected")
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
	LogInfo("MAIN", "Connected")
	var persons []ZAPSI_PERS
	db.Table("ZAPSI_PERS").Find(&persons)
	var operations []ZAPSI_OPERACE
	db.Table("ZAPSI_OPERACE").Find(&operations)
	return persons, operations
}
