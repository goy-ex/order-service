CREATE TYPE pair_status AS ENUM('active', 'inactive');

CREATE TYPE order_side AS ENUM('bid', 'ask');

CREATE TABLE orders (
    id             UUID        PRIMARY KEY,
    user_id        UUID        NOT NULL,
    pair_base      TEXT        NOT NULL,
    pair_quote     TEXT        NOT NULL,
    side           order_side  NOT NULL,
    price          BIGINT      NOT NULL,
    qty            BIGINT      NOT NULL,
    remaining_qty  BIGINT      NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL
);

CREATE TYPE aggregate_type AS ENUM('order');
CREATE TYPE event_type AS ENUM('order_created', 'order_cancelled');
CREATE TYPE event_status AS ENUM('pending', 'resolved');

CREATE TABLE events (
    id             UUID           PRIMARY KEY,
    aggregate_id   UUID           NOT NULL,
    aggregate_type aggregate_type NOT NULL,
    partition_key  TEXT,
    event_type     event_type     NOT NULL,
    payload        JSONB,
    created_at     TIMESTAMPTZ    NOT NULL,
    "status"       event_status   NOT NULL DEFAULT 'pending'
);
