package routes

import (
	"main/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())

	// Productos
	notes := api.Group("/notes")
	notes.Get("/", handlers.GetAllNotes)
	notes.Post("/", handlers.NewNote)
	notes.Put("/:id", handlers.EditNote)
	notes.Delete("/:id", handlers.DeleteNote)
}
