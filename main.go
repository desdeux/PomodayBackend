package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const filepath = "base.txt"

type Task struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	Archived   bool   `json:"archived"`
	Tag        string `json:"tag"`
	Title      string `json:"title"`
	Status     int    `json:"status"`
	Lastaction int    `json:"lastaction"`
	Logs       Logs   `json:"logs"`
}

type Log struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Logs []Log
type Tasks []Task

var db Tasks

func (tasks Tasks) readTasks(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		file.Close()

		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return err
	}

	return nil
}
func (tasks Tasks) saveTasks(filename string) error {
	data, err := json.Marshal(db)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func getList(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"tasks": db,
	})
}

func putList(c *fiber.Ctx) error {
	params := new(struct {
		Tasks Tasks `json:"tasks"`
	})
	c.BodyParser(params)
	db = (*params).Tasks
	db.saveTasks(filepath)

	return c.JSON(&fiber.Map{
		"tasks": db,
	})
}

func main() {
	port := flag.String("port", "3000", "server's port")
	login := flag.String("login", "root", "user's login")
	password := flag.String("password", "root", "user's password")

	flag.Parse()

	db = []Task{}

	err := db.readTasks(filepath)
	if err != nil {
		panic(err)
	}
	defer db.saveTasks(filepath)

	app := fiber.New()

	app.Use(cors.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			*login: *password,
		},
	}))

	app.Get("/list", getList)
	app.Put("/list", putList)

	app.Listen(":" + *port)
}
