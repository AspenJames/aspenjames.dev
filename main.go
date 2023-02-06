package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

var (
	content           string
	darkModeCookieKey string = "aj-dot-dev::dark-mode"
	domain            string
	port              int
)

type navLink struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Active bool
}

func main() {
	flag.StringVar(&content, "content", "/usr/src/content", "Website content directory")
	flag.StringVar(&domain, "domain", "aspenjames.dev", "Website domain")
	flag.IntVar(&port, "port", 80, "HTTP port")
	flag.Parse()

	// Init app & template engine.
	app := fiber.New(fiber.Config{
		Views:       html.New(fmt.Sprintf("%s/templates", content), ".html"),
		ViewsLayout: "layouts/main",
	})

	// Middlewares.
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(favicon.New(favicon.Config{
		File: fmt.Sprintf("%s/static/favicon.ico", content),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(cache.New(cache.Config{
		Expiration:   7 * 24 * time.Hour,
		CacheControl: true,
		Next: func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() >= 400
		},
	}))
	app.Use(darkModeMiddleware)

	links, err := parseNavLinks()
	if err != nil {
		os.Exit(1)
	}
	app.Use(navLinkMiddleware(links))

	// Static files.
	app.Static("/static", fmt.Sprintf("%s/static", content))

	// Routes.
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/:path", func(c *fiber.Ctx) error {
		path := c.Params("path")
		templatePath := strings.TrimSuffix(path, ".html")
		err := c.Render(templatePath, fiber.Map{})
		if err != nil {
			if match, _ := regexp.MatchString("template .* does not exist", err.Error()); match {
				return c.Next()
			}
		}
		return err
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("errors/404", fiber.Map{}, "layouts/error")
	})

	// Run app.
	appPort := fmt.Sprintf(":%d", port)
	app.Listen(appPort)
}

// Sets default color theme cookie & bind template vars.
func darkModeMiddleware(c *fiber.Ctx) error {
	if c.Cookies(darkModeCookieKey) == "" {
		expires := time.Now().Add(time.Hour * 24 * 365).UTC()
		darkModeCookie := &fiber.Cookie{
			Name:     darkModeCookieKey,
			Domain:   domain,
			SameSite: "Strict",
			Expires:  expires,
			Value:    "light",
		}
		c.Cookie(darkModeCookie)
	}
	c.Bind(fiber.Map{
		"DarkModeCookieKey": darkModeCookieKey,
		"DarkMode":          c.Cookies(darkModeCookieKey) == "dark",
	})
	return c.Next()
}

// Define navigation links & bind template var.
func navLinkMiddleware(navLinks []*navLink) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentPath := c.Path()
		for _, l := range navLinks {
			l.Active = l.Path == currentPath
		}
		c.Bind(fiber.Map{
			"NavLinks":   navLinks,
			"ActivePath": c.Path(),
		})
		return c.Next()
	}
}

func parseNavLinks() ([]*navLink, error) {
	links := make([]*navLink, 0)
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/routes.json", content))
	if err != nil {
		return links, err
	}
	err = json.Unmarshal(data, &links)
	if err != nil {
		return links, err
	}
	return links, nil
}
