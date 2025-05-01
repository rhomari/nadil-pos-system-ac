package main

import (
	"encoding/json"
	"time"
)

type Category struct {
	Category_id   int    `json:"categoryid"`
	Category_name string `json:"categoryname"`
}
type Product struct {
	Productid           int     `json:"productid"`
	Isfavorite          bool    `json:"isfavorite"`
	Name                string  `json:"text"`
	Category            int     `json:"category"`
	Image               string  `json:"image"`
	Price               string  `json:"price"`
	Subscriberprice     *string `json:"subscriberprice"`
	Goldsubscriberprice *string `json:"goldsubscriberprice"`
}
type Response struct {
	Prefix  string `json:"Prefix"`
	Message string `json:"Message"`
	Data    string `json:"data"`
}
type CustomSession struct {
	Data map[string]interface{}
}

var globalsession = make(map[string]CustomSession)

type Settings struct {
	GeneralSettings  string `json:"generalsettings"`
	ManagerSettings  string `json:"managersettings"`
	PrintingSettings string `json:"printingsettings"`
	SecuritySettings string `json:"securitysettings"`
}
type PrintedObject struct {
	GeneralSettings string          `json:"GeneralSettings"`
	Profile         PrintingProfile `json:"PrintingProfile"`
	Order           Order
}
type SitutionObject struct {
	Situation string          `json:"Situation"`
	Profile   PrintingProfile `json:"PrintingProfile"`
}
type User struct {
	Username string
	Role     string
	Password string
	Code     string
	ID       string
}

type UserInfos struct {
	ID       string
	Username string
	Role     string
	Code     string
}
type Order struct {
	Number       int             `json:"number"`
	Totalticket  float64         `json:"totalticket"`
	Waiter       string          `json:"waiter"`
	Content      json.RawMessage `json:"content"`
	Creator      string          `json:"creator"`
	Status       string          `json:"status"`
	Creationdate time.Time       `json:"creationdate"`
	Member       string          `json:"member"`
	MemberID     string          `json:"memberid"`
}
type GeneralSettings struct {
	Posname  string `json:"posname"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Wificode string `json:"wificode"`
}
type OrderContent struct {
	Text       string  `json:"text"`
	Count      int     `json:"count"`
	Price      float64 `json:"price"`
	Comment    string  `json:"comment"`
	Categoryid int     `json:"categoryid"`
}

type PrintingProfile struct {
	Profile          string `json:"profile"`
	PrintingMode     string `json:"printingmode"`
	PrintingFunction string `json:"printingfunction"`
	PrintingOptions  []int  `json:"printingoptions"`
	PrinterIP        string `json:"printerip"`
	PrinterPort      string `json:"printerport"`
	Copies           int    `json:"copies"`
	Active           bool   `json:"active"`
}
type PrintingSettings struct {
	PrintingProfiles []PrintingProfile `json:"printingprofiles"`
}
type Session struct {
	Username  string    `json:"username"`
	SessionID string    `json:"sessionid"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
}
type Situation struct {
	TotalPaid      float64 `json:"totalpaid"`
	CountPaid      int     `json:"countpaid"`
	TotalInpaid    float64 `json:"totalinpaid"`
	CountInpaid    int     `json:"countinpaid"`
	TotalCanceled  float64 `json:"totalcanceled"`
	CountICanceled int     `json:"countcanceled"`
	Total          float64 `json:"total"`
}
type SituationToPrint struct {
	TotalPaid      float64   `json:"totalpaid"`
	CountPaid      int       `json:"countpaid"`
	TotalInpaid    float64   `json:"totalinpaid"`
	CountInpaid    int       `json:"countinpaid"`
	TotalCanceled  float64   `json:"totalcanceled"`
	CountICanceled int       `json:"countcanceled"`
	Total          float64   `json:"total"`
	Username       string    `json:"username"`
	Sessionstart   time.Time `json:"sessionstart"`
	Sessionend     time.Time `json:"sessionend"`
}
