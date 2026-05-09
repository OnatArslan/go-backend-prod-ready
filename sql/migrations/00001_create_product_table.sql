-- +goose Up
create table if not exists product (
    id uuid primary key default uuidv7(),
    name text not null,
    price_in_cents bigint not null check (price_in_cents >= 0),
    quantity integer not null default 0,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
-- +goose Down
drop table if exists product;
