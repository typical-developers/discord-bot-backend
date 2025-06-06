CREATE TABLE IF NOT EXISTS guild_activity_tracking_monthly_current (
    grant_type TEXT NOT NULL DEFAULT 'chat',
    guild_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    earned_points INT NOT NULL,
    PRIMARY KEY (grant_type, guild_id, member_id)
);

CREATE TABLE IF NOT EXISTS guild_activity_tracking_weekly_current (
    grant_type TEXT NOT NULL DEFAULT 'chat',
    guild_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    earned_points INT NOT NULL,
    PRIMARY KEY (grant_type, guild_id, member_id)
);

INSERT INTO guild_activity_tracking_monthly_current (
    guild_id, member_id, earned_points
)
SELECT guild_id, member_id, earned_points
FROM guild_activity_tracking_monthly
WHERE month_start = EXTRACT(epoch FROM date_trunc('month', now() AT TIME ZONE 'utc'));

INSERT INTO guild_activity_tracking_weekly_current (
    guild_id, member_id, earned_points
)
SELECT guild_id, member_id, earned_points
FROM guild_activity_tracking_weekly
WHERE week_start = EXTRACT(epoch FROM date_trunc('week', now() AT TIME ZONE 'utc'));

DELETE FROM guild_activity_tracking_monthly
WHERE month_start = EXTRACT(epoch FROM date_trunc('month', now() AT TIME ZONE 'utc'));

DELETE FROM guild_activity_tracking_weekly
WHERE week_start = EXTRACT(epoch FROM date_trunc('week', now() AT TIME ZONE 'utc'));