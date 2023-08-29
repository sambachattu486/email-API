package main

import (
	"bytes"
	"encoding/base64"
	io "io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// type Service struct {
// 	Port   string
// 	Router *fiber.App
// }

type Email struct {
	router *fiber.App
}

func loadConfig() {
	// Config is stored as Base64 encoded string. Decoding config

	cfg, err := io.ReadFile(".config")
	if err != nil {
		panic(err)
	}
	cfgData, err := base64.StdEncoding.DecodeString(string(cfg))
	if err != nil {
		panic(err)
	}
	// Read the config data from the decoded json string
	viper.SetConfigType("json")
	if err = viper.ReadConfig(bytes.NewReader(cfgData)); err != nil {
		panic(err)
	}
}

func main() {

	loadConfig()

	app := fiber.New()
	str := &Email{router: app}
	
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE, UPDATE",
		AllowHeaders:     "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	}))

	application := EmailService{
		service: *str,
	}

	application.Bootstrap()

	app.Listen(viper.GetString("service.port"))
}
