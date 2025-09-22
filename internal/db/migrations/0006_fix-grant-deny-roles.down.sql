-- Add the columns back to guild_profiles
ALTER TABLE guild_profiles
ADD COLUMN chat_grant_deny TEXT[] DEFAULT '{}' NOT NULL,
ADD COLUMN voice_grant_deny TEXT[] DEFAULT '{}' NOT NULL;

-- Rename the chat_activity_* columns back to original names
ALTER TABLE guild_settings
RENAME COLUMN chat_activity_tracking TO activity_tracking;
ALTER TABLE guild_settings
RENAME COLUMN chat_activity_grant TO activity_tracking_grant;
ALTER TABLE guild_settings
RENAME COLUMN chat_activity_cooldown TO activity_tracking_cooldown;

-- Drop the new chat and voice columns from guild_settings
ALTER TABLE guild_settings
DROP COLUMN chat_activity_deny_roles,
DROP COLUMN voice_activity_tracking,
DROP COLUMN voice_activity_grant,
DROP COLUMN voice_activity_cooldown,
DROP COLUMN voice_grant_deny;
