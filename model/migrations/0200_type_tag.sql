-- +goose Up

-- Seed: concrete type 'tag', child of 'entity'.
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM types WHERE slug = 'entity') THEN
    RAISE EXCEPTION 'seed for tag requires parent slug ''entity'' to exist';
  END IF;
END $$;
-- +goose StatementEnd

INSERT INTO types (slug, parent_id, concrete, name, description)
SELECT
  'tag',
  id,
  true,
  'Tag',
  'A labelled annotation applied to another entity.'
FROM types WHERE slug = 'entity';
