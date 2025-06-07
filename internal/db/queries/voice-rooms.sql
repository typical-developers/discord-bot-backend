-- name: GetVoiceRoomLobbies :many
SELECT *
FROM guild_voice_rooms_settings
WHERE
    guild_voice_rooms_settings.guild_id = @guild_id
LIMIT 1;

-- name: CreateVoiceRoomLobby :one
INSERT INTO guild_voice_rooms_settings (guild_id, voice_channel_id)
VALUES (@guild_id, @voice_channel_id)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetVoiceRoomLobby :one
SELECT *
FROM guild_voice_rooms_settings
WHERE
    guild_id = @guild_id
    AND voice_channel_id = @voice_channel_id
LIMIT 1;

-- name: UpdateVoiceRoomLobby :one
UPDATE guild_voice_rooms_settings
SET
    user_limit = COALESCE(sqlc.narg(user_limit), guild_voice_rooms_settings.user_limit),
    can_rename = COALESCE(sqlc.narg(can_rename), guild_voice_rooms_settings.can_rename),
    can_lock = COALESCE(sqlc.narg(can_lock), guild_voice_rooms_settings.can_lock),
    can_adjust_limit = COALESCE(sqlc.narg(can_adjust_limit), guild_voice_rooms_settings.can_adjust_limit)
WHERE
    guild_id = @guild_id
    AND voice_channel_id = @voice_channel_id
RETURNING *;

-- name: DeleteVoiceRoomLobby :exec
DELETE FROM guild_voice_rooms_settings
WHERE
    guild_id = @guild_id
    AND voice_channel_id = @voice_channel_id;