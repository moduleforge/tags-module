-- Optional filters use the (sqlc.narg('name')::type IS NULL OR col = sqlc.narg('name')::type)
-- pattern so the generated params struct has readable, nullable field names (pgtype.Text,
-- pgtype.Int8, etc.) instead of positional Column1/Column2 names.

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

-- name: UpdateTagColor :one
UPDATE tags
SET color = @color
WHERE entity_id = @entity_id
RETURNING entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE entity_id = $1;

-- name: ListTagsBySubjectEntityID :many
SELECT entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at
FROM tags
WHERE subject_id = sqlc.arg('subject_id')
  AND (sqlc.narg('purpose')::text IS NULL OR purpose = sqlc.narg('purpose')::text)
ORDER BY created_at ASC;

-- name: SearchTags :many
SELECT entity_id, owner_id, subject_id, purpose, value, color, created_at, updated_at
FROM tags
WHERE (sqlc.narg('owner_id')::bigint IS NULL OR owner_id = sqlc.narg('owner_id')::bigint)
  AND (sqlc.narg('subject_id')::bigint IS NULL OR subject_id = sqlc.narg('subject_id')::bigint)
  AND (sqlc.narg('purpose')::text IS NULL OR purpose = sqlc.narg('purpose')::text)
  AND (sqlc.narg('value')::text IS NULL OR value = sqlc.narg('value')::text)
ORDER BY created_at ASC;

-- name: CountTagsBySubjectEntityID :one
SELECT COUNT(*)
FROM tags
WHERE subject_id = $1;
