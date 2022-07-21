package main

import (
	"github.com/dParikesit/bnmo-backend/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalln("Env load failed")
	}

	if err = utils.Db.InitDB(); err != nil {
		log.Fatalln("DB Connection error")
	}

	if err = utils.Db.InitSeeding(); err != nil {
		log.Fatalln("Seeding error")
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}
