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

CREATE TABLE IF NOT EXISTS user_session (
  id         CHAR(27)    NOT NULL PRIMARY KEY,
  user_id    CHAR(27)    NOT NULL,
  payload    VARCHAR(64) NOT NULL,
  ip         INET        NOT NULL,
  expires_in TIMESTAMP   NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS user_session_id_idx ON user_session (id DESC);

CREATE INDEX IF NOT EXISTS user_session_user_id_idx ON user_session (user_id DESC);
