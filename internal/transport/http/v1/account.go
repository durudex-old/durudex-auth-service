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

package v1

import (
	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

// AccountHandler interface stores methods for interacting with the API account handler.
type AccountHandler interface {
	// RegisterAccountRoutes method registers accounts API routes.
	RegisterAccountRoutes(fiber.Router)
}

// accountHandler structure implements methods for interacting with the account API.
type accountHandler struct{ service service.Account }

// NewAccount returns a new account service.
func NewAccount(service service.Account) AccountHandler {
	return &accountHandler{service: service}
}

// RegisterAccountRoutes method registers accounts API routes.
func (h *accountHandler) RegisterAccountRoutes(router fiber.Router) {
	user := router.Group("/accounts")
	{
		user.Get("/check", h.userCheck)
	}
}

// UserCheckRequest structure stores request data for checking account login credential.
type UserCheckRequest struct {
	// LoginType field stores the user's login credential type number.
	LoginType domain.Login `json:"loginType"`

	// Login field stores the value of the user's login credentials.
	Login string `json:"login"`
}

// userCheck method implements a handler for checking account login credentials.
func (h *accountHandler) userCheck(ctx *fiber.Ctx) error {
	var input UserCheckRequest

	// Parsing request body.
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error parsing request body")
	}

	// Checking user login credential.
	return h.service.Check(input.LoginType, input.Login)
}
