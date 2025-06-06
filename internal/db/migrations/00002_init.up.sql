ALTER TABLE guild_activity_roles
ADD COLUMN grant_type TEXT;

UPDATE guild_activity_roles
SET grant_type = 'chat';