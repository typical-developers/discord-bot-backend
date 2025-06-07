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

-- name: ResetWeeklyActivityLeaderboard :exec
MERGE INTO guild_activity_tracking_weekly AS archive
USING (
    SELECT
        EXTRACT(epoch FROM date_trunc('week', now() AT TIME ZONE 'utc') - INTERVAL '1 week')::INT AS week_start,
        *
    FROM guild_activity_tracking_weekly_current
) AS current_leaderboard
ON
    archive.week_start = current_leaderboard.week_start
    AND archive.guild_id = current_leaderboard.guild_id
    AND archive.member_id = current_leaderboard.member_id
WHEN MATCHED THEN
    UPDATE SET
        earned_points = current_leaderboard.earned_points
WHEN NOT MATCHED THEN
    INSERT (week_start, grant_type, guild_id, member_id, earned_points)
    VALUES (
        current_leaderboard.week_start,
        current_leaderboard.grant_type,
        current_leaderboard.guild_id,
        current_leaderboard.member_id,
        current_leaderboard.earned_points
    );

-- name: TruncateWeeklyActivityLeaderboard :exec
TRUNCATE TABLE guild_activity_tracking_weekly_current;

-- name: ResetMonthlyActivityLeaderboard :exec
MERGE INTO guild_activity_tracking_monthly AS archive
USING (
    SELECT
        EXTRACT(epoch FROM date_trunc('month', now() AT TIME ZONE 'utc') - INTERVAL '1 month')::INT AS month_start,
        *
    FROM guild_activity_tracking_monthly_current
) AS current_leaderboard
ON
    archive.month_start = current_leaderboard.month_start
    AND archive.guild_id = current_leaderboard.guild_id
    AND archive.member_id = current_leaderboard.member_id
WHEN MATCHED THEN
    UPDATE SET
        earned_points = current_leaderboard.earned_points
WHEN NOT MATCHED THEN
    INSERT (month_start, grant_type, guild_id, member_id, earned_points)
    VALUES (
        current_leaderboard.month_start,
        current_leaderboard.grant_type,
        current_leaderboard.guild_id,
        current_leaderboard.member_id,
        current_leaderboard.earned_points
    );

-- name: TruncateMonthlyActivityLeaderboard :exec
TRUNCATE TABLE guild_activity_tracking_monthly_current;