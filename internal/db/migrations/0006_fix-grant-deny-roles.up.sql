-- I have no clue what I was on when I wrote the original migration.
-- 
-- This removes the chat_grant_deny and voice_grant_deny columns from guild profiles
-- and adds them to the originally inteded table, guild_settings.
-- 
-- Also modifies old chat activity field names and adds new fields for the new voice grant columns.

ALTER TABLE guild_profiles
DROP COLUMN chat_grant_deny,
DROP COLUMN voice_grant_deny;

-- Rename existing activity grant columns to be defined as chat grant.
-- Also adds a new grant deny roles column.
ALTER TABLE guild_settings
RENAME COLUMN activity_tracking TO chat_activity_tracking;
ALTER TABLE guild_settings
RENAME COLUMN activity_tracking_grant TO chat_activity_grant;
ALTER TABLE guild_settings
RENAME COLUMN activity_tracking_cooldown TO chat_activity_cooldown;
ALTER TABLE guild_settings
ADD COLUMN chat_activity_deny_roles TEXT[] DEFAULT '{}' NOT NULL;

ALTER TABLE guild_settings
ADD COLUMN voice_activity_tracking BOOLEAN DEFAULT FALSE,
ADD COLUMN voice_activity_grant INT DEFAULT 2,
ADD COLUMN voice_activity_cooldown INT DEFAULT 15,
ADD COLUMN voice_grant_deny TEXT[] DEFAULT '{}' NOT NULL;