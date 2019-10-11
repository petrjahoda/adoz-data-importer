package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
	"strings"
	"time"
)

const version = "2019.4.1.11"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	SendMail("Program started", "Adoz-data-importer version "+version+" started")
	LogInfo("MAIN", "Program running")
	for {
		start := time.Now()
		persons, operations, K2DataDownloaded := DownloadDataFromK2()
		orders, products, users, userTypes, zapsiDataDownloaded := DownloadDataFromZapsi()

		if K2DataDownloaded && zapsiDataDownloaded {
			LogInfo("MAIN", "Data download after "+time.Now().Sub(start).String())
			LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(persons))+" persons")
			LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(operations))+" operations")
			LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(orders))+" orders")
			LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(products))+" products")
			LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(users))+" users")
			LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(userTypes))+" user types")

			userProcessingStart := time.Now()
			ProcessUsers(persons, users)
			LogInfo("MAIN", "Users processed after "+time.Now().Sub(userProcessingStart).String())

			orderProcessingStart := time.Now()
			for _, operation := range operations {
				operationInZapsi := false
				for _, order := range orders {
					if strings.Trim(operation.BARCODE, " ") == order.Barcode {
						operationInZapsi = true
						break
					}
				}
				if !operationInZapsi {
					LogInfo("MAIN", "Adding operation ["+operation.BARCODE+"]")
					// TODO: check for product, add product
					// TODO: get product id
					// TODO: add order with product
				}
			}
			LogInfo("MAIN", "Orders processed after "+time.Now().Sub(orderProcessingStart).String())
		}
		LogInfo("MAIN", "Total time "+time.Now().Sub(start).String())
		LogInfo("MAIN", "Sleeping for "+(1*time.Minute-time.Now().Sub(start)).String())
		time.Sleep(1*time.Minute - time.Now().Sub(start))
	}

}

func ProcessUsers(persons []ZAPSI_PERS, users []user) {
	for _, person := range persons {
		personInZapsi := false
		for _, zapsiUser := range users {
			if person.ID_CisP == zapsiUser.Barcode {
				personInZapsi = true
				LogInfo("MAIN", person.JMENO+" "+person.PRIJMENI+": updating rfid to ["+person.RFID+"]")
				err := UpdateUser(zapsiUser, person)
				if err != nil {
					LogError("MAIN", "Problem updating user: "+err.Error())
				}
				break
			}
		}
		if !personInZapsi {
			err := AddUser(person)
			if err != nil {
				LogError("MAIN", "Problem adding user "+person.JMENO+" "+person.PRIJMENI+": "+err.Error())
			}
		}
	}
}

func AddUser(person ZAPSI_PERS) error {
	LogInfo("MAIN", "Adding person "+person.JMENO+" "+person.PRIJMENI)
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return err
	}
	defer db.Close()
	newUser := user{FirstName: person.JMENO, Name: person.PRIJMENI, Rfid: person.RFID, Login: person.K2_UZIV, Barcode: person.ID_CisP}
	switch person.SKUPINA {
	case "kvalita":
		newUser.UserTypeID = 2
	case "manažer":
		newUser.UserTypeID = 3
	case "operátor":
		newUser.UserTypeID = 4
	case "serizovac":
		newUser.UserTypeID = 5
	case "technolog":
		newUser.UserTypeID = 6
	case "údržbár":
		newUser.UserTypeID = 7
	default:
		newUser.UserTypeID = 1
	}
	db.Table("user").NewRecord(newUser)
	db.Create(&newUser)
	return nil
}

func UpdateUser(zapsiUser user, person ZAPSI_PERS) error {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return err
	}
	defer db.Close()
	zapsiUser.Rfid = person.RFID
	db.Table("user").Save(&zapsiUser)
	return nil
}

func DownloadDataFromZapsi() ([]order, []product, []user, []user_type, bool) {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil, nil, nil, false
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
	return orders, products, users, userTypes, true
}

func DownloadDataFromK2() ([]ZAPSI_PERS, []ZAPSI_OPERACE, bool) {
	connectionString := "sqlserver://zapsi:RuruRavePivo92@sql:1433?database=K2_ADOZ"
	dialect := "mssql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil, false
	}
	defer db.Close()
	LogInfo("MAIN", "K2 database connected")
	var persons []ZAPSI_PERS
	db.Table("ZAPSI_PERS").Find(&persons)
	var operations []ZAPSI_OPERACE
	db.Table("ZAPSI_OPERACE").Find(&operations)
	return persons, operations, true
}
