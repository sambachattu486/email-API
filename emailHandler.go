package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	service Email
}

func (email EmailService) Bootstrap() {
	appRoute := email.service.router.Group("/api")
	appRoute.Post("/notify", emailSender)
}

func  emailSender(ctx *fiber.Ctx) error {
	var requestBody map[string]interface{}
	if err := ctx.BodyParser(&requestBody); err != nil {
		return err
	}
	currentTime := time.Now()
	output := currentTime.String()

	// read emailTemplate html for email subject
	file, err := os.Open("emailTemplate.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	emailSubject := doc.Find("h1").Text()
	

	var customMessage string
	for key, val := range requestBody {
		if (key) == "name" || key == "Name" {
			continue
		}
		customMessage += fmt.Sprint(strings.Title(fmt.Sprint("\n", key)), (fmt.Sprint(" : ", val)))
	}
	// parse the emailTemplate to insert dynamic data
	t, err := template.ParseFiles("emailTemplate.html")
	if err != nil {
		panic(err)
	}
	messageMap := make(map[string]string)
	messageMap["name"] = fmt.Sprint(requestBody["name"])
	messageMap["dynamicContent"] = customMessage
	messageMap["time"] = output[11:16]
	messageMap["date"] = output[:10]

	// execute the html template with the message
	var htmlTemplate bytes.Buffer
	err = t.Execute(&htmlTemplate, messageMap)
	if err != nil {
		panic(err)
	}
	//gomail service
	mail := gomail.NewMessage()
	mail.SetHeader("From", viper.GetString("gomail.email"))
	mail.SetHeader("To", fmt.Sprint(requestBody["email"]))
	mail.SetHeader("Subject", emailSubject)
	mail.SetBody("text/html", htmlTemplate.String())

	emailDialer := gomail.NewDialer(viper.GetString("gomail.serviceName"), viper.GetInt("gomail.port"), viper.GetString("gomail.email"), viper.GetString("gomail.password"))
	if err := emailDialer.DialAndSend(mail); err != nil {
		log.Print(err)
		return ctx.JSON(map[string]string{"Email not sent": fmt.Sprint(err)})
	}
	return ctx.JSON(map[string]string{"message": "Email Sent Successfully"})
}
