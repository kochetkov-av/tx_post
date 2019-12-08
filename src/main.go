package main

import (
	"io/ioutil"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

type CustomContext struct {
	echo.Context
	Validator *validator.Validate
	Config    *Config
	Db        *gorm.DB
	DbMut     sync.Mutex
	Lg        echo.Logger
}

func main() {
	e := echo.New()

	// Load configuration file.
	configBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		e.Logger.Fatal(err)
	}
	var config *Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		e.Logger.Fatal(err)
	}

	db, err := initDb(config)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Mutex used to lock parallel access to DB, to exclude concurrency problems.
	dbMut := sync.Mutex{}

	// Register custom context.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context:   c,
				Validator: validator.New(),
				Config:    config,
				Db:        db,
				DbMut:     dbMut,
			}
			return next(cc)
		}
	})

	e.POST("/tx", TxPostHandler)

	go TxCancelation(&CustomContext{Db: db, DbMut: dbMut, Config: config, Lg: e.Logger})

	e.Logger.Fatal(e.Start(":4000"))
}
