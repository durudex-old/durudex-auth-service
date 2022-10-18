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

// User session gRPC handler.
type SessionHandler struct {
	service service.Session
	v1.UnimplementedUserSessionServiceServer
}

// Creating a new user session gRPC handler.
func NewSessionHandler(service service.Session) *SessionHandler {
	return &SessionHandler{service: service}
}

// Getting a user session gRPC handler.
func (h *SessionHandler) GetUserSession(ctx context.Context, input *v1.GetUserSessionRequest) (*v1.GetUserSessionResponse, error) {
	session, err := h.service.Get(ctx, ksuid.FromBytesOrNil(input.UserId), ksuid.FromBytesOrNil(input.Id))
	if err != nil {
		return &v1.GetUserSessionResponse{}, err
	}

	return &v1.GetUserSessionResponse{
		Ip:        session.Ip,
		ExpiresIn: pbtype.New(session.ExpiresIn),
	}, nil
}

// Getting a user sessions gRPC handler.
func (h *SessionHandler) GetUserSessions(ctx context.Context, input *v1.GetUserSessionsRequest) (*v1.GetUserSessionsResponse, error) {
	sessions, err := h.service.GetList(ctx, ksuid.FromBytesOrNil(input.UserId), domain.SortOptions{
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

// Deleting a user session.
func (h *SessionHandler) DeleteUserSession(ctx context.Context, input *v1.DeleteUserSessionRequest) (*v1.DeleteUserSessionResponse, error) {
	err := h.service.Delete(ctx, ksuid.FromBytesOrNil(input.UserId), ksuid.FromBytesOrNil(input.Id))
	return &v1.DeleteUserSessionResponse{}, err
}

// Getting total user session count gRPC handler.
func (h *SessionHandler) GetTotalUserSessionCount(ctx context.Context, input *v1.GetTotalUserSessionCountRequest) (*v1.GetTotalUserSessionCountResponse, error) {
	count, err := h.service.GetTotalCount(ctx, ksuid.FromBytesOrNil(input.UserId))
	if err != nil {
		return &v1.GetTotalUserSessionCountResponse{}, err
	}

	return &v1.GetTotalUserSessionCountResponse{Count: count}, nil
}
