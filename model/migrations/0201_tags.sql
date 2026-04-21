-- tags is the single-subject tag subtype table. Identity and lifecycle columns
-- (uuid, archived_at) live on the parent entities row, not here.

CREATE TABLE tags (
  entity_id   BIGINT PRIMARY KEY REFERENCES entities(id) ON DELETE RESTRICT,
  owner_id    BIGINT NOT NULL REFERENCES entities(id) ON DELETE RESTRICT,
  subject_id  BIGINT NOT NULL REFERENCES entities(id) ON DELETE RESTRICT,
  purpose     TEXT NOT NULL CHECK (char_length(purpose) <= 512),
  value       TEXT NOT NULL CHECK (char_length(value) <= 512),
  color       TEXT CHECK (color SIMILAR TO '#[0-9A-Fa-f]{8}'),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX tags_owner_subject_purpose_idx
  ON tags (owner_id, subject_id, purpose);

CREATE INDEX tags_subject_id_idx ON tags (subject_id);
-- owner_id alone is already covered by the leading column of the unique index;
-- no separate index needed.

-- Enforce that only entities whose fundamental type descends from 'tag'
-- may be inserted into this table.
CREATE FUNCTION tags_check_type() RETURNS TRIGGER AS $$
DECLARE
  v_type_id BIGINT;
BEGIN
  SELECT fundamental_type_id INTO v_type_id FROM entities WHERE id = NEW.entity_id;
  IF NOT FOUND THEN
    RAISE EXCEPTION 'tags: entity_id % does not reference a known entity', NEW.entity_id;
  END IF;
  IF NOT type_is_or_descends_from(v_type_id, 'tag') THEN
    RAISE EXCEPTION 'tags: entity % fundamental type does not descend from tag', NEW.entity_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tags_type_check
  BEFORE INSERT ON tags
  FOR EACH ROW EXECUTE FUNCTION tags_check_type();

-- Enforce that owner_id, subject_id, purpose, and value are immutable after insert.
CREATE FUNCTION tags_reject_immutable_changes() RETURNS TRIGGER AS $$
BEGIN
  IF OLD.owner_id IS DISTINCT FROM NEW.owner_id THEN
    RAISE EXCEPTION 'tags: owner_id is immutable after insert';
  END IF;
  IF OLD.subject_id IS DISTINCT FROM NEW.subject_id THEN
    RAISE EXCEPTION 'tags: subject_id is immutable after insert';
  END IF;
  IF OLD.purpose IS DISTINCT FROM NEW.purpose THEN
    RAISE EXCEPTION 'tags: purpose is immutable after insert';
  END IF;
  IF OLD.value IS DISTINCT FROM NEW.value THEN
    RAISE EXCEPTION 'tags: value is immutable after insert';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tags_reject_immutable_changes
  BEFORE UPDATE ON tags
  FOR EACH ROW EXECUTE FUNCTION tags_reject_immutable_changes();

CREATE TRIGGER tags_set_updated_at
  BEFORE UPDATE ON tags
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
