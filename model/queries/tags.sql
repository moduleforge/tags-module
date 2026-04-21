-- Optional filters use the ($N::type IS NULL OR col = $N::type) pattern,
-- matching the convention established in users-module/model/queries/user_accounts.sql.

-- name: CreateTag :one
INSERT INTO tags (entity_id, owner_id, subject_id, purpose, value, color)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at;

-- name: GetTagByEntityID :one
SELECT entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at
FROM tags
WHERE entity_id = $1;

-- name: GetTagByEntityUUID :one
SELECT t.entity_id, t.owner_id, t.subject_id, t.purpose, t.value, t.color,
       t.created_at, t.updated_at, e.uuid
FROM tags t
JOIN entities e ON e.id = t.entity_id
WHERE e.uuid = $1;

-- name: UpdateTagColor :exec
UPDATE tags
SET color = $2
WHERE entity_id = $1;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE entity_id = $1;

-- name: ListTagsBySubjectEntityID :many
SELECT entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at
FROM tags
WHERE subject_id = $1
  AND ($2::text IS NULL OR purpose = $2::text)
ORDER BY created_at ASC;

-- name: SearchTags :many
SELECT entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at
FROM tags
WHERE ($1::bigint IS NULL OR owner_id = $1::bigint)
  AND ($2::bigint IS NULL OR subject_id = $2::bigint)
  AND ($3::text IS NULL OR purpose = $3::text)
  AND ($4::text IS NULL OR value = $4::text)
ORDER BY created_at ASC;

-- name: CountTagsBySubjectEntityID :one
SELECT COUNT(*)
FROM tags
WHERE subject_id = $1;
