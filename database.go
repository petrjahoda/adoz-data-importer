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

type user struct {
	OID        int    `gorm:"column:OID"`
	Login      string `gorm:"column:Login"`
	Password   string `gorm:"column:Password"`
	Name       string `gorm:"column:Name"`
	FirstName  string `gorm:"column:FirstName"`
	Rfid       string `gorm:"column:Rfid"`
	Barcode    string `gorm:"column:Barcode"`
	Pin        string `gorm:"column:Pin"`
	Function   string `gorm:"column:Function"`
	UserTypeID int    `gorm:"column:UserTypeID"`
	Email      string `gorm:"column:Email"`
	Phone      string `gorm:"column:Phone"`
	UserRoleID int    `gorm:"column:UserRoleID"`
}

type user_type struct {
	OID  int    `gorm:"column:OID"`
	Name string `gorm:"column:Name"`
}

type order struct {
	OID            int     `gorm:"column:OID"`
	Name           string  `gorm:"column:Name"`
	Barcode        string  `gorm:"column:Barcode"`
	ProductID      int     `gorm:"column:ProductID"`
	OrderStatusID  int     `gorm:"column:OrderStatusID"`
	CountRequested int     `gorm:"column:CountRequested"`
	Cavity         int     `gorm:"column:Cavity"`
	OpCode         string  `gorm:"column:OpCode"`
	OpNormaVyr     float32 `gorm:"column:OpNormaVyr"`
	OpNormaPrip    float32 `gorm:"column:OpNormaPrip"`
	Pruvodka       string  `gorm:"column:Pruvodka"`
}

type product struct {
	OID             int     `gorm:"column:OID"`
	Name            string  `gorm:"column:Name"`
	Barcode         string  `gorm:"column:Barcode"`
	Cycle           float32 `gorm:"column:Cycle"`
	IdleFromTime    int     `gorm:"column:IdleFromTime"`
	ProductStatusID int     `gorm:"column:ProductStatusID"`
	Deleted         int     `gorm:"column:Deleted"`
	ProductGroupID  int     `gorm:"column:ProductGroupID"`
}
