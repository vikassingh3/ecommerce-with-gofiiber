package controllers

import (
	"os"

	"github.com/vikas/config"
	"github.com/vikas/model"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
)

func Product(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION"))

	data := new(model.Product)

	err := c.BodyParser(&data)

	// if error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	result, err := productCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert product",
			"error":   err,
		})
	}

	// get the inserted data
	product := &model.Product{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	productCollection.FindOne(c.Context(), query).Decode(product)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"product": product,
		},
	})
}
