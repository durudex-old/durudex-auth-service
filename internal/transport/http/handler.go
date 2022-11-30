/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package http

import (
	"github.com/durudex/durudex-auth-service/internal/service"
	v1 "github.com/durudex/durudex-auth-service/internal/transport/http/v1"

	"github.com/gofiber/fiber/v2"
)

// Handler interface stores methods for interacting with the API handler.
type Handler interface {
	// RegisterAllRoutes method registers all API routes.
	RegisterAllRoutes(fiber.Router)
}

// handler structure implements methods for interacting with the API.
type handler struct{ service *service.Service }

// NewHandler function returns a new API handler.
func NewHandler(service *service.Service) Handler { return &handler{service: service} }

// RegisterAllRoutes method registers all API routes.
func (h *handler) RegisterAllRoutes(router fiber.Router) {
	v1.NewHandler(h.service).RegisterAllRoutes(router)
}
