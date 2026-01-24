-- This index doesn't seem to make much of a difference.
-- Will revisit this in the future, most likely was not understanding this completely.
DROP INDEX IF EXISTS guild_profiles_rankings_index;

-- Renames the activty_points and last_grant_epoch columns to be more specific.
ALTER TABLE guild_profiles
RENAME COLUMN activity_points TO chat_activity;
ALTER TABLE guild_profiles
RENAME COLUMN last_grant_epoch TO last_chat_activity_grant;

-- Adds new voice activity point columns & role blacklist columns.
ALTER TABLE guild_profiles
ADD COLUMN voice_activity INT DEFAULT 0 NOT NULL,
ADD COLUMN last_voice_activity_grant INT DEFAULT 0 NOT NULL,
ADD COLUMN chat_grant_deny TEXT[] DEFAULT '{}' NOT NULL,
ADD COLUMN voice_grant_deny TEXT[] DEFAULT '{}' NOT NULL;