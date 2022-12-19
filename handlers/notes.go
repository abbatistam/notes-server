package handlers

import (
	"main/database"
	"main/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllNotes(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := database.Mg.Db.Collection("notes").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var notes []models.Note = make([]models.Note, 0)

	if err := cursor.All(c.Context(), &notes); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(notes)
}

func NewNote(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("notes")

	note := new(models.Note)

	if err := c.BodyParser(note); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	note.ID = ""

	insertionResult, err := collection.InsertOne(c.Context(), note)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdNote := &models.Note{}
	createdRecord.Decode(createdNote)

	return c.Status(201).JSON(createdNote)
}

func EditNote(c *fiber.Ctx) error {
	idParam := c.Params("id")
	noteID, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		return c.SendStatus(400)
	}

	note := new(models.Note)

	if err := c.BodyParser(note); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: noteID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: note.Name},
				{Key: "body", Value: note.Body},
			},
		},
	}
	err = database.Mg.Db.Collection("notes").FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.SendStatus(500)
	}

	note.ID = idParam
	return c.Status(200).JSON(note)
}

func DeleteNote(c *fiber.Ctx) error {
	noteID, err := primitive.ObjectIDFromHex(
		c.Params("id"),
	)

	if err != nil {
		return c.SendStatus(400)
	}

	query := bson.D{{Key: "_id", Value: noteID}}
	result, err := database.Mg.Db.Collection("notes").DeleteOne(c.Context(), &query)

	if err != nil {
		return c.SendStatus(500)
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}

	return c.SendStatus(204)
}
