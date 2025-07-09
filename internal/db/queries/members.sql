-- name: CreateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id, chat_activity)
    VALUES (@guild_id, @member_id, @chat_activity)
RETURNING *;

-- name: GetMemberRankings :one
SELECT
    ranking.member_id,
    ranking.chat_rank
FROM (
    SELECT
        member_id,
        ROW_NUMBER() OVER (
            ORDER BY chat_activity DESC
        ) AS chat_rank
    FROM guild_profiles
    WHERE
        guild_id = @guild_id
) AS ranking
WHERE
    member_id = @member_id;

-- name: GetMemberProfile :one
SELECT
    card_style,
    chat_activity,
    last_chat_activity_grant
FROM guild_profiles
WHERE
    guild_id = @guild_id
    AND member_id = @member_id;

-- name: IncrememberMemberChatActivityPoints :one
UPDATE guild_profiles
SET
    chat_activity = chat_activity + @points,
    last_chat_activity_grant = EXTRACT(EPOCH FROM now() AT TIME ZONE 'utc')
WHERE
    guild_id = @guild_id
    AND member_id = @member_id
RETURNING *;

-- name: SetMemberChatActivityPoints :one
UPDATE guild_profiles
SET
    chat_activity = @points,
    last_chat_activity_grant = EXTRACT(EPOCH FROM now() AT TIME ZONE 'utc')
WHERE
    guild_id = @guild_id
    AND member_id = @member_id
RETURNING *;

-- name: MigrateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id, card_style, chat_activity, last_chat_activity_grant)
    SELECT
        guild_id,
        @to_id,
        card_style,
        chat_activity,
        last_chat_activity_grant
    FROM guild_profiles
    WHERE
        guild_profiles.guild_id = @guild_id AND
        guild_profiles.member_id = @from_id
ON CONFLICT (guild_id, member_id) DO UPDATE
SET
    card_style = EXCLUDED.card_style,
    chat_activity = guild_profiles.chat_activity + EXCLUDED.chat_activity,
    last_chat_activity_grant = EXCLUDED.last_chat_activity_grant
RETURNING *;

-- name: ResetOldMemberProfile :exec
UPDATE guild_profiles
SET
    card_style = 0,
    chat_activity = 0,
    last_chat_activity_grant = 0
WHERE
    guild_id = @guild_id AND
    member_id = @member_id;