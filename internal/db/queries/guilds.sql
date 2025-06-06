-- name: GetAllTimeChatActivityRankings :many
SELECT
    CAST (
        ROW_NUMBER() OVER (ORDER BY guild_profiles.activity_points DESC) AS INT
    ) AS rank,
    member_id,
    activity_points
FROM guild_profiles
WHERE
    guild_id = @guild_id
ORDER BY activity_points DESC
LIMIT 15
OFFSET @offset_by;

-- name: GetWeeklyActivityLeaderboard :many
SELECT
    CAST (
        ROW_NUMBER() OVER (ORDER BY guild_activity_tracking_weekly_current.earned_points DESC) AS INT
    ) AS rank,
    member_id,
    earned_points
FROM guild_activity_tracking_weekly_current
WHERE
    guild_id = @guild_id
ORDER BY earned_points DESC
LIMIT 15
OFFSET @offset_by;

-- name: GetMonthlyActivityLeaderboard :many
SELECT
    CAST (
        ROW_NUMBER() OVER (ORDER BY guild_activity_tracking_monthly_current.earned_points DESC) AS INT
    ) AS rank,
    member_id,
    earned_points
FROM guild_activity_tracking_monthly_current
WHERE
    guild_id = @guild_id
ORDER BY earned_points DESC
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