package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/charmap"
)

var biscuit http.Cookie
var usersinfos []UserInfos
var settings Settings
var PS []PrintingProfile

func NewUUID() string {
	return uuid.NewV4().String()

}

func setupRoutes() {

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir := filepath.Dir(executablePath)

	staticDir := filepath.Join(executableDir, "static")

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)
	http.HandleFunc("/getcategories", getCategories)
	http.HandleFunc("/saveNewCategory", server.saveNewCategory) //admin only
	http.HandleFunc("/deletecategory", server.deleteCategory)   //admin only
	http.HandleFunc("/editCategory", server.editCategory)       //admin only
	http.HandleFunc("/addNewProduct", server.addNewProduct)     //admin only
	http.HandleFunc("/editProduct", server.editProduct)         //admin only
	http.HandleFunc("/getMenu", server.getMenu)
	http.HandleFunc("/deleteProduct", server.deleteProduct) //admin only
	http.HandleFunc("/saveTicket", server.saveTicket)
	http.HandleFunc("/adminlogin", server.adminLogin)
	http.HandleFunc("/saveuser", server.saveUser) //admin only
	http.HandleFunc("/getusers", server.getUsers)
	http.HandleFunc("/edituser", server.editUser)     //admin only
	http.HandleFunc("/deleteuser", server.deleteUser) //admin only
	http.HandleFunc("/auth", server.auth)
	http.HandleFunc("/opensession", server.openSession)
	http.HandleFunc("/endsession", server.endSession)
	http.HandleFunc("/getOrders", server.getOrders)
	http.HandleFunc("/printreceipt", server.printReceipt)
	http.HandleFunc("/cancelorder", cancelOrder)
	http.HandleFunc("/checkoutorder", checkoutOrder)
	http.HandleFunc("/savegeneralsettings", savegeneralsettings)
	http.HandleFunc("/loadsettings", loadSettings)
	http.HandleFunc("/notauthorised", notAuthorized)
	http.HandleFunc("/savemanagersettings", saveManagerSettings)
	http.HandleFunc("/saveprintingsettings", server.savePrintingSettings)
	http.HandleFunc("/changesessionstatus", changeSessionStatus)
	http.HandleFunc("/getsessionsstatus", getSessionsStatus)
	http.HandleFunc("/getsituation", getSituation)
	http.HandleFunc("/printsituation", server.printReceipt)
	http.HandleFunc("/changeproductcategory", server.changeProductCategory)
	http.HandleFunc("/savesecuritysettings", saveSecuritySettings)
	log.Println("[setupRoutes] : Succes")
}
func (server *Server) changeProductCategory(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	server.connectToDB()

	_, err := server.DATABASE.Exec("UPDATE products_table SET category=$1 WHERE productid=$2;", r.FormValue("categoryid"), r.FormValue("elementid"))
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "UPDATE", "Catégorie modifiée", "UPDATE_PC_OK"))

}
func networkPrintSituation(situation string, printerIP string, printerPort string, copies int) {

	address := net.JoinHostPort(printerIP, printerPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Erreur de connexion à l'imprimante: %v\n", err)
		return
	}
	defer conn.Close()

	receiptData := formatSituation(situation)
	for i := 1; i <= copies; i++ {
		_, err = conn.Write(receiptData)
		if err != nil {
			log.Printf("Echec d'envoi des données du reçu: %v\n", err)
			return
		}
	}

}
func formatSituation(situation string) []byte {
	var stc SituationToPrint

	err := json.Unmarshal([]byte(situation), &stc)
	if err != nil {
		log.Println(err)

	}

	initPrinter := []byte{0x1B, 0x40}
	selectCharset := []byte{0x1B, 0x52, 0x00}
	alignCenter := []byte{0x1B, 0x61, 0x01}
	alignLeft := []byte{0x1B, 0x61, 0x00}

	boldOn := []byte{0x1B, 0x45, 0x01}
	boldOff := []byte{0x1B, 0x45, 0x00}

	doubleHeightOn := []byte{0x1D, 0x21, 0x10}
	doubleHeightOff := []byte{0x1D, 0x21, 0x00}
	cutPaper := []byte{0x1D, 0x56, 0x00}
	separator := []byte("------------------------------------------------\n")

	storeHeader := append(append(append(
		alignCenter,
		doubleHeightOn...),
		[]byte(toCP437("Situation des commandes"+"\n\n"))...),
		doubleHeightOff...,
	)
	var items []byte
	items = append(items, []byte(spacePAD(string(toCP437("Payées")), 30, strconv.FormatInt(int64(stc.CountPaid), 10), 15, strconv.FormatFloat(stc.TotalPaid, 'f', 2, 32), 3, ""))...)
	items = append(items, []byte(spacePAD(string(toCP437("Impayées")), 30, strconv.FormatInt(int64(stc.CountInpaid), 10), 15, strconv.FormatFloat(stc.TotalInpaid, 'f', 2, 32), 3, ""))...)
	items = append(items, []byte(spacePAD(string(toCP437("Annulées")), 30, strconv.FormatInt(int64(stc.CountICanceled), 10), 15, strconv.FormatFloat(stc.TotalCanceled, 'f', 2, 32), 3, ""))...)

	storeHeader = append(storeHeader, alignLeft...)
	storeHeader = append(storeHeader, []byte(toCP437("Utilisateur : "+stc.Username+"\n"))...)
	storeHeader = append(storeHeader, separator...)
	storeHeader = append(storeHeader, []byte(toCP437("Début : "+stc.Sessionstart.Format("02-01-2006 15:04:05")+"\n"))...)

	storeHeader = append(storeHeader, []byte(toCP437("Fin   : "+stc.Sessionend.Format("02-01-2006 15:04:05")+"\n"))...)
	storeHeader = append(storeHeader, separator...)

	itemsHeader := append(append(alignLeft,
		[]byte(toCP437("Status  commandes       Nombre          Total\n"))...),
		separator...,
	)

	totals := []byte(spaceSoloPAD("Total :", 48, strconv.FormatFloat(float64(stc.Total), 'f', 2, 32)+" MAD\n"))

	receipt := append(initPrinter, selectCharset...)
	receipt = append(receipt, storeHeader...)
	receipt = append(receipt, itemsHeader...)
	receipt = append(receipt, items...)
	receipt = append(receipt, separator...)
	receipt = append(receipt, boldOn...)
	receipt = append(receipt, totals...)
	receipt = append(receipt, doubleHeightOff...)
	receipt = append(receipt, boldOff...)
	receipt = append(receipt, separator...)
	receipt = append(receipt, []byte("\n\n\n\n\n\n")...)
	receipt = append(receipt, cutPaper...)

	return receipt
}

