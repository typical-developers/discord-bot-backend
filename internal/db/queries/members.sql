-- name: CreateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id, activity_points)
    VALUES (@guild_id, @member_id, @activity_points)
RETURNING *;

-- name: GetMemberRankings :one
WITH member_rankings AS (
    SELECT
        member_id,
        CAST(ROW_NUMBER() OVER (ORDER BY activity_points DESC) AS BIGINT) AS chat_rank
    FROM guild_profiles
    WHERE guild_id = @guild_id
)
SELECT * FROM member_rankings
WHERE member_id = @member_id;

-- name: GetMemberProfile :one
SELECT
    card_style,
    activity_points,
    last_grant_epoch
FROM guild_profiles
WHERE
    guild_id = @guild_id
    AND member_id = @member_id;

-- name: IncrememberMemberChatActivityPoints :one
UPDATE guild_profiles
SET
    activity_points = activity_points + @points,
    last_grant_epoch = EXTRACT(EPOCH FROM now() AT TIME ZONE 'utc')
WHERE
    guild_id = @guild_id
    AND member_id = @member_id
RETURNING *;

-- name: SetMemberChatActivityPoints :one
UPDATE guild_profiles
SET
    activity_points = @points,
    last_grant_epoch = EXTRACT(EPOCH FROM now() AT TIME ZONE 'utc')
WHERE
    guild_id = @guild_id
    AND member_id = @member_id
RETURNING *;