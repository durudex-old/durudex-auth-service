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
	"github.com/durudex/durudex-auth-service/internal/repository/postgres"
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
	// User SignOut.
	SignOut(ctx context.Context, token, secret string) error
	// Refresh user token.
	RefreshToken(ctx context.Context, token, secret string) (string, error)
	// Getting a user session.
	GetSession(ctx context.Context, id, userId ksuid.KSUID) (domain.UserSession, error)
	// Getting a user sessions list.
	GetSessions(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error)
	// Deleting a user session.
	DeleteSession(ctx context.Context, token, secret string) error
	// Getting total user session count.
	GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error)
}

// User service structure.
type UserService struct {
	// User auth repository.
	auth postgres.User
	// Service client.
	client *client.Client
	// Auth config variables.
	cfg *config.AuthConfig
}

// Creating a new user service.
func NewUserService(repos postgres.User, client *client.Client, cfg *config.AuthConfig) *UserService {
	return &UserService{auth: repos, client: client, cfg: cfg}
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
	if err := s.auth.Create(ctx, domain.UserSession{
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

	return domain.UserTokens{Refresh: r.Token(sessionId.String()), Access: access}, nil
}

// User SignOut.
func (s *UserService) SignOut(ctx context.Context, token, secret string) error {
	// Parsing refresh token string.
	r, sId, err := refresh.Parse(token)
	if err != nil {
		return err
	}

	// Parsing session id string.
	id, err := ksuid.Parse(sId)
	if err != nil {
		return err
	}

	// Hashing refresh token by secret key.
	payload := r.Hash([]byte(secret))

	// Deleting a user session.
	return s.auth.Delete(ctx, id, fmt.Sprintf("%x", payload))
}

// Refresh user token.
func (s *UserService) RefreshToken(ctx context.Context, token, secret string) (string, error) {
	// Parsing refresh token string.
	r, sId, err := refresh.Parse(token)
	if err != nil {
		return "", err
	}

	// Parsing session id string.
	id, err := ksuid.Parse(sId)
	if err != nil {
		return "", err
	}

	// Getting a user session.
	session, err := s.auth.Get(ctx, id)
	if err != nil {
		return "", err
	}

	// Hashing refresh token by secret key.
	payload := r.Hash([]byte(secret))

	// Checking user session payload for similar input payload.
	if session.Payload != fmt.Sprintf("%x", payload) {
		return "", &domain.Error{Code: domain.CodeInvalidArgument, Message: "Session payload is not similar"}
	}

	// Generating a new jwt access token.
	access, err := auth.GenerateAccessToken(session.UserId.String(), s.cfg.JWT.SigningKey, s.cfg.JWT.TTL)
	if err != nil {
		return "", err
	}

	return access, nil
}

// Getting a user session.
func (s *UserService) GetSession(ctx context.Context, id, userId ksuid.KSUID) (domain.UserSession, error) {
	// Getting a user session.
	session, err := s.auth.Get(ctx, id)
	if err != nil {
		return domain.UserSession{}, err
	}

	// Checking user session owner for similar input user id.
	if session.UserId != userId {
		return domain.UserSession{}, &domain.Error{Code: domain.CodeInvalidArgument, Message: "User id is not similar"}
	}

	return session, nil
}

// Getting a user sessions list.
func (s *UserService) GetSessions(ctx context.Context, userId ksuid.KSUID, sort domain.SortOptions) ([]domain.UserSession, error) {
	// Checking is first and last are set.
	if sort.First == nil && sort.Last == nil {
		return nil, &domain.Error{Message: "Must be `first` or `last`", Code: domain.CodeInvalidArgument}
	}

	// Getting a user sessions list.
	sessions, err := s.auth.GetList(ctx, userId, sort)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// Deleting a user session.
func (s *UserService) DeleteSession(ctx context.Context, token, secret string) error {
	// Parsing refresh token string.
	r, sId, err := refresh.Parse(token)
	if err != nil {
		return err
	}

	// Parsing session id string.
	id, err := ksuid.Parse(sId)
	if err != nil {
		return err
	}

	// Hashing refresh token by secret key.
	payload := r.Hash([]byte(secret))

	// Deleting a user session.
	return s.auth.Delete(ctx, id, fmt.Sprintf("%x", payload))
}

// Getting total user session count.
func (s *UserService) GetTotalCount(ctx context.Context, userId ksuid.KSUID) (int32, error) {
	return s.auth.GetTotalCount(ctx, userId)
}
