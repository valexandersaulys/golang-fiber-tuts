package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"os"
	"time"
)

// note: no commas!
type Person struct {
	Name  string `json:"name" xml:"name" form:"name"`
	Email string `json:"email" xml:"email" form:"email"`
}

func LogTimeMiddlware(c *fiber.Ctx) error {
	fmt.Println("Date: ", time.Now())
	// fiber.IsChild() =>  returns true meaning a child process in Fiber.Config{ Prefork:True }
	return c.Next()
}

func CreateApp() fiber.App {
	engine := django.New(
		"./views",
		".jinja2",
	)

	app := fiber.New(fiber.Config{
		Views:         engine, // pass in a Views engine for Load & Render methods
		ServerHeader:  "MY Server",
		StrictRouting: true,            // treat /foo and /foo/ as different routes
		CaseSensitive: false,           // treat /Foo and /foo as different routes
		Immutable:     true,            // make all handler values immutable
		BodyLimit:     4 * 1024 * 1024, // maximum allowed size for a request body
		// Concurrency:   256 * 1024,      // maximum allowed size for a request body
		// ReadTimeout:     nil,   // supposed to default to nil
		// WriteTimeout:    nil,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		ProxyHeader:     "",    // enables ctx.IP to return the value of the given header key (e.g. `X-Forwarded-*`)
		GETOnly:         false, // reject all non-GET requests
		// ErrorHandler:    // for handling e.g. 404   https://docs.gofiber.io/guide/error-handling/#custom-error-handler
		DisableStartupMessage: false,
	})

	app.Use(LogTimeMiddlware)

	app.Static(
		"/static",
		"./public",
	)

	// I think the whole second arg is a function that takes and returns an error?
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	app.Get("/first", func(c *fiber.Ctx) error {
		return c.Render("first", fiber.Map{"Title": "Hey!"})
	})

	app.Get("/named-routes/:one/:two?", func(c *fiber.Ctx) error {
		return c.Render("second", fiber.Map{
			"one": c.Params("one"),
			"two": c.Params("two"),
		})
	})

	app.Get("/greedy/", func(c *fiber.Ctx) error {
		return c.SendString(c.Query("oogah"))
	})

	app.Post("/person/", func(c *fiber.Ctx) error {
		person := new(Person) // instantiate new person
		if err := c.BodyParser(person); err != nil {
			return err
		}
		// c.Set("Content-Type", "application/json")
		return c.JSON(person)
	})

	app.Mount("/users", UserRoutes())
	MountOtherRoutes(app)

	if os.Getenv("PRINT_STACK") == "TRUE" {
		data, _ := json.MarshalIndent(app.Stack(), "", " ")
		fmt.Println(string(data))
	}
	if os.Getenv("CONFIG_KEY") != "" {
		fmt.Println(app.Config()) // sadly not indexable
	}

	return *app
}

func main() {
	app := CreateApp()
	app.Listen(":3000")
}
