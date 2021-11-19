package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikas/auth"
	"github.com/vikas/controllers"
)

func Route(app *fiber.App) {
	app.Post("api/product", controllers.Product)
	app.Get("/api/cart", controllers.GetCartInfo)
	app.Delete("/api/cart/{id}", controllers.DeleteItemFromCart)
	app.Delete("/api/cart", controllers.ResetCart)
	app.Post("/api/cart/:id", controllers.AddItemToCart)
	app.Post("/api/order", controllers.PlaceOrders)
	app.Post("/login", auth.Login)
	app.Post("/signup", auth.SignUp)
	app.Post("/logout/:email", auth.LogOut)
}
