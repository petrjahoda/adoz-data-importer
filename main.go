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
	SendMail("Program started", "Zapsi Service version "+version+" started")
	LogInfo("MAIN", "Program running")
	persons, operations := DownloadDataFromK2()

	for _, person := range persons {
		LogInfo("MAIN", person.ID_CisP+" "+person.JMENO+" "+person.PRIJMENI)
	}
	for _, operation := range operations {
		LogInfo("MAIN", operation.ID+" "+operation.BARCODE+" "+operation.PRODUKT_NAZ)
	}
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
