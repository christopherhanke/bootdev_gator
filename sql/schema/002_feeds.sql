-- +goose Up
create table feeds(
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    name text not null,
    url text unique not null,
    user_id uuid not null references users on delete cascade,
    foreign key (user_id) references users(id)
);

alter table feeds add column last_fetched_at timestamp;

-- +goose Down
drop table feeds;