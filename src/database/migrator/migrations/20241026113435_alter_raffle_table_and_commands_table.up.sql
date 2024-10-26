ALTER TABLE raffles
    ADD COLUMN reward_price BIGINT NOT NULL DEFAULT 0;

CREATE TABLE commands_used (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    guild_id VARCHAR(50) NOT NULL,
    command_name VARCHAR(30) NOT NULL,
    message_content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX commands_used_idx ON commands_used (user_id, channel_id, guild_id, command_name);