package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
)

type LogEntry struct {
	Time      string `json:"time"`
	IP        string `json:"ip"`
	Status    int    `json:"status"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	RequestID string `json:"request_id"` // Added field for request ID
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ServerHeader:  "X-go-server",
		AppName:       "GO-RestAPI",
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})
	app.Use(recover.New())

	app.Use(requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: fiber.HeaderXRequestID,
		Generator:  utils.UUID,
	}))
	app.Use(healthcheck.New())

	app.Use(func(c *fiber.Ctx) error {
		// Process request
		if err := c.Next(); err != nil {
			return err
		}

		logEntry := LogEntry{
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			IP:        c.IP(),
			Status:    c.Response().StatusCode(),
			Method:    c.Method(),
			Path:      c.Path(),
			RequestID: c.GetRespHeader(fiber.HeaderXRequestID),
		}

		// Convert log entry to JSON
		logJSON, _ := json.Marshal(logEntry)
		os.Stdout.Write(append(logJSON, '\n'))

		return nil
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	port := "8000"

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// check if certs folder exists
	if _, err := os.Stat("certs"); os.IsNotExist(err) {
		log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
	} else {
		// Load SSL certificate
		cer, err := tls.LoadX509KeyPair("certs/ssl.cert", "certs/ssl.key")
		if err != nil {
			log.Fatal(err)
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		ln, err := tls.Listen("tcp", fmt.Sprintf(":%s", port), config)
		if err != nil {
			panic(err)
		}

		log.Fatal(app.Listener(ln))
	}

	<-idleConnsClosed

	fmt.Println("Running cleanup tasks...")

}
