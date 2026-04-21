-- Seed: concrete type 'tag', child of 'entity'.
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM types WHERE slug = 'entity') THEN
    RAISE EXCEPTION 'seed for tag requires parent slug ''entity'' to exist';
  END IF;
END $$;

INSERT INTO types (slug, parent_id, concrete, name, description)
SELECT
  'tag',
  id,
  true,
  'Tag',
  'A labelled annotation applied to another entity.'
FROM types WHERE slug = 'entity';
