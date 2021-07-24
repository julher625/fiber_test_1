package main

import (
	"fmt"
	"log"

	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	// projectID      = "mechDesing"
	projectID      = "mechdesing-1644d"
	credentialFile = "mechdesing-1644d-firebase-adminsdk-8aqoj-406b7479bc.json"
)

func initFireStore(ctx *context.Context) (*firestore.Client, *firebase.App, error) {

	opt := option.WithCredentialsFile(credentialFile)
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(*ctx, conf, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(*ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client, app, nil

}

type Sensor struct {
	Name string `json:"name" xml:"name" form:"name"`
}

func main() {
	app := fiber.New(
		fiber.Config{
			Prefork:       true,
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Fiber",
			AppName:       "Mech Desing API",
		})

	ctx := context.Background()

	client, _, _ := initFireStore(&ctx)

	app.Server().MaxConnsPerIP = 1

	app.Post("/api/sensor", func(c *fiber.Ctx) error {

		sensor := Sensor{}
		if err := c.BodyParser(&sensor); err != nil {
			// return c.Send(c.Body())
			return c.SendStatus(fiber.StatusBadRequest)
		}
		_, _, err := client.Collection("sensors").Add(ctx, map[string]interface{}{
			"name": sensor.Name,
		})
		if err != nil {
			log.Fatalf("Failed adding lovelace: %v", err)
		}
		str := fmt.Sprintf("Added: %v", sensor.Name)
		log.Println(sensor.Name)
		return c.Send([]byte(str))
	})

	app.Get("/api/sensor/:id", func(c *fiber.Ctx) error {
		// id := c.Params("id")
		aux := Sensor{}
		iter := client.Collection("sensors").Documents(ctx)
		var str string
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			err = doc.DataTo(&aux)
			if err != nil {
				log.Fatalf("Failed Data to : %v", err)
			}
			res := fmt.Sprintf("%v", aux.Name)
			str = str + res

		}
		return c.Send([]byte(str))

		// return nil
		// response := fmt.Sprintf("ID: %s", id)
		// return json.NewEncoder(c.Response().BodyWriter()).Encode(tr)
		// return c.Send([]byte([]byte(response)))
	})

	app.Listen("localhost:8000")
	defer client.Close()
}
