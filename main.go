package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"strings"
	"time"
)

const version = "2021.1.3.3"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	SendMail("Program started", "Adoz-data-importer version "+version+" started")
	for {
		LogInfo("MAIN", "Program running")
		start := time.Now()
		persons, operations, K2DataDownloaded := DownloadDataFromK2()
		orders, users, zapsiDataDownloaded := DownloadDataFromZapsi()
		if K2DataDownloaded && zapsiDataDownloaded {
			LogInfo("MAIN", "Data download after "+time.Now().Sub(start).String())
			userProcessingStart := time.Now()
			ProcessUsers(persons, users)
			LogInfo("MAIN", "Users processed after "+time.Now().Sub(userProcessingStart).String())
			orderProcessingStart := time.Now()
			ProcessOperations(operations, orders)
			LogInfo("MAIN", "Orders processed after "+time.Now().Sub(orderProcessingStart).String())
		}
		LogInfo("MAIN", "Total time "+time.Now().Sub(start).String())
		LogInfo("MAIN", "Sleeping for "+(1*time.Minute-time.Now().Sub(start)).String())
		DeleteOldLogFiles()
		time.Sleep(1*time.Minute - time.Now().Sub(start))
	}

}

func ProcessOperations(operations []ZAPSI_OPERACE, orders []order) {
	for _, operation := range operations {
		operationInZapsi := false
		for idx, order := range orders {
			if strings.Trim(operation.BARCODE, " ") == order.Barcode {
				operationInZapsi = true
				orders = append(orders[0:idx], orders[idx+1:]...)
				break
			}
		}
		if !operationInZapsi {
			LogInfo("MAIN", "Adding operation ["+operation.BARCODE+"]")
			err := AddOrder(operation)
			if err != nil {
				LogError("MAIN", "Problem adding operation "+operation.BARCODE+": "+err.Error())
			}
		}
	}
}

func AddOrder(operation ZAPSI_OPERACE) error {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return err
	}
	defer db.Close()
	newProduct := product{}
	db.Table("product").Where("Name = ?", operation.PRODUKT_NAZ).First(&newProduct)
	if newProduct.OID == 0 {
		LogInfo("MAIN", "Adding new product "+operation.PRODUKT_NAZ)
		newProduct.Barcode = operation.PRODUKT_ZKR
		db.NewRecord(newProduct)
		db.Table("product").Create(&newProduct)
	}
	alteredProduct := product{}
	db.Table("product").Where("Name = ?", operation.PRODUKT_NAZ).First(&alteredProduct)
	productOID := alteredProduct.OID
	countRequested, countError := strconv.Atoi(operation.PLAN_KS)
	if countError != nil {
		LogError("MAIN", "Problem parsing plan_ks "+operation.PLAN_KS+", "+countError.Error())
		countRequested = 0
	}
	trimmedNormaPrip := strings.Replace(operation.NORMA_PRIP, ",", "", -1)
	opNormaPrip, normaPripErr := strconv.ParseFloat(trimmedNormaPrip, 64)
	if normaPripErr != nil {
		LogError("MAIN", "Problem parsing norma_prip "+operation.NORMA_PRIP+", "+normaPripErr.Error())
		opNormaPrip = 0
	}
	trimmedNormaVyr := strings.Replace(operation.NORMA_VYR, ",", "", -1)
	opNormaVyr, normaVyrErr := strconv.ParseFloat(trimmedNormaVyr, 64)
	if normaVyrErr != nil {
		LogError("MAIN", "Problem parsing vyr "+operation.NORMA_PRIP+", "+normaVyrErr.Error())
		opNormaVyr = 0
	}
	newOrder := order{Name: strings.Trim(operation.BARCODE, " "), Barcode: strings.Trim(operation.BARCODE, " "), Pruvodka: operation.PRUVODKA,
		OpCode: operation.OPCODE, CountRequested: countRequested, OpNormaPrip: opNormaPrip,
		OpNormaVyr: opNormaVyr, ProductID: productOID, OrderStatusID: 1, Cavity: 1}
	db.Debug().NewRecord(newOrder)
	db.Debug().Table("order").Create(&newOrder)
	return nil
}

func ProcessUsers(persons []ZAPSI_PERS, users []user) {
	for _, person := range persons {
		personInZapsi := false
		for idx, zapsiUser := range users {
			if person.ID_CisP == zapsiUser.Barcode {
				personInZapsi = true
				err := UpdateUser(zapsiUser, person)
				users = append(users[0:idx], users[idx+1:]...)
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
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return err
	}
	newUser := user{FirstName: person.JMENO, Name: person.PRIJMENI, Rfid: person.RFID, Login: person.K2_UZIV, Barcode: person.ID_CisP, UserRoleID: 2}
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
	db.NewRecord(newUser)
	db.Table("user").Create(&newUser)
	return nil
}

func UpdateUser(zapsiUser user, person ZAPSI_PERS) error {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return err
	}
	zapsiUser.Rfid = person.RFID
	db.Table("user").Save(&zapsiUser)
	return nil
}

func DownloadDataFromZapsi() ([]order, []user, bool) {
	connectionString := "zapsi_uzivatel:zapsi@tcp(localhost:3306)/zapsi2?charset=utf8&parseTime=True&loc=Local"
	dialect := "mysql"
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil, false
	}
	LogInfo("MAIN", "Zapsi database connected")
	var orders []order
	db.Table("order").Find(&orders)
	var users []user
	db.Table("user").Find(&users)
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(orders))+" orders")
	LogInfo("MAIN", "Zapsi Found "+strconv.Itoa(len(users))+" users")
	return orders, users, true
}

func DownloadDataFromK2() ([]ZAPSI_PERS, []ZAPSI_OPERACE, bool) {
	connectionString := "sqlserver://zapsi:RuruRavePivo92@sql:1433?database=K2_ADOZ"
	dialect := "mssql"
	db, err := gorm.Open(dialect, connectionString)
	defer db.Close()
	if err != nil {
		LogError("MAIN", "Problem opening database "+connectionString+", "+err.Error())
		return nil, nil, false
	}
	LogInfo("MAIN", "K2 database connected")
	var persons []ZAPSI_PERS
	db.Table("ZAPSI_PERS").Find(&persons)
	var operations []ZAPSI_OPERACE
	db.Table("ZAPSI_OPERACE").Find(&operations)
	LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(persons))+" persons")
	LogInfo("MAIN", "K2 Found "+strconv.Itoa(len(operations))+" operations")
	return persons, operations, true
}
