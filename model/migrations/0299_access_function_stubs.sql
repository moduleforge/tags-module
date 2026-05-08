-- +goose Up
--
-- Access function stub for tag resource.
--
-- Defines the signature that tag list/search queries JOIN against for
-- row-level scoping. The body here is a stub that returns the empty set.
-- At application startup, the chosen Authorizer implementation replaces
-- this body via CREATE OR REPLACE FUNCTION with the real policy.
--
-- Phase 2.2: signature gains a second parameter p_op_ids INT[] carrying the
-- satisfied-by closure for the requested operation. The old 1-arg form is
-- explicitly dropped before the 2-arg form is created so Postgres overloading
-- does not leave a stale 1-arg variant in the schema.
--
-- See core-module/docs/architecture/authorization-design.md "Row-level scoping".

-- +goose StatementBegin
DROP FUNCTION IF EXISTS accessible_tag_ids_for_actor(BIGINT);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION accessible_tag_ids_for_actor(p_actor_entity_id BIGINT, p_op_ids INT[])
RETURNS TABLE(entity_id BIGINT) LANGUAGE sql STABLE AS $$
    SELECT 0::BIGINT AS entity_id WHERE FALSE
$$;
-- +goose StatementEnd

-- +goose Down

DROP FUNCTION IF EXISTS accessible_tag_ids_for_actor(BIGINT, INT[]);
