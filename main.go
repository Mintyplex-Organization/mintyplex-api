package main

import (
	"mintyplex-api/internal/database"
	"mintyplex-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func main(){
	engine := fiber.New()

	middleware.CorsMiddleware(engine)

	engine.Get("/plexer", func (ctx *fiber.Ctx) error {
		return ctx.SendString("Plexer Baby")
	})

	DB := database.MongoClient()
	engine.Use(middleware.IngestDb(DB))

	

	engine.Listen(":8081")


}