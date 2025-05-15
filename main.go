package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Logger middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		fmt.Printf("[%s] %s %s - %v\n",
			c.Method(),
			c.Path(),
			c.IP(),
			duration,
		)
		return err
	}
}

// Message represents a simple message structure
type Message struct {
	Content string `json:"content" validate:"required,min=1,max=100"`
	Author  string `json:"author" validate:"required,min=1,max=50"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {
	// Create new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(ErrorResponse{
				Error:   "Request Error",
				Message: err.Error(),
			})
		},
	})

	// Use logger middleware
	app.Use(Logger())

	// Use recover middleware
	app.Use(recover.New())

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Add a simple API endpoint
	app.Get("/api/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app_name":  "Simple API",
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
			"status":    "running",
		})
	})

	// Add message endpoint
	app.Post("/api/message", func(c *fiber.Ctx) error {
		message := new(Message)

		// Parse request body
		if err := c.BodyParser(message); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		// Validate message
		if message.Content == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Message content is required")
		}
		if message.Author == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Author is required")
		}
		if len(message.Content) > 100 {
			return fiber.NewError(fiber.StatusBadRequest, "Message content too long")
		}
		if len(message.Author) > 50 {
			return fiber.NewError(fiber.StatusBadRequest, "Author name too long")
		}

		// Return success response
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Message received successfully",
			"data":      message,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
