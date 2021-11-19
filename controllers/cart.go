package controllers

import (
	"fmt"
	"os"

	"github.com/vikas/config"

	"github.com/gofiber/fiber/v2"
	"github.com/vikas/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//GetCartInfo func will return all item of the cart with total price
func GetCartInfo(c *fiber.Ctx) error {
	cartCollection := config.MI.DB.Collection(os.Getenv("ADD_TO_CART_COLLECTION"))

	//declaring variable in which data will be stored
	var cartData []model.Item
	var totalPrice float64

	//selecting DB collection/table & taking to a variable

	//finding all documents/rows from DB
	cursor, err := cartCollection.Find(c.Context(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	//getting multiple documents(rows)
	//Iterating through the cursor allows us to decode one document at a time
	for cursor.Next(c.Context()) {
		//creating a temporary variable in which the single document can be decoded
		temp := new(model.Item)
		err := c.BodyParser(&temp)
		if err != nil {
			fmt.Println("can not decode model.temp")
		}
		totalPrice += temp.Total //price of each item will be added together and stored to totalPrice variable

		cartData = append(cartData, *temp) //finally taking this single document to the slice of documents
	}

	//preparing data for json response
	returnData := struct {
		Items []model.Item `json:"items"`
		Total float64      `json:"total"`
	}{
		Items: cartData,
		Total: totalPrice,
	}

	//encoding data to json
	/* 	rData, err := json.Marshal(returnData)
	   	if err != nil {
	   		fmt.Println(err)
	   	}

	   	// c.Request().Header.Set("Content-Type", "application/json") //setting content type as application/json
	   	c.Write(rData) //finally response back to client */
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": returnData,
	})

}

//AddItemToCart func will insert an item to the cart and return the inserted id
func AddItemToCart(c *fiber.Ctx) error {
	//declaring variable in which data will be stored
	var item model.Item

	//receiving request body

	c.BodyParser(&item) //getting json data to the variable

	//calculating total price of this item based on price and unit
	item.Total = item.Price * float64(item.Unit)

	//connecting to DB
	cartCollection := config.MI.DB.Collection(os.Getenv("ADD_TO_CART_COLLECTION"))
	/* 	productCollection := config.MI.DB.Collection(os.Getenv("PRODUCT_COLLECTION"))

	   	paramsId := c.Params("id")
	   	id, err := primitive.ObjectIDFromHex(paramsId)
	   	if err != nil {
	   		fmt.Println("can not parse Id")
	   	}
	   	product := &model.Product{}

	   	query := bson.D{{Key: "_id", Value: id}}

	   	err = productCollection.FindOne(c.Context(), query).Decode(product)

	   	if err != nil {
	   		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	   			"success": false,
	   			"message": "product not found",
	   			"error":   err,
	   		})
	   	} */

	//setting up options for insert/update
	opts := options.Update().SetUpsert(true)
	update := bson.M{
		"$set": bson.M{
			"unit":  item.Unit,
			"price": item.Price,
			"total": item.Total,
		},
	}

	//inserting item (document/row) to DB, if already exist then modify
	res, err := cartCollection.UpdateOne(c.Context(), bson.M{"code": item.Code}, update, opts)
	if err != nil {
		fmt.Println("unable to update data inside cart")
	}

	//preparing data for json response
	var returnData model.ResponseData

	if res.MatchedCount == 0 { //new entry for this item
		//getting inserted id as string from the result of insert query to DB
		upsertedID := res.UpsertedID.(primitive.ObjectID).Hex() //type assertion && Calling Hex func

		returnData.Status = "success"
		returnData.ID = upsertedID
		returnData.Message = "Item successfully added to cart. Inserted id: " + upsertedID
	} else { //exsisting item, update only
		returnData.Status = "success"
		returnData.Message = "Existing item. Item successfully updated in the cart."
	}

	//encoding data to json
	/* rData, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println(err)
	}

	// w.Header().Set("Content-Type", "application/json") //setting content type as application/json
	c.Write(rData) //finally response back to client */
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": returnData,
	})

}

//DeleteItemFromCart func will delete an item from the cart and return the number of deleted item with the id of deleted item
func DeleteItemFromCart(c *fiber.Ctx) error {
	//getting item id from request parameter

	paramId := c.Params("id")

	//converting requested id from string to primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		fmt.Println(err)
	}

	//connecting to DB
	cartCollection := config.MI.DB.Collection(os.Getenv("ADD_TO_CART_COLLECTION"))

	//deleting specific item (document/row) from DB
	res, err := cartCollection.DeleteOne(c.Context(), bson.M{"_id": id})
	if err != nil {
		fmt.Println(err)
	}

	//preparing data for json response
	returnData := new(model.ResponseData)

	if res.DeletedCount == 0 { //if no item found with the provided id
		returnData.Status = "error"
		returnData.Message = "No product found with ID: " + paramId
	} else { //if item found
		returnData.Status = "success"
		returnData.ID = paramId
		returnData.Message = "Item successfully deleted from cart. Total item deleted: " + fmt.Sprintf("%d", res.DeletedCount)
	}

	//encoding data to json
	/* rData, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println(err)
	}

	//w.Header().Set("Content-Type", "application/json") //setting content type as application/json
	c.Write(rData) //finally response back to client */
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": returnData,
	})

}

//ResetCart func will delete all item(s) from the cart and return the number of deleted item(s)
func ResetCart(c *fiber.Ctx) error {
	//connecting to DB
	cartCollection := config.MI.DB.Collection(os.Getenv("ADD_TO_CART_COLLECTION"))

	//deleting all item(s) (documents/rows) from cart
	res, err := cartCollection.DeleteMany(c.Context(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	//preparing data for json response
	returnData := new(model.ResponseData)

	if res.DeletedCount == 0 { //if cart is empty
		returnData.Status = "error"
		returnData.Message = "No product found. Cart is already empty."
	} else { //got some item(s) in the cart and deleted
		returnData.Status = "success"
		returnData.Message = "All item(s) successfully deleted from cart. Total item(s) deleted: " + fmt.Sprintf("%d", res.DeletedCount)
	}

	//encoding data to json
	/* rData, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println(err)
	}

	//w.Header().Set("Content-Type", "application/json") //setting content type as application/json
	c.Write(rData) //finally response back to client */
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": returnData,
	})

}
