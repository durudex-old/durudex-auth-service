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

package client

import (
	"github.com/durudex/durudex-auth-service/internal/config"
	v1 "github.com/durudex/durudex-auth-service/pkg/pb/durudex/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Client structure.
type Client struct {
	User  *UserClient
	Code  *CodeClient
	Email *EmailClient
}

// User client structure.
type UserClient struct {
	v1.UserServiceClient
	conn *grpc.ClientConn
}

// Code client structure.
type CodeClient struct {
	v1.UserCodeServiceClient
	conn *grpc.ClientConn
}

// Email client structure.
type EmailClient struct {
	v1.EmailUserServiceClient
	conn *grpc.ClientConn
}

// Creating a new client.
func NewClient(cfg config.ServiceConfig) *Client {
	log.Debug().Msg("Creating a new client...")

	userServiceConn := ConnectToGRPCService(cfg.User)
	codeServiceConn := ConnectToGRPCService(cfg.Code)
	emailServiceConn := ConnectToGRPCService(cfg.Email)

	return &Client{
		User: &UserClient{
			v1.NewUserServiceClient(userServiceConn),
			userServiceConn,
		},
		Code: &CodeClient{
			v1.NewUserCodeServiceClient(codeServiceConn),
			codeServiceConn,
		},
		Email: &EmailClient{
			v1.NewEmailUserServiceClient(emailServiceConn),
			emailServiceConn,
		},
	}
}

// Closing a client connections.
func (c *Client) Close() {
	log.Info().Msg("Closing a client connections")

	// Closing user service connection.
	if err := c.User.conn.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close user service connection")
	}
	// Closing code service connection.
	if err := c.Code.conn.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close code service connection")
	}
	// Closing email service connection.
	if err := c.Email.conn.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to close email service connection")
	}
}
