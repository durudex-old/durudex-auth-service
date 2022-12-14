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

package postgres

import (
	"context"
	"errors"

	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/pkg/database/postgres"

	"github.com/jackc/pgx/v4"
	"github.com/leporo/sqlf"
	"github.com/segmentio/ksuid"
)

// User session repository interface.
type Session interface {
	// Creating a new user session.
	Create(ctx context.Context, session domain.UserSession) error
	// Getting a user session.
	Get(ctx context.Context, userId, id ksuid.KSUID) (domain.UserSession, error)
	// Getting a user sessions list.
	GetList(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error)
	// Deleting a user session.
	Delete(ctx context.Context, userId, id ksuid.KSUID) error
	// Getting total user session count.
	GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error)
}

// User session repository structure.
type SessionRepository struct{ psql postgres.Postgres }

// Creating a new use session postgres repository.
func NewSessionRepository(psql postgres.Postgres) *SessionRepository {
	return &SessionRepository{psql: psql}
}

// Creating a new user session.
func (r *SessionRepository) Create(ctx context.Context, session domain.UserSession) error {
	query := "INSERT INTO user_session (id, user_id, payload, ip, expires_in) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.psql.Exec(ctx, query, session.Id, session.UserId, session.Payload, session.Ip, session.ExpiresIn)

	return err
}

// Getting a user session.
func (r *SessionRepository) Get(ctx context.Context, userId, id ksuid.KSUID) (domain.UserSession, error) {
	var session domain.UserSession

	query := "SELECT payload, ip, expires_in FROM user_session WHERE user_id=$1 AND id=$2"
	row := r.psql.QueryRow(ctx, query, userId, id)

	// Scanning query row.
	if err := row.Scan(&session.Payload, &session.Ip, &session.ExpiresIn); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserSession{}, &domain.Error{Code: domain.CodeNotFound, Message: "Session not found"}
		}

		return domain.UserSession{}, err
	}

	return session, nil
}

// Getting a user sessions list.
func (r *SessionRepository) GetList(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error) {
	var n int32

	qb := sqlf.Select("id, ip, expires_in").From("user_session").Where("user_id = ?", userId)

	// Added first or last sort option.
	if sort.First != nil {
		n = *sort.First
		qb.OrderBy("user_id ASC").Limit(*sort.First)
	} else if sort.Last != nil {
		n = *sort.Last
		qb.OrderBy("user_id DESC").Limit(*sort.Last)
	}

	// Added before sort option.
	if sort.Before != ksuid.Nil {
		qb.Where("id < ?", sort.Before)
	}
	// Added after sort option.
	if sort.After != ksuid.Nil {
		qb.Where("id > ?", sort.After)
	}

	sessions := make([]domain.UserSession, n)

	// Query for getting author posts by author id.
	rows, err := r.psql.Query(ctx, qb.String(), qb.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var i int

	// Scanning query rows.
	for rows.Next() {
		var session domain.UserSession

		// Scanning query row.
		if err := rows.Scan(&session.Id, &session.Ip, &session.ExpiresIn); err != nil {
			return nil, err
		}

		sessions[i] = session

		i++
	}

	// Check is rows error.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Check for fullness of the slice.
	if i == int(n) {
		return sessions, nil
	}

	res := make([]domain.UserSession, i)
	copy(res, sessions[:i])

	return res, nil
}

// TODO: fix payload.
// Deleting a user session.
func (r *SessionRepository) Delete(ctx context.Context, userId, id ksuid.KSUID) error {
	// Deleting user session.
	_, err := r.psql.Exec(ctx, "DELETE FROM user_session WHERE id=$1", id)
	return err
}

// Getting total user session count.
func (r *SessionRepository) GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error) {
	var count int32

	// Get total user session count.
	query := "SELECT count(*) FROM user_session WHERE author_id=$1"
	row := r.psql.QueryRow(ctx, query, userId)

	// Scanning query row.
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