func getSituation(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}

	role, err := getUserRole(r.FormValue("username"))
	if err != nil {
		log.Println(err)
		return
	}

	server.connectToDB()

	start := r.FormValue("start")

	var rows *sql.Rows
	if role == "Gérant" {

		rows, err = server.DATABASE.Query("SELECT * FROM orders WHERE creator=$1 AND waiter='' AND creationdate >= $2 ORDER BY number;", r.FormValue("username"), start)
		if err != nil {
			log.Println(err)
			return
		}

	}
	if role == "Serveur" {

		rows, err = server.DATABASE.Query("SELECT * FROM orders WHERE waiter=$1 AND creationdate >= $2 ORDER BY number;", r.FormValue("username"), start)
		if err != nil {
			log.Println(err)
			return
		}

	}

	var order Order

	var situation Situation
	for rows.Next() {

		err = rows.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate)
		if err != nil {
			log.Println(err)
			return

		}

		switch order.Status {
		case "PAYE":
			situation.TotalPaid += order.Totalticket
			situation.CountPaid += 1
			situation.Total += order.Totalticket
		case "IMPAYE":
			situation.TotalInpaid += order.Totalticket
			situation.CountInpaid += 1
			situation.Total += order.Totalticket

		case "ANNULE":
			situation.TotalCanceled += order.Totalticket
			situation.CountICanceled += 1

		}

	}
	w.Header().Set("Content-Type", "application/json")
	jsonresponse, err := json.Marshal(situation)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonresponse)

}
func getSessionsStatus(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}

	server.connectToDB()

	var sessions []Session
	var session Session
	for _, userinfos := range usersinfos {

		rows, err := server.DATABASE.Query("SELECT * FROM sessions WHERE username=$1 ORDER BY sessionid DESC LIMIT 1;", userinfos.Username)
		if err != nil {
			log.Println(err)
			return
		}
		for rows.Next() {
			rows.Scan(&session.Username, &session.SessionID, &session.Start, &session.End)
			sessions = append(sessions, session)
			session = Session{}
		}
	}
	sessionsjson, err := json.Marshal(sessions)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(sessionsjson)

}
func changeSessionStatus(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)

	if username != r.FormValue("username") {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT THE SAME USER"))
		return
	}
	server.connectToDB()

	row := server.DATABASE.QueryRow("SELECT * FROM sessions WHERE username=$1 ORDER BY sessionid DESC LIMIT 1;", r.FormValue("username"))
	var _user string
	var sessionid int64
	var start, end time.Time
	err := row.Scan(&_user, &sessionid, &start, &end)
	if err == sql.ErrNoRows || !end.IsZero() {

		_, err := server.DATABASE.Exec("INSERT INTO sessions (username,start) VALUES ($1,$2);", r.FormValue("username"), time.Now())
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(MakeJsonResponse(w, "INSERT", "Votre session a commencé!", "SESSION_START"))

		return
	}

	_, err = server.DATABASE.Exec("UPDATE sessions SET \"end\"=$1 WHERE sessionid=$2;", time.Now(), sessionid)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "UPDATE", "Votre session a terminé!", "SESSION_END"))

}
func (server *Server) savePrintingSettings(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)

	if username != "Admin" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT_THE_ADMIN"))
		return
	}
	server.connectToDB()

	var count int
	err := server.DATABASE.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		log.Println(err)
		return
	}

	if count == 0 {
		server.DATABASE.Exec("INSERT INTO settings (printingsettings) VALUES($1);", r.FormValue("printingsettings"))
		w.Write(MakeJsonResponse(w, "INSERT", "Les nouveaux paramètres d'impresssion ont été enregistrés", "SAVE_SETTINGS_OK"))
		return
	}

	server.DATABASE.Exec("UPDATE settings SET printingsettings=$1;", r.FormValue("printingsettings"))
	w.Write(MakeJsonResponse(w, "UPDATE", "Les nouveaux paramètres d'impression ont été enregistrés", "UPDATE_SETTINGS_OK"))

}
func saveSecuritySettings(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)

	if username != "Admin" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT_THE_ADMIN"))
		return
	}
	server.connectToDB()

	var count int
	err := server.DATABASE.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		log.Println(err)
		return
	}

	if count == 0 {
		server.DATABASE.Exec("INSERT INTO settings (securitysettings) VALUES($1);", r.FormValue("securitysettings"))
		w.Write(MakeJsonResponse(w, "INSERT", "Les nouveaux paramètres de sécurité ont été enregistrés", "SAVE_SETTINGS_OK"))
		return
	}
	server.DATABASE.Exec("UPDATE settings SET securitysettings=$1;", r.FormValue("securitysettings"))
	w.Write(MakeJsonResponse(w, "UPDATE", "Les paramètres d'impression ont été enregistrés", "UPDATE_SETTINGS_OK"))

}

