CREATE TYPE pair_status AS ENUM('halted', 'trading');

CREATE TYPE order_side AS ENUM('bid', 'ask');

CREATE TABLE
    orders (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        pair_base TEXT NOT NULL,
        pair_quote TEXT NOT NULL,
        side order_side NOT NULL,
        price BIGINT NOT NULL,
        amount BIGINT NOT NULL,
        remaining BIGINT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE
    outbox (
        id UUID PRIMARY KEY,
        aggregate_id UUID NOT NULL,
        aggregate_type TEXT NOT NULL,
        event_type TEXT NOT NULL,
        payload JSONB NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        is_commited BOOLEAN NOT NULL DEFAULT false
    )