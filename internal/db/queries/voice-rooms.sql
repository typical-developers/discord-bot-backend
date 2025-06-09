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

-- name: RegisterVoiceRoom :one
INSERT INTO guild_active_voice_rooms (
    guild_id, origin_channel_id,
    channel_id, created_by_user_id, current_owner_id
)
VALUES (
    @guild_id, @origin_channel_id,
    @channel_id, @created_by_user_id, @current_owner_id
)
RETURNING *;

-- name: GetVoiceRooms :many
SELECT * FROM guild_active_voice_rooms
WHERE
    guild_id = @guild_id
    AND origin_channel_id = @origin_channel_id;

-- name: UpdateVoiceRoom :one
UPDATE guild_active_voice_rooms
SET
    current_owner_id = COALESCE(sqlc.narg(current_owner_id), guild_active_voice_rooms.current_owner_id),
    is_locked = COALESCE(sqlc.narg(is_locked), guild_active_voice_rooms.is_locked)
WHERE
    guild_id = @guild_id
    AND channel_id = @channel_id
RETURNING *;

-- name: DeleteVoiceRoom :exec
DELETE FROM guild_active_voice_rooms
WHERE
    guild_id = @guild_id
    AND channel_id = @channel_id;