func saveManagerSettings(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)

	if username != "Admin" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT_THE_ADMIN"))
		return
	}
	server.connectToDB()

	var count int
	err := server.DATABASE.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		log.Println(err)
		return
	}

	if count == 0 {
		server.DATABASE.Exec("INSERT INTO settings (managersettings) VALUES($1);", r.FormValue("managersettings"))
		w.Write(MakeJsonResponse(w, "INSERT", "Les nouveaux paramètres ont été enregistrés", "SAVE_SETTINGS_OK"))
		return
	}
	server.DATABASE.Exec("UPDATE settings SET managersettings=$1;", r.FormValue("managersettings"))
	w.Write(MakeJsonResponse(w, "UPDATE", "Les nouveaux paramètres ont été enregistrés", "UPDATE_SETTINGS_OK"))

}
func notAuthorized(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)

	fmt.Fprintf(w, "Not Authorized. You need a proper browser")
}
func loadSettings(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	server.connectToDB()

	settings = Settings{GeneralSettings: "", ManagerSettings: "", PrintingSettings: "", SecuritySettings: ""}
	if err := server.DATABASE.QueryRow("SELECT * FROM settings;").Scan(&settings.GeneralSettings, &settings.ManagerSettings, &settings.PrintingSettings, &settings.SecuritySettings); err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonresponse, err := json.Marshal(settings)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal([]byte(settings.PrintingSettings), &PS)

	if err != nil {
		log.Println(err)
	}
	w.Write(jsonresponse)

}
func savegeneralsettings(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)

	if username != "Admin" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT_THE_ADMIN"))
		return
	}
	server.connectToDB()

	var count int
	err := server.DATABASE.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		log.Println(err)
		return
	}

	if count == 0 {
		server.DATABASE.Exec("INSERT INTO settings (generalsettings) VALUES($1);", r.FormValue("generalsettings"))
		w.Write(MakeJsonResponse(w, "INSERT", "Les nouveaux paramètres ont été enregistrés", "SAVE_SETTINGS_OK"))
		return
	}
	server.DATABASE.Exec("UPDATE settings SET generalsettings=$1;", r.FormValue("generalsettings"))
	w.Write(MakeJsonResponse(w, "UPDATE", "Les nouveaux paramètres ont été enregistrés", "UPDATE_SETTINGS_OK"))

}
func cancelOrder(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)
	role, err := getUserRole(username)
	if err != nil {
		log.Println(err)
		return
	}
	if role != "Gérant" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT_THE_SAME_WAITER_OR_MANAGER"))
		return
	}
	server.connectToDB()

	_, err = server.DATABASE.Exec("UPDATE orders SET status=$1 WHERE number=$2", "ANNULE", r.FormValue("number"))
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "UPDATE", "Le Ticket est annulé.", "CANCEL_OK"))

}
func checkoutOrder(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NO SESSION"))
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)
	role, err := getUserRole(username)
	if err != nil {
		log.Println(err)
		return
	}
	if username != r.FormValue("waiter") && role != "Gérant" {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", "NOT THE SAME WAITER OR MANAGER"))
		return
	}
	server.connectToDB()

	_, err = server.DATABASE.Exec("UPDATE orders SET status=$1 WHERE number=$2", "PAYE", r.FormValue("number"))
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "UPDATE", "Le Ticket est marqué comme payé.", "CHECKOK"))

}
func PrinterMessenger(jsonData []byte, printerport string, printerip string) {

	resp, err := http.Post("http://"+printerip+":"+printerport+"/api/print", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {

		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

}

func (server *Server) printReceipt(w http.ResponseWriter, r *http.Request) {

	setHeaders(w)

	if !checkValidLogin(r) {
		w.Header().Set("Refresh", "0")
		w.Write(MakeJsonResponse(w, "PRINT", "Not logged in", "PRINTFAILED"))

		return
	}

	for _, profile := range PS {
		if !profile.Active {
			continue
		}
		if profile.PrintingMode == "USB" && profile.PrintingFunction == "Reçu" && r.FormValue("number") != "" {
			number := r.FormValue("number")
			go func(n string) {

				server.connectToDB()

				row := server.DATABASE.QueryRow("SELECT * FROM orders WHERE number=$1;", n)
				var order Order
				err := row.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate)
				if err != nil {
					log.Println(err)
				}

				printedobject := PrintedObject{GeneralSettings: settings.GeneralSettings, Profile: profile, Order: order}
				jsondata, _ := json.Marshal(printedobject)
				PrinterMessenger(jsondata, profile.PrinterPort, profile.PrinterIP)

			}(number)

		}
		if profile.PrintingMode == "LAN" && profile.PrintingFunction == "Reçu" && r.FormValue("number") != "" {
			number := r.FormValue("number")
			go func(n string) {

				server.connectToDB()

				row := server.DATABASE.QueryRow("SELECT * FROM orders WHERE number=$1;", number)
				var order Order
				err := row.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate)
				if err != nil {
					log.Println(err)
				}

				printerIP := profile.PrinterIP
				printerPort := profile.PrinterPort

				address := net.JoinHostPort(printerIP, printerPort)
				conn, err := net.Dial("tcp", address)
				if err != nil {
					log.Printf("Erreur de connexion à l'imprimante: %v\n", err)
					return
				}
				defer conn.Close()

				receiptData := formatReceipt(order, profile)
				if receiptData == nil {
					return
				}
				for i := 1; i <= PS[0].Copies; i++ {
					_, err = conn.Write(receiptData)
					if err != nil {
						log.Printf("Echec d'envoi des données du reçu: %v\n", err)
						return
					}
				}

			}(number)

		}
		if profile.PrintingMode == "USB" && profile.PrintingFunction == "Situation" && r.FormValue("situation") != "" {
			situationtoprint := r.FormValue("situation")

			go func(s string) {
				fmt.Println("Printing situation")
				situationobject := SitutionObject{Situation: s, Profile: profile}
				jsondata, _ := json.Marshal(situationobject)
				PrinterMessenger([]byte(jsondata), profile.PrinterPort, profile.PrinterIP)

			}(situationtoprint)

		}
		if profile.PrintingMode == "LAN" && profile.PrintingFunction == "Situation" && r.FormValue("situation") != "" {

			situationtoprint := r.FormValue("situation")
			go func(s string) {
				networkPrintSituation(s, profile.PrinterIP, profile.PrinterPort, profile.Copies)

			}(situationtoprint)

		}
		if profile.PrintingMode == "USB" && profile.PrintingFunction == "Personnalisé" && r.FormValue("number") != "" {
			number := r.FormValue("number")
			go func(n string) {

				server.connectToDB()

				row := server.DATABASE.QueryRow("SELECT * FROM orders WHERE number=$1;", n)
				var order Order
				err := row.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate)
				if err != nil {
					log.Println(err)
				}

				printedobject := PrintedObject{GeneralSettings: settings.GeneralSettings, Profile: profile, Order: order}
				jsondata, _ := json.Marshal(printedobject)
				PrinterMessenger(jsondata, profile.PrinterPort, profile.PrinterIP)

			}(number)

		}
		if profile.PrintingMode == "LAN" && profile.PrintingFunction == "Personnalisé" && r.FormValue("number") != "" {
			number := r.FormValue("number")
			go func(n string) {

				server.connectToDB()

				row := server.DATABASE.QueryRow("SELECT * FROM orders WHERE number=$1;", number)
				var order Order
				err := row.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate)
				if err != nil {
					log.Println(err)
				}

				printerIP := profile.PrinterIP
				printerPort := profile.PrinterPort

				address := net.JoinHostPort(printerIP, printerPort)
				conn, err := net.Dial("tcp", address)
				if err != nil {
					log.Printf("Erreur de connexion à l'imprimante: %v\n", err)
					return
				}
				defer conn.Close()

				receiptData := formatReceipt(order, profile)
				for i := 1; i <= PS[0].Copies; i++ {
					_, err = conn.Write(receiptData)
					if err != nil {
						log.Printf("Echec d'envoi des données du reçu: %v\n", err)
						return
					}
				}

			}(number)

		}

	}
	w.Write(MakeJsonResponse(w, "PRINT", "Document imprimé.", "PRINTSUCESS"))

}
func (server *Server) getOrders(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	token := r.Header.Get("cookie")[8:]
	var user CustomSession
	user, ok := globalsession[token]
	if !ok {
		log.Println("Non autorisé")
		return
	}
	username := user.Data["user"].(string)
	role, err := getUserRole(username)
	if err != nil {
		log.Println(err)
		return
	}
	server.connectToDB()

	var rows *sql.Rows
	startdate := r.FormValue("startdate")
	starttime := r.FormValue("starttime")
	enddate := r.FormValue("enddate")
	endtime := r.FormValue("endtime")

	if role == "Serveur" {
		rows, err = server.DATABASE.Query("SELECT number, totalticket, waiter, content, creator, status, creationdate, member, memberid FROM orders WHERE waiter=$1 AND (creationdate BETWEEN $2 AND $3) ORDER BY number;", username, startdate+" "+starttime, enddate+" "+endtime)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if role == "Gérant" {
		rows, err = server.DATABASE.Query("SELECT number, totalticket, waiter, content, creator, status, creationdate, member, memberid FROM orders WHERE creationdate BETWEEN $1 AND $2 ORDER BY number;", startdate+" "+starttime, enddate+" "+endtime)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var order Order
	var orders []Order
	for rows.Next() {
		err = rows.Scan(&order.Number, &order.Totalticket, &order.Waiter, &order.Content, &order.Creator, &order.Status, &order.Creationdate, &order.Member, &order.MemberID)
		if err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, order)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonresponse, err := json.Marshal(orders)
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonresponse)
}
func (server *Server) openSession(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	user := r.FormValue("user")
	newuuid := NewUUID()
	biscuit = http.Cookie{
		Name:     "biscuit",
		Value:    newuuid,
		Path:     "/",
		HttpOnly: true,
	}

	globalsession[newuuid] = CustomSession{map[string]interface{}{"user": user, "logintime": time.Now(), "lifetime": time.Now().Add(time.Minute * 10)}}
	http.SetCookie(w, &biscuit)
	w.Write(MakeJsonResponse(w, "LOGIN", "Session created", "OK"))

}
func getUserRole(user string) (string, error) {

	for _, usr := range usersinfos {
		if usr.Username == user {
			return usr.Role, nil

		}
	}
	return "", fmt.Errorf("%s", "pas de rôle trouvé pour cet utilisateur "+user)
}
func (server *Server) endSession(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	biscuit := r.Header.Get("cookie")
	delete(globalsession, biscuit[8:])

	w.Write(MakeJsonResponse(w, "INFO", "Session terminée", "OFF"))

}
func (server *Server) auth(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	user := r.FormValue("user")
	password := r.FormValue("password")

	server.connectToDB()

	var dbpassword string
	row := server.DATABASE.QueryRow("SELECT password FROM users WHERE UPPER(username)=$1;", strings.ToUpper(user))
	err := row.Scan(&dbpassword)
	if err != nil {
		log.Println("select password", err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {

		w.Write(MakeJsonResponse(w, "ERROR", "Echec de l'authentification, veuillez vérifier vos données de connexion.", "NO"))

		return
	}
	newuuid := NewUUID()
	biscuit = http.Cookie{
		Name:     "biscuit",
		Value:    newuuid,
		Path:     "/",
		HttpOnly: true,
	}

	globalsession[newuuid] = CustomSession{map[string]interface{}{"user": user, "logintime": time.Now(), "lifetime": time.Now().Add(time.Minute * 10)}}
	http.SetCookie(w, &biscuit)

	w.Write(MakeJsonResponse(w, "LOGIN", "Authentification réussie.", "OK"))

}
func (server *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var user User

	user.ID = r.FormValue("id")

	server.connectToDB()

	_, err := server.DATABASE.Exec("DELETE FROM users  WHERE id=$1;", user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "DELETE", "Utilisateur Supprimé ", user.Username))

}
func (server *Server) editUser(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var user User
	user.Username = r.FormValue("username")
	user.Role = r.FormValue("role")
	user.Password = r.FormValue("password")
	user.Code = r.FormValue("code")
	user.ID = r.FormValue("id")

	server.connectToDB()

	bcpwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	_, err := server.DATABASE.Exec("UPDATE users SET username=$1, password=$2, role=$3, code=$4 WHERE id=$5;", user.Username, bcpwd, user.Role, user.Code, user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "UPDATE", "Utilisateur modifié ", user.Username))

}
func (server *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	server.connectToDB()

	rows, err := server.DATABASE.Query("SELECT username,role,code,id FROM users;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var userinfo UserInfos
	usersinfos = nil
	for rows.Next() {

		err := rows.Scan(&userinfo.Username, &userinfo.Role, &userinfo.Code, &userinfo.ID)

		if err != nil {
			log.Fatalln(err)
		}

		usersinfos = append(usersinfos, userinfo)

	}
	jsonresp, err := json.Marshal(usersinfos)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(jsonresp)

}
func (server *Server) saveUser(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var user User
	user.Username = r.FormValue("username")
	user.Role = r.FormValue("role")
	user.Password = r.FormValue("password")
	user.Code = r.FormValue("code")

	server.connectToDB()

	var count int32
	row := server.DATABASE.QueryRow("SELECT COUNT(*) FROM users WHERE UPPER(username)=$1;", strings.ToUpper(user.Username))
	if err := row.Scan(&count); err != nil {
		log.Println(err)
		return
	}
	if count > 0 {

		w.Write(MakeJsonResponse(w, "ERROR", "Cet utilisateur existe déjà", user.Username))
		return
	}
	bcpwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	_, err := server.DATABASE.Exec("INSERT INTO users (username,password,role,code) VALUES ($1,$2,$3,$4);", user.Username, bcpwd, user.Role, user.Code)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(MakeJsonResponse(w, "INSERT", "Utilisateur ajouté ", user.Username))

}
func (server *Server) adminLogin(w http.ResponseWriter, r *http.Request) {

	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	pwd := r.FormValue("pwd")

	server.connectToDB()

	row := server.DATABASE.QueryRow("SELECT password FROM admin;")
	var pass string
	err := row.Scan(&pass)
	if err != nil {
		log.Println(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(pwd))
	if err != nil {

		w.Write(MakeJsonResponse(w, "LOGIN", "Mot de passe erroné", "ACCES REFUSE"))
		return
	}
	newuuid := NewUUID()
	biscuit = http.Cookie{
		Name:     "biscuit",
		Value:    newuuid,
		Path:     "/",
		HttpOnly: true,
	}

	globalsession[newuuid] = CustomSession{map[string]interface{}{"user": "Admin", "logintime": time.Now(), "lifetime": time.Now().Add(time.Minute * 10)}}
	http.SetCookie(w, &biscuit)

	w.Write(MakeJsonResponse(w, "LOGIN", "Login granted", "OK"))

}
func checkValidLogin(r *http.Request) bool {
	cookie := r.Header.Get("Cookie")
	if cookie == "" {
		log.Println("no cookie")
		return false
	}
	_, ok := globalsession[cookie[8:]]
	return ok

}
func (server *Server) saveTicket(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if !checkValidLogin(r) {
		w.Write(MakeJsonResponse(w, "ERROR", "Vous n'êtes pas autorisé.", ""))
		return
	}
	server.connectToDB()

	jsonorder := r.FormValue("content")
	totalticket := r.FormValue("totalticket")
	waiter := r.FormValue("waiter")
	creator := r.FormValue("creator")
	member := r.FormValue("member")
	memberid := r.FormValue("memberid")

	var number int
	if err := server.DATABASE.QueryRow("INSERT INTO orders (totalticket,waiter,content,creator,creationdate,status,member,memberid) VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING number;",
		totalticket, waiter, jsonorder, creator, time.Now().Format(time.RFC3339), "IMPAYE", member, memberid).Scan(&number); err != nil {
		log.Println(err)
		return
	}

	w.Write(MakeJsonResponse(w, "INSERT", "Le ticket a été enregistré", strconv.Itoa(number)))

}

var categories []Category

func getCategories(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if categories != nil && !menuchanged {
		w.Header().Set("Content-Type", "application/json")
		jsonresp, err := json.Marshal(categories)
		if err != nil {
			log.Fatalln(err)
		}

		w.Write(jsonresp)
		return
	}
	server.connectToDB()

	rows, err := server.DATABASE.Query("SELECT category_id,category_name FROM categories_table ORDER BY category_id ASC;")
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()
	categories = []Category{}
	var category Category
	for rows.Next() {

		err := rows.Scan(&category.Category_id, &category.Category_name)

		if err != nil {
			log.Fatalln(err)
		}

		categories = append(categories, category)

	}

	w.Header().Set("Content-Type", "application/json")
	jsonresp, err := json.Marshal(categories)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(jsonresp)

}
func (server *Server) saveNewCategory(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.connectToDB()

	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
	}
	if category.Category_name == "" {
		w.Write(MakeJsonResponse(w, "EMPTY", "Le nom de la catégorie ne doit pas être vide.", category.Category_name))
		return
	}
	row := server.DATABASE.QueryRow("SELECT COUNT(*) FROM categories_table WHERE UPPER(category_name) = $1", strings.ToUpper(category.Category_name))
	var count int
	if err := row.Scan(&count); err != nil {
		log.Println(err)
		return
	}
	if count > 0 {
		w.Write(MakeJsonResponse(w, "ERROR", "La catégorie existe déjà", category.Category_name))

		return
	}
	_, err = server.DATABASE.Exec("INSERT INTO categories_table (category_name) VALUES ($1);", category.Category_name)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true
	w.Write(MakeJsonResponse(w, "INSERT", "La catrégorie a été enregistrée", category.Category_name))

}
func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

}
func (server *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.connectToDB()

	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
	}
	_, err = server.DATABASE.Exec("DELETE FROM categories_table  WHERE category_id=$1;", category.Category_id)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = server.DATABASE.Exec("DELETE FROM products_table  WHERE category=$1;", category.Category_id)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true
	w.Write(MakeJsonResponse(w, "DELETE", "La catrégory a été supprimée", strconv.Itoa(category.Category_id)))

}
func (server *Server) editCategory(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.connectToDB()

	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		log.Println(err)
	}

	if category.Category_name == "" {
		w.Write([]byte("EMPTY:Le nom de la catégorie ne peut pas être vide."))
		return
	}
	row := server.DATABASE.QueryRow("SELECT COUNT(*) FROM categories_table WHERE UPPER(category_name) = $1", strings.ToUpper(category.Category_name))
	var count int
	if err := row.Scan(&count); err != nil {
		log.Println(err)
		return
	}
	if count > 0 {

		w.Write([]byte("ERROR:La catégorie " + category.Category_name + " existe déjà."))
		return
	}
	_, err = server.DATABASE.Exec("UPDATE categories_table  SET category_name = $1 WHERE category_id = $2;", category.Category_name, category.Category_id)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true
	w.Write(MakeJsonResponse(w, "EDIT", "La catrégorie a été modifiée", ""))

}
func (server *Server) addNewProduct(w http.ResponseWriter, r *http.Request) {

	setHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
	}
	details := r.FormValue("details")
	var product Product
	err = json.Unmarshal([]byte(details), &product)
	if err != nil {
		log.Println(err)
	}
	server.connectToDB()

	row := server.DATABASE.QueryRow("SELECT COUNT(*) FROM products_table WHERE UPPER(name) = $1 AND category=$2", strings.ToUpper(product.Name), product.Category)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Println(err)
		return
	}
	if count > 0 {
		log.Println("ERROR:Le produit " + product.Name + " existe déjà.")
		w.Write(MakeJsonResponse(w, "ERROR", "Le produit "+product.Name+" existe déjà.", product.Name))
		return
	}
	imagefile, filehandler, err := r.FormFile("imagefile")

	if err != nil {
		log.Println(err)
		return
	}

	defer imagefile.Close()
	ex, _ := os.Executable()
	filename := NewUUID() + filepath.Ext(filehandler.Filename)
	f, err := os.OpenFile(filepath.Dir(ex)+"\\static\\img\\"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}

	io.Copy(f, imagefile)
	_, err = server.DATABASE.Exec("INSERT INTO products_table (name,isfavorite,image,category,price) VALUES ($1,$2,$3,$4,$5);", product.Name, product.Isfavorite, filename, product.Category, product.Price)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true
	w.Write(MakeJsonResponse(w, "INSERT", "Le produit "+product.Name+" a été ajouté.", product.Name))

}

