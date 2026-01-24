CREATE OR REPLACE FUNCTION insert_guild_settings()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO guild_voice_activity_settings (guild_id)
    VALUES (NEW.guild_id);

    INSERT INTO guild_chat_activity_settings (guild_id)
    VALUES (NEW.guild_id);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TABLE guild_message_embeds_settings;