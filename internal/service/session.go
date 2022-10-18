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

import (
	"context"

	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/internal/repository/postgres"

	"github.com/segmentio/ksuid"
)

// User session service.
type Session interface {
	// Creating a new user session.
	Create(ctx context.Context, session domain.UserSession) error
	// Getting user session.
	Get(ctx context.Context, userId, id ksuid.KSUID) (domain.UserSession, error)
	// Getting user sessions list.
	GetList(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error)
	// Deleting user session.
	Delete(ctx context.Context, userId, id ksuid.KSUID) error
	// Getting total user session count.
	GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error)
}

// User session service structure.
type SessionService struct{ repos postgres.Session }

// Creating a new user session service.
func NewSessionService(repos postgres.Session) *SessionService {
	return &SessionService{repos: repos}
}

// Creating a new user session.
func (s *SessionService) Create(ctx context.Context, session domain.UserSession) error {
	return s.repos.Create(ctx, session)
}

// Getting user session.
func (s *SessionService) Get(ctx context.Context, userId, id ksuid.KSUID) (domain.UserSession, error) {
	return s.repos.Get(ctx, userId, id)
}

// Getting user sessions list.
func (s *SessionService) GetList(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error) {
	// Checking is first and last are set.
	if sort.First == nil && sort.Last == nil {
		return nil, &domain.Error{Message: "Must be `first` or `last`", Code: domain.CodeInvalidArgument}
	}

	return s.repos.GetList(ctx, userId, sort)
}

// Deleting user session.
func (s *SessionService) Delete(ctx context.Context, userId, id ksuid.KSUID) error {
	return s.repos.Delete(ctx, userId, id)
}

// Getting total user session count.
func (s *SessionService) GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error) {
	return s.repos.GetTotalCount(ctx, userId)
}
