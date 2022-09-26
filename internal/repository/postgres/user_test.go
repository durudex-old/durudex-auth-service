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

package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/internal/repository/postgres"

	"github.com/pashagolub/pgxmock"
	"github.com/segmentio/ksuid"
)

// Testing creating a new user session.
func TestUserRepository_Create(t *testing.T) {
	// Creating a new mock pool connection.
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error creating a new mock pool connection: %s", err.Error())
	}
	defer mock.Close()

	// Testing args.
	type args struct{ session domain.UserSession }

	// Test behavior.
	type mockBehavior func(args args)

	// Creating a new repository.
	repos := postgres.NewUserRepository(mock)

	// Tests structures.
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			args: args{session: domain.UserSession{
				Id:        ksuid.New(),
				UserId:    ksuid.New(),
				Payload:   "0000000000000000000000000000000000000000000000000000000000000000",
				Ip:        "0.0.0.0",
				ExpiresIn: time.Now(),
			}},
			mockBehavior: func(args args) {
				mock.ExpectExec("INSERT INTO user_session").
					WithArgs(args.session.Id, args.session.UserId, args.session.Payload, args.session.Ip, args.session.ExpiresIn).
					WillReturnResult(pgxmock.NewResult("", 1))
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			// Creating a new user session.
			err := repos.Create(context.Background(), tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("error creating a new user session: %s", err.Error())
			}
		})
	}
}

// Testing getting a user session.
func TestUserRepository_Get(t *testing.T) {
	// Creating a new mock pool connection.
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error creating a new mock pool connection: %s", err.Error())
	}
	defer mock.Close()

	// Testing args.
	type args struct{ id, userId ksuid.KSUID }

	// Test behavior.
	type mockBehavior func(args args, session domain.UserSession)

	// Creating a new repository.
	repos := postgres.NewUserRepository(mock)

	// Tests structures.
	tests := []struct {
		name         string
		args         args
		want         domain.UserSession
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			args: args{id: ksuid.New(), userId: ksuid.New()},
			want: domain.UserSession{
				Payload:   "0000000000000000000000000000000000000000000000000000000000000000",
				Ip:        "0.0.0.0",
				ExpiresIn: time.Now(),
			},
			mockBehavior: func(args args, session domain.UserSession) {
				rows := mock.NewRows([]string{"user_id", "payload", "ip", "expires_in"}).AddRow(
					args.userId, session.Payload, session.Ip, session.ExpiresIn)

				mock.ExpectQuery("SELECT (.+) FROM user_session").
					WithArgs(args.id).
					WillReturnRows(rows)
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args, tt.want)
			tt.want.UserId = tt.args.userId

			// Getting a user session.
			got, err := repos.Get(context.Background(), tt.args.id, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("error getting user session: %s", err.Error())
			}

			// Check for similarity of user session.
			if !reflect.DeepEqual(got, tt.want) {
				t.Error("error user session are not similar")
			}
		})
	}
}

// Testing getting a user session list.
func TestUserRepository_GetList(t *testing.T) {
	// Creating a new mock pool connection.
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error creating a new mock pool connection: %s", err.Error())
	}
	defer mock.Close()

	// Testing args.
	type args struct {
		userId ksuid.KSUID
		sort   domain.SortOptions
	}

	// Test behavior.
	type mockBehavior func(args args, want []domain.UserSession)

	// Creating a new repository.
	repos := postgres.NewUserRepository(mock)

	// Query filter.
	var limit int32 = 12

	// Tests structures.
	tests := []struct {
		name         string
		args         args
		want         []domain.UserSession
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name: "OK",
			args: args{
				userId: ksuid.New(),
				sort: domain.SortOptions{
					First:  &limit,
					Before: ksuid.New(),
				},
			},
			want: []domain.UserSession{
				{
					Id:        ksuid.New(),
					Payload:   "0000000000000000000000000000000000000000000000000000000000000000",
					Ip:        "0.0.0.0",
					ExpiresIn: time.Now(),
				},
			},
			mockBehavior: func(args args, want []domain.UserSession) {
				rows := mock.NewRows([]string{"id", "payload", "ip", "expires_in"}).AddRow(
					want[0].Id, want[0].Payload, want[0].Ip, want[0].ExpiresIn,
				)

				mock.ExpectQuery("SELECT (.+) FROM user_session").
					WithArgs(args.userId, args.sort.Before, *args.sort.First).
					WillReturnRows(rows)
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args, tt.want)

			// Getting a user session list.
			got, err := repos.GetList(context.Background(), tt.args.userId, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("error getting user sessions: %s", err.Error())
			}

			// Check for similarity of posts.
			if !reflect.DeepEqual(got, tt.want) {
				t.Error("error user sessions are not similar")
			}
		})
	}
}

// Testing deleting a user session.
func TestUserRepository_Delete(t *testing.T) {
	// Creating a new mock pool connection.
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error creating a new mock pool connection: %s", err.Error())
	}
	defer mock.Close()

	// Testing args.
	type args struct{ id, userId ksuid.KSUID }

	// Test behavior.
	type mockBehavior func(args args)

	// Creating a new repository.
	repos := postgres.NewUserRepository(mock)

	// Tests structures.
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockBehavior mockBehavior
	}{
		{
			name:    "OK",
			args:    args{id: ksuid.New(), userId: ksuid.New()},
			wantErr: false,
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"user_id"}).AddRow(args.userId)

				mock.ExpectBegin()

				mock.ExpectQuery("SELECT (.+) FROM user_session").
					WithArgs(args.id).
					WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM user_session").
					WithArgs(args.id).
					WillReturnResult(pgxmock.NewResult("", 1))

				mock.ExpectCommit()
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			// Deleting a user session.
			err := repos.Delete(context.Background(), tt.args.id, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("error deleting user session: %s", err.Error())
			}
		})
	}
}