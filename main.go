package main

import (
	"fmt"

	"github.com/chagaphongk/register-api/constant"
	"github.com/chagaphongk/register-api/service"
	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.ERROR)

	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(constant.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/login" || c.Path() == "/signup" {
				return true
			}
			return false
		},
	}))

	// Database connection
	connString := fmt.Sprintf("%v://%v:%v/%v", constant.MongoUser, constant.MongoHost, constant.MongoPort, constant.DBCollection)
	db, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Create indices
	if err = db.Copy().DB(constant.DBName).C(constant.DBCollection).EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	h := &service.Handler{DB: db}

	// Routes
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
