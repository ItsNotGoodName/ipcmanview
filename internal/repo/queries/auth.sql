-- name: AuthCreateUser :one
INSERT INTO
  users (
    email,
    username,
    password,
    created_at,
    updated_at,
    disabled_at
  )
VALUES
  (?, ?, ?, ?, ?, ?) RETURNING id;

-- name: AuthGetUser :one
SELECT
  *
FROM
  users
WHERE
  id = ?;

-- name: AuthGetUserByUsernameOrEmail :one
SELECT
  *
FROM
  users
WHERE
  username = sqlc.arg ('username_or_email')
  OR email = sqlc.arg ('username_or_email');

-- name: AuthListUsersByGroup :many
SELECT
  users.*
FROM
  users
  LEFT JOIN group_users ON group_users.user_id = id
WHERE
  group_users.group_id = ?;

-- name: AuthPatchUser :one
UPDATE users
SET
  username = coalesce(sqlc.narg ('username'), username),
  email = coalesce(sqlc.narg ('email'), email),
  password = coalesce(sqlc.narg ('password'), password),
  updated_at = sqlc.arg ('updated_at')
WHERE
  id = sqlc.arg ('id') RETURNING id;

-- name: AuthUpdateUserDisabledAt :one
UPDATE users
SET
  disabled_at = ?
WHERE
  id = ? RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = ?;

-- name: AuthCreateUserSession :exec
INSERT INTO
  user_sessions (
    user_id,
    session,
    user_agent,
    ip,
    last_ip,
    last_used_at,
    created_at,
    expired_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?);

-- name: AuthGetUserSessionForContext :one
SELECT
  user_sessions.id as id,
  user_sessions.user_id as user_id,
  users.username,
  admins.user_id IS NOT NULL as 'admin',
  user_sessions.last_ip,
  user_sessions.last_used_at,
  users.disabled_at AS 'users_disabled_at',
  user_sessions.session
FROM
  user_sessions
  LEFT JOIN users ON users.id = user_sessions.user_id
  LEFT JOIN admins ON admins.user_id = user_sessions.user_id
WHERE
  session = ?
  AND expired_at > sqlc.arg ('now');

-- name: AuthDeleteUserSessionForUser :exec
DELETE FROM user_sessions
WHERE
  id = ?
  AND user_id = ?;

-- name: AuthDeleteUserSessionByExpired :exec
DELETE FROM user_sessions
WHERE
  expired_at < ?;

-- name: AuthListUserSessionsForUserAndNotExpired :many
SELECT
  *
FROM
  user_sessions
WHERE
  user_id = ?
  AND expired_at > sqlc.arg ('now');

-- name: AuthUpdateUserSession :exec
UPDATE user_sessions
SET
  last_ip = ?,
  last_used_at = ?
WHERE
  id = ?;

-- name: AuthDeleteUserSessionForUserAndNotSession :exec
DELETE FROM user_sessions
WHERE
  user_id = ?
  AND id != ?;

-- name: AuthDeleteUserSessionBySession :exec
DELETE FROM user_sessions
WHERE
  session = ?;

-- name: AuthListGroupsForUser :many
SELECT
  g.*,
  gu.created_at AS joined_at
FROM
  groups AS g
  LEFT JOIN group_users AS gu ON gu.group_id = g.id
WHERE
  gu.user_id = ?;

-- name: AuthCountGroup :one
SELECT
  count(*)
FROM
  groups;

-- name: AuthGetGroup :one
SELECT
  *
FROM
  groups
where
  id = ?;

-- name: AuthCreateGroup :one
INSERT INTO
  groups (name, description, created_at, updated_at)
VALUES
  (?, ?, ?, ?) RETURNING id;

-- name: AuthUpdateGroup :one
UPDATE groups
SET
  name = ?,
  description = ?,
  updated_at = ?
WHERE
  id = ? RETURNING id;

-- name: AuthDeleteGroup :exec
DELETE FROM groups
WHERE
  id = ?;

-- name: AuthUpdateGroupDisabledAt :one
UPDATE groups
SET
  disabled_at = ?
WHERE
  id = ? RETURNING id;

-- name: AuthUpsertAdmin :one
INSERT OR IGNORE INTO
  admins (user_id, created_at)
VALUES
  (?, ?) RETURNING user_id;

-- name: AuthDeleteAdmin :exec
DELETE FROM admins
WHERE
  user_id = ?;