var products []Product
var menuchanged bool

func (server *Server) getMenu(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if products != nil && !menuchanged {
		jsonresp, err := json.Marshal(products)
		if err != nil {
			log.Fatalln(err)
		}

		w.Write(jsonresp)
		return
	}
	server.connectToDB()
	products = []Product{}
	rows, err := server.DATABASE.Query("SELECT * FROM products_table ORDER BY productid ASC;")
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	var product Product
	for rows.Next() {

		err := rows.Scan(&product.Productid, &product.Isfavorite, &product.Name, &product.Category, &product.Price, &product.Image, &product.Subscriberprice, &product.Goldsubscriberprice, &product.SubscriberDiscount, &product.GoldSubscriberDiscount)

		if err != nil {
			log.Fatalln(err)
		}

		products = append(products, product)

	}

	w.Header().Set("Content-Type", "application/json")
	jsonresp, err := json.Marshal(products)
	if err != nil {
		log.Fatalln(err)
	}
	menuchanged = false
	w.Write(jsonresp)

}
func MakeJsonResponse(w http.ResponseWriter, prefix string, message string, data string) []byte {
	w.Header().Set("Content-Type", "application/json")
	jsonresponse, err := json.Marshal(Response{prefix, message, data})
	if err != nil {
		log.Println(err)
	}
	return jsonresponse

}
func (server *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	server.connectToDB()

	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
	}

	_, err = server.DATABASE.Exec("DELETE FROM products_table  WHERE Productid=$1;", product.Productid)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true
	w.Write(MakeJsonResponse(w, "DELETE", "Suppression bien effectuée", strconv.Itoa(product.Productid)))
}
func (server *Server) editProduct(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
	}
	details := r.FormValue("details")
	var product Product

	err = json.Unmarshal([]byte(details), &product)
	if err != nil {
		log.Println(err)
	}
	server.connectToDB()

	_, err = server.DATABASE.Exec("UPDATE products_table SET name=$1, isfavorite=$2, category=$3, price=$4 WHERE productid=$5;", product.Name, product.Isfavorite, product.Category, product.Price, product.Productid)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(MakeJsonResponse(w, "UPDATE", "Le produit "+product.Name+" a été mis à jour.", product.Name))
	menuchanged = true
	imagefile, filehandler, err := r.FormFile("imagefile")

	if err != nil {
		log.Println(err)
		return
	}

	defer imagefile.Close()
	ex, _ := os.Executable()
	filename := NewUUID() + filepath.Ext(filehandler.Filename)
	f, err := os.OpenFile(filepath.Dir(ex)+"\\static\\img\\"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}

	io.Copy(f, imagefile)
	_, err = server.DATABASE.Exec("UPDATE products_table SET image=$1 WHERE productid=$2;", filename, product.Productid)
	if err != nil {
		log.Println(err)
		return
	}
	menuchanged = true

}
func contains(array []int, value int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
func formatReceipt(order Order, profile PrintingProfile) []byte {

	var generalsettings GeneralSettings

	err := json.Unmarshal([]byte(settings.GeneralSettings), &generalsettings)
	if err != nil {
		log.Println(err)

	}
	initPrinter := []byte{0x1B, 0x40}
	selectCharset := []byte{0x1B, 0x52, 0x00}
	alignCenter := []byte{0x1B, 0x61, 0x01}
	alignLeft := []byte{0x1B, 0x61, 0x00}

	boldOn := []byte{0x1B, 0x45, 0x01}
	boldOff := []byte{0x1B, 0x45, 0x00}
	var OrderContent []OrderContent

	err = json.Unmarshal(order.Content, &OrderContent)
	if err != nil {
		log.Println(err)
	}
	personalised := profile.PrintingFunction == "Personnalisé"
	var items []byte

	for _, singlecontent := range OrderContent {

		if personalised {
			print := contains(profile.PrintingOptions, singlecontent.Categoryid)
			if !print {
				continue
			}
			padstr := pad(string(toCP437(singlecontent.Text)), 45, strconv.FormatInt(int64(singlecontent.Count), 10))
			items = append(items, boldOn...)
			items = append(items, []byte(padstr+"\n")...)
			items = append(items, alignLeft...)
			items = append(items, boldOff...)
			if singlecontent.Comment != "" {
				lines := strings.Split(singlecontent.Comment, "\n")
				comment := strings.Join(lines, "\n - ")
				items = append(items, []byte(string(toCP437(" - "+comment+"\n")))...)

			}

			continue
		}
		items = append(items, []byte(spacePAD(string(toCP437(singlecontent.Text)), 30, strconv.FormatInt(int64(singlecontent.Count), 10), 8, strconv.FormatFloat(singlecontent.Price, 'f', 2, 32), 10, strconv.FormatFloat(float64(singlecontent.Price*float64(singlecontent.Count)), 'f', 2, 32))+"\n")...)
	}
	if len(items) == 0 {
		return nil
	}

	doubleHeightOn := []byte{0x1D, 0x21, 0x10}
	doubleHeightOff := []byte{0x1D, 0x21, 0x00}
	cutPaper := []byte{0x1D, 0x56, 0x00}
	separator := []byte("------------------------------------------------\n")
	title := generalsettings.Posname
	headers := "Produit                    Qté    Prix     Total\n"

	if personalised {
		title = profile.Profile
		headers = "Produit                               Quantité\n"
	}
	storeHeader := append(append(append(
		alignCenter,
		doubleHeightOn...),
		[]byte(toCP437(title+"\n"))...),
		doubleHeightOff...,
	)
	if order.Waiter == "" {
		order.Waiter = "        "
	}
	if !personalised {
		storeHeader = append(storeHeader, []byte(toCP437(generalsettings.Address+"\n"))...)
		storeHeader = append(storeHeader, []byte(toCP437(generalsettings.Phone+"\n"))...)

	}
	storeHeader = append(storeHeader, separator...)

	storeHeader = append(storeHeader, alignLeft...)
	storeHeader = append(storeHeader, []byte(toCP437("Ticket n° : "+strconv.FormatInt(int64(order.Number), 10)+"\n"))...)
	storeHeader = append(storeHeader, []byte(toCP437("Serveur : "+order.Waiter+"    "))...)

	storeHeader = append(storeHeader, []byte(toCP437("Date : "+getVFDate(order.Creationdate)+"\n"))...)
	storeHeader = append(storeHeader, separator...)

	itemsHeader := append(append(alignLeft,
		[]byte(toCP437(headers))...),
		separator...,
	)

	totals := []byte(spaceSoloPAD("Total :", 48, strconv.FormatFloat(float64(order.Totalticket), 'f', 2, 32)+" MAD\n"))

	footer := append(append(append(
		alignLeft,
		boldOff...),
		[]byte(toCP437("Merci de votre visite\nA bientôt!\n"))...),
		boldOff...,
	)

	receipt := append(initPrinter, selectCharset...)
	receipt = append(receipt, storeHeader...)
	receipt = append(receipt, itemsHeader...)
	receipt = append(receipt, items...)
	receipt = append(receipt, separator...)
	if !personalised {
		receipt = append(receipt, boldOn...)
		receipt = append(receipt, totals...)
		receipt = append(receipt, doubleHeightOff...)
		receipt = append(receipt, boldOff...)
		receipt = append(receipt, separator...)
		receipt = append(receipt, footer...)
		receipt = append(receipt, alignCenter...)
		receipt = append(receipt, []byte("Code WIFI : "+generalsettings.Wificode+"\n")...)

	}
	receipt = append(receipt, []byte("\n\n\n\n\n\n")...)
	receipt = append(receipt, cutPaper...)

	return receipt
}
func toCP437(input string) []byte {
	encoder := charmap.CodePage437.NewEncoder()
	convertedText, err := encoder.String(input)
	if err != nil {
		log.Println("Encoding conversion error:", err)
		return []byte(input)
	}

	return []byte(convertedText)
}
func pad(product string, lenght1 int, quantity string) string {
	var result = product
	for len(result+quantity) < lenght1 {
		result += " "

	}
	result += quantity
	return result

}
func spacePAD(product string, lenght1 int, quantity string, lenght2 int, price string, lenght3 int, total string) string {
	var result = product
	for len(result+quantity) < lenght1 {
		result += " "

	}
	result += quantity
	for len(result+price) < lenght1+lenght2 {
		result += " "
	}
	result += price
	for len(result+total) < lenght1+lenght2+lenght3 {
		result += " "
	}
	result += total
	return result
}
func getVFDate(dt time.Time) string {
	enDateStr := dt.String()

	layout := "2006-01-02 15:04:05 -0700 -0700"
	parsedTime, err := time.Parse(layout, enDateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
		return ""
	}

	frDate := parsedTime.Format("02/01/2006 15:04:05")

	return frDate
}
func spaceSoloPAD(text1 string, lenght int, text2 string) string {
	var result = text1
	for len(result+text2) < lenght {
		result += " "

	}

	result += text2
	return result
}
