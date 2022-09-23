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
	err := h.service.SignOut(ctx, domain.UserSignOutInput{
		Refresh: input.Refresh,
		Secret:  input.Secret,
		Ip:      input.Ip,
	})

	return &v1.UserSignOutResponse{}, err
}

// Refresh user authentication token gRPC handler.
func (h *UserHandler) RefreshUserToken(ctx context.Context, input *v1.RefreshUserTokenRequest) (*v1.RefreshUserTokenResponse, error) {
	access, err := h.service.RefreshToken(ctx, domain.UserRefreshTokenInput{
		Refresh: input.Refresh,
		Secret:  input.Secret,
		Ip:      input.Ip,
	})
	if err != nil {
		return &v1.RefreshUserTokenResponse{}, err
	}

	return &v1.RefreshUserTokenResponse{Access: access}, nil
}
