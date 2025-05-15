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

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Create new Fiber app
	app := fiber.New()

	// Use logger middleware
	app.Use(Logger())

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

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

	// Add create user endpoint
	app.Post("/api/user", func(c *fiber.Ctx) error {
		user := new(User)

		// Parse request body
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate required fields
		if user.Name == "" || user.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Name and email are required",
			})
		}

		// Set creation time
		user.CreatedAt = time.Now()
		user.ID = fmt.Sprintf("user_%d", time.Now().Unix())

		return c.Status(201).JSON(user)
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
