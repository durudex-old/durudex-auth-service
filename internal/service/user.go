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
	"fmt"
	"time"

	"github.com/durudex/durudex-auth-service/internal/client"
	"github.com/durudex/durudex-auth-service/internal/config"
	"github.com/durudex/durudex-auth-service/internal/domain"
	"github.com/durudex/durudex-auth-service/pkg/auth"
	v1 "github.com/durudex/durudex-auth-service/pkg/pb/durudex/v1"

	"github.com/durudex/go-refresh"
	"github.com/segmentio/ksuid"
)

// User auth service interface.
type User interface {
	// User SignUp.
	SignUp(ctx context.Context, input domain.UserSignUpInput) (domain.UserTokens, error)
	// User SignIn.
	SignIn(ctx context.Context, input domain.UserSignInInput) (domain.UserTokens, error)
	// Creating a new user session.
	CreateSession(ctx context.Context, userId ksuid.KSUID, ip, secret string) (domain.UserTokens, error)
	// Refresh user token.
	RefreshToken(ctx context.Context, token, secret string) (string, error)
}

// User service structure.
type UserService struct {
	session Session
	// Service client.
	client *client.Client
	// Auth config variables.
	cfg *config.AuthConfig
}

// Creating a new user service.
func NewUserService(session Session, client *client.Client, cfg *config.AuthConfig) *UserService {
	return &UserService{session: session, client: client, cfg: cfg}
}

// User SignUp.
func (s *UserService) SignUp(ctx context.Context, input domain.UserSignUpInput) (domain.UserTokens, error) {
	// Verifying user email code.
	emailResponse, err := s.client.Code.VerifyUserEmailCode(ctx, &v1.VerifyUserEmailCodeRequest{
		Email: input.Email,
		Code:  input.Code,
	})
	if err != nil || !emailResponse.Status {
		return domain.UserTokens{}, err
	}

	// Creating a new user.
	userResponse, err := s.client.User.CreateUser(ctx, &v1.CreateUserRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return domain.UserTokens{}, err
	}

	// Creating a new user session.
	tokens, err := s.CreateSession(ctx, ksuid.FromBytesOrNil(userResponse.Id), input.Ip, input.Secret)
	if err != nil {
		return domain.UserTokens{}, err
	}

	// Sending an email to a user with register.
	if _, err := s.client.Email.SendEmailUserRegister(ctx, &v1.SendEmailUserRegisterRequest{
		Email:    input.Email,
		Username: input.Username,
	}); err != nil {
		return domain.UserTokens{}, err
	}

	return tokens, nil
}

// User SignIn.
func (s *UserService) SignIn(ctx context.Context, input domain.UserSignInInput) (domain.UserTokens, error) {
	// Getting a user by credentials.
	userResponse, err := s.client.User.GetUserByCreds(ctx, &v1.GetUserByCredsRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return domain.UserTokens{}, err
	}

	// Creating a new user session.
	tokens, err := s.CreateSession(ctx, ksuid.FromBytesOrNil(userResponse.Id), input.Ip, input.Secret)
	if err != nil {
		return domain.UserTokens{}, err
	}

	// Sending an email to a user with logged in.
	if _, err := s.client.Email.SendEmailUserLoggedIn(ctx, &v1.SendEmailUserLoggedInRequest{
		Email: userResponse.Email,
		Ip:    input.Ip,
	}); err != nil {
		return domain.UserTokens{}, err
	}

	return tokens, nil
}

// Creating a new user session.
func (s *UserService) CreateSession(ctx context.Context, userId ksuid.KSUID, ip, secret string) (domain.UserTokens, error) {
	// Generating a new refresh token.
	r, err := refresh.New()
	if err != nil {
		return domain.UserTokens{}, err
	}

	// Generate a new session id.
	sessionId := ksuid.New()

	// Creating a new user session.
	if err := s.session.Create(ctx, domain.UserSession{
		Id:        sessionId,
		UserId:    userId,
		Payload:   fmt.Sprintf("%x", r.Hash([]byte(secret))),
		Ip:        ip,
		ExpiresIn: time.Now().Add(s.cfg.Session.TTL),
	}); err != nil {
		return domain.UserTokens{}, err
	}

	// Generating a new jwt access token.
	access, err := auth.GenerateAccessToken(userId.String(), s.cfg.JWT.SigningKey, s.cfg.JWT.TTL)
	if err != nil {
		return domain.UserTokens{}, err
	}

	return domain.UserTokens{Refresh: r.Token(sessionId.String(), userId.String()), Access: access}, nil
}

// Refresh user token.
func (s *UserService) RefreshToken(ctx context.Context, token, secret string) (string, error) {
	// Parsing refresh token string.
	r, err := refresh.Parse(token)
	if err != nil {
		return "", err
	}

	// Parsing session id string.
	id, err := ksuid.Parse(r.Session)
	if err != nil {
		return "", err
	}

	// Parsing user id string.
	userId, err := ksuid.Parse(r.Object)
	if err != nil {
		return "", err
	}

	// Getting a user session.
	session, err := s.session.Get(ctx, userId, id)
	if err != nil {
		return "", err
	}

	// Checking user session payload for similar input payload.
	if session.Payload != fmt.Sprintf("%x", r.Payload.Hash([]byte(secret))) {
		return "", &domain.Error{Code: domain.CodeInvalidArgument, Message: "Session payload is not similar"}
	}

	// Generating a new jwt access token.
	access, err := auth.GenerateAccessToken(session.UserId.String(), s.cfg.JWT.SigningKey, s.cfg.JWT.TTL)
	if err != nil {
		return "", err
	}

	return access, nil
}
