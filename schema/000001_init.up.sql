CREATE TABLE users
(
    id serial not null unique,
    name varchar(255) not null,
    grade int not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE petitions_lists
(
    id serial not null unique,
    title varchar(255) not null,
    date timestamp not null default now(),
    timeEnd timestamp not null default now() + interval '14 day',
    text varchar(255) not null,
    answer varchar(255) not null default ''
);

CREATE TABLE users_lists
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    petition_id int references petitions_lists (id) on delete cascade not null
);

CREATE TABLE subs_items
(
    id serial not null unique,
    date timestamp not null default now()
);

CREATE TABLE petitions_items
(
    id serial not null unique,
    sub_id int references subs_items (id) on delete cascade not null,
    petition_id int references petitions_lists (id) on delete cascade not null,
    user_id int references users (id) on delete cascade not null unique
);
