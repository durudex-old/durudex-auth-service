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
	"context"

	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/internal/service"
	v1 "github.com/durudex/durudex-auth-service/pkg/pb/durudex/v1"

	"github.com/durudex/go-protobuf-type/pbtype"
	"github.com/segmentio/ksuid"
)

// User auth gRPC handler.
type UserHandler struct {
	service service.User
	v1.UnimplementedUserAuthServiceServer
}

// Creating a new user auth gRPC handler.
func NewUserHandler(service service.User) *UserHandler {
	return &UserHandler{service: service}
}

// User Sign Up gRPC handler.
func (h *UserHandler) UserSignUp(ctx context.Context, input *v1.UserSignUpRequest) (*v1.UserSignUpResponse, error) {
	tokens, err := h.service.SignUp(ctx, domain.UserSignUpInput{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Secret:   input.Secret,
		Code:     input.Code,
		Ip:       input.Ip,
	})
	if err != nil {
		return &v1.UserSignUpResponse{}, err
	}

	return &v1.UserSignUpResponse{Access: tokens.Access, Refresh: tokens.Refresh}, nil
}

// User Sign In gRPC handler.
func (h *UserHandler) UserSignIn(ctx context.Context, input *v1.UserSignInRequest) (*v1.UserSignInResponse, error) {
	tokens, err := h.service.SignIn(ctx, domain.UserSignInInput{
		Username: input.Username,
		Password: input.Password,
		Secret:   input.Secret,
		Ip:       input.Ip,
	})
	if err != nil {
		return &v1.UserSignInResponse{}, err
	}

	return &v1.UserSignInResponse{Access: tokens.Access, Refresh: tokens.Refresh}, nil
}

// User Sign Out gRPC handler.
func (h *UserHandler) UserSignOut(ctx context.Context, input *v1.UserSignOutRequest) (*v1.UserSignOutResponse, error) {
	err := h.service.SignOut(ctx, input.Refresh, input.Secret)

	return &v1.UserSignOutResponse{}, err
}

// Refresh user authentication token gRPC handler.
func (h *UserHandler) RefreshUserToken(ctx context.Context, input *v1.RefreshUserTokenRequest) (*v1.RefreshUserTokenResponse, error) {
	access, err := h.service.RefreshToken(ctx, input.Refresh, input.Secret)
	if err != nil {
		return &v1.RefreshUserTokenResponse{}, err
	}

	return &v1.RefreshUserTokenResponse{Access: access}, nil
}

// Getting a user session gRPC handler.
func (h *UserHandler) GetUserSession(ctx context.Context, input *v1.GetUserSessionRequest) (*v1.GetUserSessionResponse, error) {
	session, err := h.service.GetSession(ctx, ksuid.FromBytesOrNil(input.Id), ksuid.FromBytesOrNil(input.UserId))
	if err != nil {
		return &v1.GetUserSessionResponse{}, err
	}

	return &v1.GetUserSessionResponse{Ip: session.Ip, ExpiresIn: pbtype.New(session.ExpiresIn)}, nil
}

// Getting a user sessions gRPC handler.
func (h *UserHandler) GetUserSessions(ctx context.Context, input *v1.GetUserSessionsRequest) (*v1.GetUserSessionsResponse, error) {
	sessions, err := h.service.GetSessions(ctx, ksuid.FromBytesOrNil(input.UserId), domain.SortOptions{
		First:  input.SortOptions.First,
		Last:   input.SortOptions.Last,
		Before: ksuid.FromBytesOrNil(input.SortOptions.Before),
		After:  ksuid.FromBytesOrNil(input.SortOptions.After),
	})
	if err != nil {
		return &v1.GetUserSessionsResponse{}, err
	}

	responseSessions := make([]*v1.UserSession, len(sessions))

	for i, session := range sessions {
		responseSessions[i] = &v1.UserSession{
			Id:        session.Id.Bytes(),
			Ip:        session.Ip,
			ExpiresIn: pbtype.New(session.ExpiresIn),
		}
	}

	return &v1.GetUserSessionsResponse{Sessions: responseSessions}, nil
}

// Deleting a user session gRPC handler.
func (h *UserHandler) DeleteUserSession(ctx context.Context, input *v1.DeleteUserSessionRequest) (*v1.DeleteUserSessionResponse, error) {
	err := h.service.DeleteSession(ctx, input.Refresh, input.Secret)

	return &v1.DeleteUserSessionResponse{}, err
}

// Getting total user session count gRPC handler.
func (h *UserHandler) GetTotalUserSessionCount(ctx context.Context, input *v1.GetTotalUserSessionCountRequest) (*v1.GetTotalUserSessionCountResponse, error) {
	count, err := h.service.GetTotalCount(ctx, ksuid.FromBytesOrNil(input.UserId))
	if err != nil {
		return &v1.GetTotalUserSessionCountResponse{}, err
	}

	return &v1.GetTotalUserSessionCountResponse{Count: count}, nil
}
