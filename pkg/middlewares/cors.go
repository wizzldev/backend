package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/url"
	"slices"
)

var corsOriginDomains = []string{"wizzl.app", "localhost"}

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowHeaders:     "Authorization,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,X-Frontend-Client,X-File-Access-Token",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		MaxAge:           60 * 5, // 5 minutes
		AllowOriginsFunc: func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			return slices.Contains(corsOriginDomains, u.Hostname())
		},
	})
}
