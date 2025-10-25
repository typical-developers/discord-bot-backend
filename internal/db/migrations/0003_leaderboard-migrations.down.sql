-- Insert the data back into the original monthly table
INSERT INTO guild_activity_tracking_monthly (
    guild_id, member_id, earned_points, month_start
)
SELECT 
    guild_id, 
    member_id, 
    earned_points,
    EXTRACT(epoch FROM date_trunc('month', now() AT TIME ZONE 'utc')) AS month_start
FROM guild_activity_tracking_monthly_current;

-- Insert the data back into the original weekly table
INSERT INTO guild_activity_tracking_weekly (
    guild_id, member_id, earned_points, week_start
)
SELECT 
    guild_id, 
    member_id, 
    earned_points,
    EXTRACT(epoch FROM date_trunc('week', now() AT TIME ZONE 'utc')) AS week_start
FROM guild_activity_tracking_weekly_current;

-- Optionally, drop the current tables if they are no longer needed
DROP TABLE IF EXISTS guild_activity_tracking_monthly_current;
DROP TABLE IF EXISTS guild_activity_tracking_weekly_current;
