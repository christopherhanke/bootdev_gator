-- +goose Up
create table feed_follows(
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    user_id uuid not null references users on delete cascade,
    feed_id uuid not null references feeds on delete cascade,
    foreign key (user_id) references users(id),
    foreign key (feed_id) references feeds(id),
    unique (user_id, feed_id)
);

-- +goose Down
drop table feed_follows;