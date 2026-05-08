-- +goose Up
--
-- Access function stub for tag resource.
--
-- Defines the signature that tag list/search queries JOIN against for
-- row-level scoping. The body here is a stub that returns the empty set.
-- At application startup, the chosen Authorizer implementation replaces
-- this body via CREATE OR REPLACE FUNCTION with the real policy.
--
-- See core-module/docs/architecture/authorization-design.md "Row-level scoping".

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION accessible_tag_ids_for_actor(p_actor_entity_id BIGINT)
RETURNS TABLE(entity_id BIGINT) LANGUAGE sql STABLE AS $$
    SELECT 0::BIGINT AS entity_id WHERE FALSE
$$;
-- +goose StatementEnd

-- +goose Down

DROP FUNCTION IF EXISTS accessible_tag_ids_for_actor(BIGINT);
