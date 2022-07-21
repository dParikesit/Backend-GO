package main

import (
	"github.com/dParikesit/bnmo-backend/handlers"
	middleware2 "github.com/dParikesit/bnmo-backend/middleware"
	"github.com/dParikesit/bnmo-backend/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	jwtconfig := middleware.JWTConfig{Claims: &utils.CustomClaims{}, SigningKey: []byte(os.Getenv("SECRET"))}

	auth := e.Group("/auth")
	auth.POST("/register", handlers.Registration)
	auth.POST("/login", handlers.Login)

	admin := e.Group("/admin")
	admin.Use(middleware.JWTWithConfig(jwtconfig))
	admin.Use(middleware2.CheckUser)
	admin.Use(middleware2.CheckAdmin)
	admin.POST("/verify/user", handlers.VerifyUser)
	admin.POST("/verify/req", handlers.VerifyRequest)
	admin.GET("/user", handlers.GetUsers)

	customer := e.Group("/customer")
	customer.Use(middleware.JWTWithConfig(jwtconfig))
	customer.Use(middleware2.CheckUser)
	customer.Use(middleware2.CheckCustomer)
	customer.GET("/info", handlers.GetInfo)
	customer.GET("/history/:historyType", handlers.GetHistory)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}
