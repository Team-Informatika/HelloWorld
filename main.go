package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

func main() {
	// Create new Fiber app
	app := fiber.New()

	// Use logger middleware
	app.Use(Logger())

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

	// Start server
	log.Fatal(app.Listen(":3000"))
}
