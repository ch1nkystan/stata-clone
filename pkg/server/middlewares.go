package server

import (
	"github.com/gofiber/fiber/v2"
)

const (
	XAdminToken = "X-Admin-Token"
)

func (s *Server) apiMiddlewareBackend(c *fiber.Ctx) error {
	if _, ok := s.cfg.BackendTokens[c.Get(XAdminToken)]; !ok {
		return s.Unauthorized(c)
	}

	return c.Next()
}

func (s *Server) apiMiddlewareFrontend(c *fiber.Ctx) error {
	if _, ok := s.cfg.FrontendTokens[c.Get(XAdminToken)]; !ok {
		return s.Unauthorized(c)
	}

	return c.Next()
}
