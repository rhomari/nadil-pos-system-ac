//go:build windows
// +build windows

package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func openLogfile() *os.File {

	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	logFilePath := filepath.Join(exeDir, "server.log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)

	}
	return file
}
func setupMode() {

	if strings.ToLower(os.Args[1]) == "setup" {
		fmt.Println("Vous avez demander la configuration du serveur. Choisisez parmi les choix suivants :")
		fmt.Println("1) Activation du logiciel")
		fmt.Println("2) Changer le mot de passe administrateur")

		var choice string
		var key string

		fmt.Scanf("%s", &choice)
		switch choice {
		case "1":
			fmt.Print("Entrez la clé d'activation : ")
			fmt.Scan(&key)

			var macs string
			as, err := getMacAddr()

			if err != nil {
				fmt.Println(err)
			}

			macs = strings.Join(as, "|")

			server.connectToDB()
			r, err := server.DATABASE.Exec("UPDATE setup SET macs=$1 WHERE key=$2;", &macs, &key)
			if err != nil {
				fmt.Println(err)
				return
			}
			count, err := r.RowsAffected()
			if err != nil {
				fmt.Println(err)
				return
			}
			if count == 0 {
				fmt.Println("Clé d'activation erronée.")
				os.Exit(-1)
				return
			}
			fmt.Println("Logiciel activé. Veuillez relancer l'application")
			os.Exit(0)

		case "2":
			server.setupAdminPassword()
			os.Exit(0)

		default:
			fmt.Println("Choix inéxistant: ")
			os.Exit(0)

		}

	}

}
func promptPassword() string {
	var newpassword string
	var confirmnewpassword string
	fmt.Println("Entrez le nouveau mot de passe :")
	fmt.Scan(&newpassword)
	fmt.Println("Confirmez le nouveau mot de passe :")
	fmt.Scan(&confirmnewpassword)
	for newpassword != confirmnewpassword {
		fmt.Println("Le mot de passe et sa confirmation ne sont pas identiques")
		fmt.Println("Entrez le nouveau mot de passe :")
		fmt.Scan(&newpassword)
		fmt.Println("Confirmez le nouveau mot de passe :")
		fmt.Scan(&confirmnewpassword)

	}
	return newpassword
}

func (server *Server) setupAdminPassword() {

	var password string
	server.connectToDB()
	if err := server.DATABASE.QueryRow("SELECT password FROM admin").Scan(&password); err != nil {
		fmt.Println(err)
	}

	if password == "" {
		newpassword := promptPassword()
		bytespass, err := bcrypt.GenerateFromPassword([]byte(newpassword), 14)
		if err != nil {
			fmt.Println(err)
		}
		_, err = server.DATABASE.Exec("INSERT INTO admin  (password) VALUES($1);", string(bytespass))
		if err != nil {
			fmt.Println(err)

		}

		return

	}
	var oldpassword string
	fmt.Println("Entrez l'ancien mot de passe :")
	fmt.Scan(&oldpassword)

	if bcrypt.CompareHashAndPassword([]byte(password), []byte(oldpassword)) != nil {
		fmt.Println("Ce mot de passe ne correspond pas à l'ancien.")
		return

	}
	newpassword := promptPassword()
	bytesnewpassword, _ := bcrypt.GenerateFromPassword([]byte(newpassword), 14)
	_, err := server.DATABASE.Exec("UPDATE admin set password=$1 WHERE password=$2 ;", string(bytesnewpassword), password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Le mot de passe administrateur a été mis à jour.")

}
func (server *Server) checkLicence() bool {
	server.connectToDB()

	var macs string
	if err := server.DATABASE.QueryRow("SELECT macs FROM setup LIMIT 1;").Scan(&macs); err != nil {
		fmt.Println(err)
		return false
	}
	if macs == "" {
		fmt.Println("Vous n'avez pas de licence pour lancer ce logiciel.\nVeuillez contacter Mr Tariq Rhomari sur l'adresse mail tariqrhomari@gmail.com ou par téléphone +212671626226.")
		return false
	}

	as, err := getMacAddr()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(as)
	for _, inter := range as {
		if strings.Contains(macs, inter) || strings.Contains(macs, "OK") {
			return true
		}
	}

	fmt.Println("Votre licence n'est pas valide.\nVeuillez contacter Mr Tariq Rhomari sur l'adresse mail tariqrhomari@gmail.com ou par téléphone +212671626226.")
	return false

}
func main() {

	if !server.checkLicence() && len(os.Args) == 1 {
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		setupMode()
		return
	}

	port := os.Getenv("NADIL_SERVER_PORT")

	setupRoutes()
	fmt.Println("Starting service on port", port)
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		fmt.Println(err)
	}

}
