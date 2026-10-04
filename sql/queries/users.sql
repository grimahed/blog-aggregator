-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
	$1,
	$2,
	$3,
	$4
)

RETURNING *;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url,
			user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)

RETURNING *;

-- name: GetUser :one
SELECT
   *
FROM
   users
WHERE
   name = $1;

-- name: GetUsers :many
SELECT
   *
FROM
   users
ORDER BY
   name ASC;

-- name: GetAllFeedsWithUserName :many
SELECT
   users.name,
   feeds.*
FROM
   users
RIGHT JOIN feeds ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS
(
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)

VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)


RETURNING
   *
)

SELECT
  inserted_feed_follow.*,
  users.name AS name_of_user,
  feeds.name AS name_of_feed
FROM
  inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many

SELECT
  *,
  users.name AS users_name,
  feeds.name AS feed_name
FROM
  feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
GROUP BY
  feed_follows.id,
  users.id,
  feeds.id
HAVING
  users.id = $1;

-- name: GetFeed :one
SELECT
   *
FROM
   feeds
WHERE
   url = $1;

-- name: DeleteAll :exec
DELETE FROM
   users;
