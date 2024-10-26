SET statement_timeout = 0;

CREATE TABLE raffles (
    id SERIAL PRIMARY KEY,
    raffle_type INT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ends_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ NULL
);

CREATE TABLE raffle_tickets (
    id SERIAL PRIMARY KEY,
    bought_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(50) NOT NULL,
    raffle_id INT NOT NULL REFERENCES raffles (id) ON DELETE CASCADE
);

CREATE INDEX raffle_tickets_idx ON raffle_tickets (user_id, bought_at, raffle_id);
CREATE INDEX raffles_idx ON raffles (ended_at, ends_at, started_at, raffle_type);

ALTER TABLE raffles 
    ADD COLUMN winner_ticket_id INT REFERENCES raffle_tickets (id) ON DELETE SET NULL;