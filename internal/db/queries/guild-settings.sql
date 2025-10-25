-- name: RegisterGuild :one
INSERT INTO guilds (guild_id)
VALUES (@guild_id)
RETURNING *;

-- name: GetGuildChatActivitySettings :one
SELECT
    is_enabled,
    grant_amount,
    grant_cooldown,
    deny_roles
FROM guild_chat_activity_settings
WHERE
    guild_chat_activity_settings.guild_id = @guild_id
LIMIT 1;

-- name: GetGuildVoiceActivitySettings :one
SELECT
    is_enabled,
    grant_amount,
    grant_cooldown,
    deny_roles
FROM guild_voice_activity_settings
WHERE
    guild_voice_activity_settings.guild_id = @guild_id
LIMIT 1;

-- name: GetGuildActivityRoles :many
SELECT
    role_id,
    required_points
FROM guild_activity_roles
WHERE
    guild_activity_roles.guild_id = @guild_id
    AND grant_type = @activity_type
GROUP BY grant_type, role_id, required_points
ORDER BY required_points ASC;

-- name: UpdateGuildChatActivitySettings :exec
UPDATE guild_chat_activity_settings SET
    is_enabled = COALESCE(sqlc.narg(is_enabled), guild_chat_activity_settings.is_enabled),
    grant_amount = COALESCE(sqlc.narg(grant_amount), guild_chat_activity_settings.grant_amount),
    grant_cooldown = COALESCE(sqlc.narg(grant_cooldown), guild_chat_activity_settings.grant_cooldown)
WHERE
    guild_id = @guild_id;

-- name: UpdateGuildVoiceActivitySettings :exec
UPDATE guild_voice_activity_settings SET
    is_enabled = COALESCE(sqlc.narg(is_enabled), guild_voice_activity_settings.is_enabled),
    grant_amount = COALESCE(sqlc.narg(grant_amount), guild_voice_activity_settings.grant_amount),
    grant_cooldown = COALESCE(sqlc.narg(grant_cooldown), guild_voice_activity_settings.grant_cooldown)
WHERE
    guild_id = @guild_id;

-- name: InsertActivityRole :exec
INSERT INTO guild_activity_roles (guild_id, grant_type, role_id, required_points)
    VALUES (@guild_id, @grant_type, @role_id, @required_points::INT);

-- name: DeleteActivityRole :exec
DELETE FROM guild_activity_roles
WHERE
    guild_id = @guild_id
    AND role_id = @role_id;