package main

type ZAPSI_PERS struct {
	ID_CisP  string `gorm:"column:ID_CisP"`  //user.barcode
	ID_OsCP  string `gorm:"column:ID_OsCP"`  // user.login
	JMENO    string `gorm:"column:JMENO"`    // user.firstname
	PRIJMENI string `gorm:"column:PRIJMENI"` // user.name
	SKUPINA  string `gorm:"column:SKUPINA"`  // user_type.name ---> user_type.oid = user.usertypeid
	RFID     string `gorm:"column:RFID"`     // user.rfid
	K2_UZIV  string `gorm:"column:K2_UZIV"`  // user.login
}

type ZAPSI_OPERACE struct {
	ID          string `gorm:"column:ID"`
	PRUVODKA    string `gorm:"column:PRUVODKA"`    // order.Pruvodka
	BARCODE     string `gorm:"column:BARCODE"`     // order.name, order.barcode
	OPCODE      string `gorm:"column:OPCODE"`      // order.opcode
	NORMA_VYR   string `gorm:"column:NORMA_VYR"`   // order.OpNormaVyr
	NORMA_PRIP  string `gorm:"column:NORMA_PRIP"`  // order.OpNormaPrip
	PRODUKT_ZKR string `gorm:"column:PRODUKT_ZKR"` // product.barcode
	PRODUKT_NAZ string `gorm:"column:PRODUKT_NAZ"` // product.name
	PLAN_KS     string `gorm:"column:PLAN_KS"`     // order.CountRequested
	OK          string `gorm:"column:OK"`
	NOK         string `gorm:"column:NOK"`
	POTVRZENI   string `gorm:"column:POTVRZENI"`
}
