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

package service

import "github.com/durudex/durudex-auth-service/internal/domain"

// Account interface stores methods for interacting with the account service.
type Account interface {
	// Check method checks the user's login credentials.
	Check(domain.Login, string) error
}

// account structure implements methods for interacting with the account service.
type account struct{}

// NewAccount function returns a new account service.
func NewAccount() Account { return &account{} }

// Check method checks the user's login credentials.
func (a *account) Check(lt domain.Login, login string) error { return nil }
