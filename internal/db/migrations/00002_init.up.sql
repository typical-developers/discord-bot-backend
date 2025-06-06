ALTER TABLE guild_activity_roles
ADD COLUMN grant_type TEXT NOT NULL DEFAULT 'chat';

UPDATE guild_activity_roles
SET grant_type = 'chat';

ALTER TABLE guild_profiles
ALTER COLUMN last_grant_epoch SET DEFAULT 0;