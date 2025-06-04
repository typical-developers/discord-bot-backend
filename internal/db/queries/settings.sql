-- name: GetGuildSettings :one
SELECT
    insert_epoch,
    activity_tracking,
    activity_tracking_grant,
    activity_tracking_cooldown
FROM guild_settings
WHERE
    guild_settings.guild_id = @guild_id
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
ORDER BY required_points DESC;

-- name: CreateGuildSettings :one
INSERT INTO guild_settings (guild_id)
VALUES (@guild_id)
RETURNING *;
