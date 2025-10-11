-- name: GetAllTimeActivityLeaderboard :many
SELECT
    rankings.member_id,
    rankings.rank,
    rankings.points
FROM (
    SELECT 
        member_id,
        ROW_NUMBER() OVER (
            ORDER BY
                CASE LOWER(@activity_type)::TEXT
                    WHEN 'chat' THEN chat_activity
                    WHEN 'voice' THEN voice_activity
                END DESC,
                CASE LOWER(@activity_type)::TEXT
                    WHEN 'chat' THEN last_chat_activity_grant
                    WHEN 'voice' THEN last_voice_activity_grant
                END DESC
        ) AS rank,
        CAST (
            CASE LOWER(@activity_type)::TEXT
                WHEN 'chat' THEN chat_activity
                WHEN 'voice' THEN voice_activity
            END AS INT
        ) AS points
    FROM guild_profiles
    WHERE
        guild_profiles.guild_id = @guild_id
) AS rankings
LIMIT @limit_by
OFFSET @offset_by;

-- name: GetWeeklyActivityLeaderboard :many
SELECT
    rankings.rank,
    rankings.member_id,
    rankings.earned_points
FROM (
    SELECT
        ROW_NUMBER() OVER (
            ORDER BY earned_points DESC
        ) AS rank,
        member_id,
        earned_points
    FROM guild_activity_tracking_weekly_current
    WHERE
        guild_id = @guild_id
        AND grant_type = @grant_type
) AS rankings
LIMIT 15
OFFSET @offset_by;

-- name: GetMonthlyActivityLeaderboard :many
SELECT
    rankings.rank,
    rankings.member_id,
    rankings.earned_points
FROM (
    SELECT
        ROW_NUMBER() OVER (
            ORDER BY earned_points DESC
        ) AS rank,
        member_id,
        earned_points
    FROM guild_activity_tracking_monthly_current
    WHERE
        guild_id = @guild_id
        AND grant_type = @grant_type
) AS rankings
LIMIT 15
OFFSET @offset_by;

-- name: IncrementWeeklyActivityLeaderboard :exec
INSERT INTO guild_activity_tracking_weekly_current (
    grant_type, guild_id, member_id, earned_points
)
VALUES (
    @grant_type, @guild_id, @member_id, @earned_points
)
ON CONFLICT (grant_type, guild_id, member_id)
DO UPDATE SET
    earned_points = guild_activity_tracking_weekly_current.earned_points + @earned_points
WHERE
    guild_activity_tracking_weekly_current.grant_type = @grant_type;

-- name: IncrementMonthlyActivityLeaderboard :exec
INSERT INTO guild_activity_tracking_monthly_current (
    grant_type, guild_id, member_id, earned_points
)
VALUES (
    @grant_type, @guild_id, @member_id, @earned_points
)
ON CONFLICT (grant_type, guild_id, member_id)
DO UPDATE SET
    earned_points = guild_activity_tracking_monthly_current.earned_points + @earned_points
WHERE
    guild_activity_tracking_monthly_current.grant_type = @grant_type;