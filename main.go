package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/dParikesit/bnmo-backend/utils"
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

	err = utils.InitDB()
	if err != nil {
		log.Fatalln("DB Connection error")
	}

	fmt.Println("Db Connected")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}
