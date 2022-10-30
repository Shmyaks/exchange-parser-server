// Package core core of app
package core

import (
	"fmt"

	"github.com/Shmyaks/exchange-parser-server/app/config"
	"github.com/Shmyaks/exchange-parser-server/app/internal"
	"github.com/Shmyaks/exchange-parser-server/app/internal/delivery/routes"
	jsoniter "github.com/json-iterator/go"

	// docs are generated by Swag CLI
	_ "github.com/Shmyaks/exchange-parser-server/app/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Server struct
type Server struct {
	app *fiber.App
}

// NewServer fabric
func NewServer() *Server {
	return &Server{app: fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		JSONDecoder: jsoniter.Unmarshal,
	})}
}

// StartControllers start goroutines of controllers
func (s *Server) StartControllers(applicationContainer *internal.ApplicationContainer) {
	go applicationContainer.Controllers.P2PController.Parse()
	go applicationContainer.Controllers.SpotController.Parse()
}

// Start Init app
func (s *Server) Start() {
	applicationContainer := internal.InitApplication()
	s.StartControllers(applicationContainer)
	routes.InitV1Routes(s.app, applicationContainer)
	s.app.Get("/swagger/*", swagger.HandlerDefault)
	s.app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none"}))
	s.app.Listen(fmt.Sprintf(": %s", config.Env.BackendPort))
}