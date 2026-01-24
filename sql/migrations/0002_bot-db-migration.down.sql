ALTER TABLE guild_activity_roles
DROP COLUMN IF EXISTS grant_type;

ALTER TABLE guild_profiles
ALTER COLUMN last_grant_epoch DROP DEFAULT;

ALTER TABLE guild_activity_tracking_monthly
DROP COLUMN IF EXISTS grant_type;

ALTER TABLE guild_activity_tracking_weekly
DROP COLUMN IF EXISTS grant_type;