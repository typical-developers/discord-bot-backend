-- name: CreateGuildSettings :one
-- Creates new generic guild settings.
INSERT INTO guild_settings (guild_id)
VALUES (@guild_id)
RETURNING *;

-- name: GetGuildSettings :one
SELECT
    insert_epoch,
    chat_activity_tracking,
    chat_activity_grant,
    chat_activity_cooldown
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

-- name: GetVoiceRoomLobbies :many
SELECT *
FROM guild_voice_rooms_settings
WHERE
    guild_voice_rooms_settings.guild_id = @guild_id;

-- name: CreateVoiceRoomLobby :one
INSERT INTO guild_voice_rooms_settings (
    guild_id, voice_channel_id,
    user_limit, can_rename, can_lock, can_adjust_limit
)
VALUES (
    @guild_id, @voice_channel_id,
    @user_limit, @can_rename, @can_lock, @can_adjust_limit
)
RETURNING *;

-- name: UpdateActivitySettings :exec
INSERT INTO
    guild_settings (
        guild_id,
        chat_activity_tracking, chat_activity_grant, chat_activity_cooldown,
        voice_activity_tracking, voice_activity_grant, voice_activity_cooldown
    )
    VALUES (
        @guild_id,
        COALESCE(@chat_activity_tracking::BOOLEAN, FALSE),
        COALESCE(@chat_activity_grant::INT, 2),
        COALESCE(@chat_activity_cooldown::INT, 15),
        COALESCE(@voice_activity_tracking::BOOLEAN, FALSE),
        COALESCE(@voice_activity_grant::INT, 2),
        COALESCE(@voice_activity_cooldown::INT, 15)
    )
ON CONFLICT (guild_id)
DO UPDATE SET
    chat_activity_tracking = COALESCE(sqlc.narg(chat_activity_tracking), guild_settings.chat_activity_tracking),
    chat_activity_grant = COALESCE(sqlc.narg(chat_activity_grant), guild_settings.chat_activity_grant),
    chat_activity_cooldown = COALESCE(sqlc.narg(chat_activity_cooldown), guild_settings.chat_activity_cooldown),
    voice_activity_tracking = COALESCE(sqlc.narg(voice_activity_tracking), guild_settings.voice_activity_tracking),
    voice_activity_grant = COALESCE(sqlc.narg(voice_activity_grant), guild_settings.voice_activity_grant),
    voice_activity_cooldown = COALESCE(sqlc.narg(voice_activity_cooldown), guild_settings.voice_activity_cooldown)
RETURNING *;

-- name: InsertActivityRole :exec
INSERT INTO guild_activity_roles (guild_id, grant_type, role_id, required_points)
    VALUES (@guild_id, @grant_type, @role_id, @required_points::INT)
ON CONFLICT DO NOTHING;