package main

import (
	"encoding/gob"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/metaclips/LetsTalk/backend/controller"
	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

func main() {
	if err := values.LoadConfiguration("./config.json"); err != nil {
		log.Fatalln("unable to load config", err)
	}

	gob.Register(time.Time{})
	model.InitDB()
	go model.HubConstruct.Run()

	router := httprouter.New()

	router.GET("/ws", controller.ServeWs)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/register", controller.RegisterUserPost)

	port := values.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	router.ServeFiles("/assets/*filepath", http.Dir("./views/assets"))
	log.Infoln("Webserver UP")

	handler := cors.New(cors.Options{
		Debug:          true,
		AllowedOrigins: values.Config.CorsAllowedOrigins,
	}).Handler(router)

	// Optional use of TLS due to Heroku serving TLS at low level.
	if values.Config.TLS.CertPath != "" && values.Config.TLS.KeyPath != "" {
		if err := http.ListenAndServeTLS(":"+port, values.Config.TLS.CertPath, values.Config.TLS.KeyPath, handler); err != nil {
			log.Fatalln(err)
		}

		return
	}

	// Note: without HTTPS users wont be able to login as SetCookie uses Secure flag.
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalln(err)
	}
}
