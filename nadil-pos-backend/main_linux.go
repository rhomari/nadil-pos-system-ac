//go:build linux
// +build linux

package main

import (
	"fmt"
	"net/http"
	"strings"

	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func setupMode() {

	if strings.ToLower(os.Args[1]) == "setup" {
		log.Println("Vous avez demander la configuration du serveur. Choisisez parmi les choix suivants :")
		log.Println("1) Activation du logiciel")
		log.Println("2) Changer le mot de passe administrateur")

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
				log.Println(err)
			}

			macs = strings.Join(as, "|")

			server.connectToDB()
			r, err := server.DATABASE.Exec("UPDATE setup SET macs=$1 WHERE key=$2;", &macs, &key)
			if err != nil {
				log.Println(err)
				return
			}
			count, err := r.RowsAffected()
			if err != nil {
				log.Println(err)
				return
			}
			if count == 0 {
				log.Println("Clé d'activation erronée.")
				os.Exit(-1)
				return
			}
			log.Println("Logiciel activé. Veuillez relancer l'application")
			os.Exit(0)

		case "2":
			server.setupAdminPassword()
			os.Exit(0)

		default:
			log.Println("Choix inéxistant: ")
			os.Exit(0)

		}

	}

}
func promptPassword() string {
	var newpassword string
	var confirmnewpassword string
	log.Println("Entrez le nouveau mot de passe :")
	fmt.Scan(&newpassword)
	log.Println("Confirmez le nouveau mot de passe :")
	fmt.Scan(&confirmnewpassword)
	for newpassword != confirmnewpassword {
		log.Println("Le mot de passe et sa confirmation ne sont pas identiques")
		log.Println("Entrez le nouveau mot de passe :")
		fmt.Scan(&newpassword)
		log.Println("Confirmez le nouveau mot de passe :")
		fmt.Scan(&confirmnewpassword)

	}
	return newpassword
}

func (server *Server) setupAdminPassword() {

	var password string
	server.connectToDB()
	if err := server.DATABASE.QueryRow("SELECT password FROM admin").Scan(&password); err != nil {
		log.Println(err)
	}

	if password == "" {
		newpassword := promptPassword()
		bytespass, err := bcrypt.GenerateFromPassword([]byte(newpassword), 14)
		if err != nil {
			log.Println(err)
		}
		_, err = server.DATABASE.Exec("INSERT INTO admin  (password) VALUES($1);", string(bytespass))
		if err != nil {
			log.Println(err)

		}

		return

	}
	var oldpassword string
	log.Println("Entrez l'ancien mot de passe :")
	fmt.Scan(&oldpassword)

	if bcrypt.CompareHashAndPassword([]byte(password), []byte(oldpassword)) != nil {
		log.Println("Ce mot de passe ne correspond pas à l'ancien.")
		return

	}
	newpassword := promptPassword()
	bytesnewpassword, _ := bcrypt.GenerateFromPassword([]byte(newpassword), 14)
	_, err := server.DATABASE.Exec("UPDATE admin set password=$1 WHERE password=$2 ;", string(bytesnewpassword), password)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Le mot de passe administrateur a été mis à jour.")

}
func (server *Server) checkLicence() bool {
	server.connectToDB()

	var macs string
	if err := server.DATABASE.QueryRow("SELECT macs FROM setup LIMIT 1;").Scan(&macs); err != nil {
		log.Println(err)
		return false
	}
	if macs == "" {
		log.Println("Vous n'avez pas de licence pour lancer ce logiciel.\nVeuillez contacter Mr Tariq Rhomari sur l'adresse mail tariqrhomari@gmail.com ou par téléphone +212671626226.")
		return false
	}

	as, err := getMacAddr()

	if err != nil {
		log.Println(err)
	}
	log.Println(as)
	for _, inter := range as {
		if strings.Contains(macs, inter) || strings.Contains(macs, "OK") {
			return true
		}
	}

	log.Println("Votre licence n'est pas valide.\nVeuillez contacter Mr Tariq Rhomari sur l'adresse mail tariqrhomari@gmail.com ou par téléphone +212671626226.")
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
	log.Println("Starting service on port", port)
	if err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil); err != nil {
		log.Println(err)
	}

}
