package main

import (
	"flag"
	"fmt"
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
	darkModeCookieKey string = "aj-dot-dev::dark-mode"
	domain            string
	port              int
)

type navLink struct {
	Name   string
	Path   string
	Active bool
}

func main() {
	flag.StringVar(&domain, "domain", "aspenjames.dev", "Website domain")
	flag.IntVar(&port, "port", 3030, "HTTP port")
	flag.Parse()

	// Init app & template engine.
	app := fiber.New(fiber.Config{
		Views:       html.New("./content/templates", ".html"),
		ViewsLayout: "layouts/main",
	})

	// Middlewares.
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(favicon.New(favicon.Config{
		File: "./content/static/favicon.ico",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(cache.New(cache.Config{
		Expiration:   60 * time.Minute,
		CacheControl: true,
	}))
	app.Use(darkModeMiddleware)
	app.Use(navLinkMiddleware)

	// Static files.
	app.Static("/static", "./content/static")

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
func navLinkMiddleware(c *fiber.Ctx) error {
	path := c.Path()
	c.Bind(fiber.Map{
		"NavLinks": []navLink{
			{
				Name:   "Home",
				Path:   "/",
				Active: path == "/",
			},
			{
				Name:   "About",
				Path:   "/about",
				Active: path == "/about",
			},
			{
				Name:   "Resume",
				Path:   "/resume",
				Active: path == "/resume",
			},
			{
				Name:   "Particles",
				Path:   "/particles",
				Active: path == "/particles",
			},
		},
	})
	return c.Next()
}
