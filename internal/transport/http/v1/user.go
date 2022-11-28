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

	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes method registers users API routes.
func (h *handler) RegisterUserRoutes(router fiber.Router) {
	user := router.Group("/users")
	{
		user.Get("/check", h.userCheck)
	}
}

// UserCheckRequest structure stores request data for checking user credential.
type UserCheckRequest struct {
	// CredentialType field stores the user's credential type number.
	CredentialType domain.Credential `json:"credentialType"`

	// Credential field stores the value of the user's credentials.
	Credential string `json:"credential"`
}

// UserCheckResponse structure stores response data for checking user credential.
type UserCheckResponse struct {
	// Used field stores the usage status of credentials.
	Used bool `json:"used"`

	// Valid field stores the validation status according to Durudex standards.
	Valid bool `json:"valid"`
}

// userCheck method implements a handler for checking user credentials.
func (h *handler) userCheck(ctx *fiber.Ctx) error {
	return ctx.SendString("hello")
}
