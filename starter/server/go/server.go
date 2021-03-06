package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/foolin/goview/supports/echoview"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

func main() {
	err := godotenv.Load("./.env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = echoview.Default()

	e.Static("/", os.Getenv("STATIC_DIR"))
	e.File("/", os.Getenv("STATIC_DIR")+"/index.html")
	e.POST("/webhook", webhookHandler)
	e.Logger.Fatal(e.Start(":4242"))
}

func webhookHandler(c echo.Context) (err error) {
	request := c.Request()
	payload, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	var event stripe.Event

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret != "" {
		event, err = webhook.ConstructEvent(payload, request.Header.Get("Stripe-Signature"), webhookSecret)
		if err != nil {
			return err
		}
	} else {
		err := json.Unmarshal(payload, &event)
		if err != nil {
			return err
		}
	}

	objectType := event.Data.Object["object"].(string)
	switch objectType {
	case "customer.created":
		// Handle logic when a customer is created
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, event)
}
