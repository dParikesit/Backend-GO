package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/dParikesit/bnmo-backend/controllers"
	"github.com/dParikesit/bnmo-backend/models"
	"github.com/dParikesit/bnmo-backend/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Login(c echo.Context) error {
	var user models.User
	var err error

	if err = c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var dbUser models.User
	dbUser, err = controllers.UserGetByUsername(user.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	expiry := time.Now().Add(time.Hour * 24)

	claims := utils.CustomClaims{
		Username:   dbUser.Username,
		Name:       dbUser.Name,
		IsAdmin:    dbUser.IsAdmin,
		IsVerified: dbUser.IsVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return err
	}

	var cookie *http.Cookie
	cookie.Name = "access_token"
	cookie.Value = tokenSigned
	cookie.Expires = expiry
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

func Registration(c echo.Context) error {
	var user models.User

	user.Username = c.FormValue("username")
	user.Name = c.FormValue("name")

	hashed, err := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), 14)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("hashing error"))
	}

	user.Password = string(hashed)

	photo, err := c.FormFile("photo")
	src, err := photo.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	dst, err := os.Create(filepath.Join("files", filepath.Base(base64.URLEncoding.EncodeToString([]byte(photo.Filename)))))
	if err != nil {
		return err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
