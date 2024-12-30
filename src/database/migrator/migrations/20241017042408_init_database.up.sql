SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id VARCHAR(40) UNIQUE PRIMARY KEY,
  money BIGINT NOT NULL DEFAULT 0,
  language VARCHAR(6) NOT NULL DEFAULT 'pt-BR',
  status_time BIGINT DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guilds (
  id VARCHAR(40) UNIQUE PRIMARY KEY,
  prefix VARCHAR(6) DEFAULT '..',
  commands_channels TEXT[],
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
CREATE INDEX guilds_idx ON guilds (id);
CREATE INDEX users_idx ON users (id, money);