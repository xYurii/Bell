SET statement_timeout = 0;

CREATE TABLE raffles (
    id SERIAL PRIMARY KEY,
    raffle_type INT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ends_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ NULL,
    winner_ticket_id BIGINT NULL
);

CREATE TABLE raffle_tickets (
    id SERIAL PRIMARY KEY,
    bought_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL
);

CREATE INDEX raffle_tickets_idx ON raffle_tickets (user_id, bought_at);
CREATE INDEX raffles_idx ON raffles (id, ended_at, ends_at, started_at, raffle_type);