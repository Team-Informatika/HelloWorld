package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
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

	// Define a route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Add a new API endpoint
	app.Get("/api/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "active",
			"message": "Server is running",
			"version": "1.0.0",
		})
	})

	// Add user info endpoint
	app.Get("/api/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		return c.JSON(fiber.Map{
			"user_id":    userID,
			"name":       "User " + userID,
			"created_at": time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
