-- First, rename the columns back so activity_points exists
ALTER TABLE guild_profiles
RENAME COLUMN chat_activity TO activity_points;

ALTER TABLE guild_profiles
RENAME COLUMN last_chat_activity_grant TO last_grant_epoch;

-- Now you can safely create the index on activity_points
CREATE INDEX IF NOT EXISTS guild_profiles_rankings_index
ON guild_profiles (guild_id, activity_points DESC);

-- Drop voice/chat grant related columns
ALTER TABLE guild_profiles
DROP COLUMN IF EXISTS voice_activity,
DROP COLUMN IF EXISTS last_voice_activity_grant,
DROP COLUMN IF EXISTS chat_grant_deny,
DROP COLUMN IF EXISTS voice_grant_deny;
