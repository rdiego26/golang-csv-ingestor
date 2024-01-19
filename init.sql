create table if not exists users
(
    id         uuid                                   not null primary key,
    first_name text                                   not null,
    last_name  text                                   not null,
    email      text unique                            not null,
    parent_id  uuid,
    created_at timestamp with time zone default now() not null,
    deleted_at timestamp with time zone,
    merged_at  timestamp with time zone
);

create index if not exists users_parent_id_index
    on users (parent_id);


