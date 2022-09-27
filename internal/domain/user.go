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

package domain

import (
	"time"

	"github.com/segmentio/ksuid"
)

// User auth session.
type UserSession struct {
	// User session id.
	Id ksuid.KSUID
	// User session owner.
	UserId ksuid.KSUID
	// User session payload.
	Payload string
	// User ip address.
	Ip string
	// User session expires in.
	ExpiresIn time.Time
}

// User SignUp auth input.
type UserSignUpInput struct {
	// Unique username.
	Username string
	// Unique user email address.
	Email string
	// User password hash.
	Password string
	// Client secret key.
	Secret string
	// Verification code.
	Code uint64
	// User ip address.
	Ip string
}

// User SignIn auth input.
type UserSignInInput struct {
	// Username.
	Username string
	// User password hash.
	Password string
	// Client secret key.
	Secret string
	// User ip address.
	Ip string
}

//	User auth tokens.
type UserTokens struct {
	// JWT access token.
	Access string
	// Refresh token.
	Refresh string
}
