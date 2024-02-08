package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `json:"name"`
}

var dummy = make(map[string]string)

func main() {
	app := fiber.New()
	db := ConnectMongoDB()
	collection := db.Collection("test")
	app.Get("/api", func(c *fiber.Ctx) error {
		reqContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		queryValue := c.Query("name")
		if queryValue != "" {
			log.Println(queryValue)
			getDataByName, _ := GetAll(reqContext, bson.M{"name": bson.M{"$regex": primitive.Regex{
				Pattern: queryValue,
				Options: "i",
			}}}, collection)
			log.Println(getDataByName)
			return c.Status(fiber.StatusCreated).JSON(bson.M{
				"Status":  "200",
				"Message": "success",
				"Data":    getDataByName,
			})
		}
		getAllData, _ := GetAll(reqContext, bson.M{}, collection)
		log.Println(getAllData)
		return c.Status(fiber.StatusCreated).JSON(bson.M{
			"Status":  "200",
			"Message": "success",
			"Data":    getAllData,
		})
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		reqContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		p := new(Person)

		if err := c.BodyParser(p); err != nil {
			return err
		}
		dummy["name"] = p.Name
		doc := map[string]string{"name": p.Name}
		_, err := collection.InsertOne(reqContext, doc)
		if err != nil {
			return err
		}

		log.Println(p.Name)
		return c.SendString("post")
	})

	app.Listen(":3000")
}

func ConnectMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb+srv://kundan16239:kundan16239@cluster0.89z140s.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	db := client.Database("interviewdb")
	log.Println("Connected to MongoDB!")
	return db
}

func GetAll(ctx context.Context, doc interface{}, collection *mongo.Collection) ([]map[string]interface{}, error) {
	// Implement user creation logic using the MongoDB database connection
	cursor, err := collection.Find(ctx, doc)

	if err != nil {
		return nil, err
	}

	var elems []map[string]interface{}
	for cursor.Next(ctx) {
		// create a value into which the single document can be decoded
		var elem map[string]interface{}
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		elems = append(elems, elem)
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	return elems, nil

}
