-- name: CreateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id, activity_points)
    VALUES (@guild_id, @member_id, @activity_points)
RETURNING *;

-- name: GetMemberRankings :one
SELECT
    ranking.member_id,
    ranking.chat_rank
FROM (
    SELECT
        member_id,
        ROW_NUMBER() OVER (
            ORDER BY activity_points DESC
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

-- name: MigrateMemberProfile :one
INSERT INTO guild_profiles (guild_id, member_id, card_style, activity_points, last_grant_epoch)
    SELECT
        guild_id,
        @to_id,
        card_style,
        activity_points,
        last_grant_epoch
    FROM guild_profiles
    WHERE
        guild_profiles.guild_id = @guild_id AND
        guild_profiles.member_id = @from_id
ON CONFLICT (guild_id, member_id) DO UPDATE
SET
    card_style = EXCLUDED.card_style,
    activity_points = guild_profiles.activity_points + EXCLUDED.activity_points,
    last_grant_epoch = EXCLUDED.last_grant_epoch
RETURNING *;

-- name: ResetOldMemberProfile :exec
UPDATE guild_profiles
SET
    card_style = 0,
    activity_points = 0,
    last_grant_epoch = 0
WHERE
    guild_id = @guild_id AND
    member_id = @member_id;