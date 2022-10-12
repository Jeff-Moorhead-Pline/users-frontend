package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	UserTypeAdmin     string = "admin"
	UserTypeClient           = "client"
	UserTypeAssociate        = "associate"
	UserTypeReadOnly         = "read-only"
)

type User struct {
	ID        string `json:"id" form:"id"`
	UserSince string `json:"userSince" form:"userSince"`
	Type      string `json:"type" form:"type"`
}

var users = []User{
	{
		ID:        "jmoorhead",
		UserSince: time.Date(2021, time.June, 21, 23, 29, 32, 0, time.UTC).Format(time.Stamp),
		Type:      UserTypeAdmin,
	},
	{
		ID:        "jmcarthur",
		UserSince: time.Date(2020, time.May, 30, 8, 43, 23, 0, time.UTC).Format(time.Stamp),
		Type:      UserTypeAssociate,
	},
	{
		ID:        "asalatto",
		UserSince: time.Now().Format(time.Stamp),
		Type:      UserTypeClient,
	},
	{
		ID:        "amcbride",
		UserSince: time.Now().Format(time.Stamp),
		Type:      UserTypeReadOnly,
	},
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/api/users", handleGetUsers)
	e.GET("/api/users/:id", handleGetUserByID)
	e.POST("/api/users/new", handleCreateUser)
	e.PUT("/api/users/:id", handleUpdateUser)
	e.DELETE("/api/users/:id", handleDeleteUser)

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleGetUsers(c echo.Context) error {

	return c.JSON(http.StatusOK, &users)
}

func handleGetUserByID(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no user id given")
	}

	for _, u := range users {
		if u.ID == id {
			return c.JSON(http.StatusOK, &u)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "user not found")
}

func handleCreateUser(c echo.Context) error {

	newUser := new(User)
	if err := c.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser.UserSince = time.Now().Format(time.Stamp)

	for _, u := range users {
		if u.ID == newUser.ID {
			return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("user %v already exists", newUser.ID))
		}
	}

	users = append(users, *newUser)

	return c.JSON(http.StatusCreated, newUser)
}

func handleUpdateUser(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no user id given")
	}

	userData := new(User)
	if err := c.Bind(userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userData.ID = id

	for i, u := range users {
		if u.ID == id {
			users[i] = *userData
			return c.JSON(http.StatusOK, userData)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %v not found", id))
}

func handleDeleteUser(c echo.Context) error {

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "no user id given")
	}

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("user %v not found", id))
}
