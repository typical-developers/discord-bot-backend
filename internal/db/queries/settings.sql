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
ORDER BY required_points ASC;

-- name: CreateGuildSettings :one
INSERT INTO guild_settings (guild_id)
VALUES (@guild_id)
RETURNING *;

-- name: UpdateActivitySettings :exec
INSERT INTO
    guild_settings (guild_id, activity_tracking, activity_tracking_grant, activity_tracking_cooldown)
    VALUES (@guild_id, @activity_tracking, @activity_tracking_grant, @activity_tracking_cooldown)
ON CONFLICT (guild_id)
DO UPDATE SET
    activity_tracking = COALESCE(sqlc.narg(activity_tracking), guild_settings.activity_tracking),
    activity_tracking_grant = COALESCE(sqlc.narg(activity_tracking_grant), guild_settings.activity_tracking_grant),
    activity_tracking_cooldown = COALESCE(sqlc.narg(activity_tracking_cooldown), guild_settings.activity_tracking_cooldown)
RETURNING *;

-- name: InsertActivityRole :exec
INSERT INTO guild_activity_roles (guild_id, grant_type, role_id, required_points)
    VALUES (@guild_id, @grant_type, @role_id, @required_points::INT)
ON CONFLICT DO NOTHING;