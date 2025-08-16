-- +goose Up

create type follow_feed_result as (
    id uuid,
    created_at timestamp,
    updated_at timestamp,
    user_name text,
    feed_name text 
);

create type feed_follows_for_users_result as (
    user_name text,
    feed_name text,
    feed_url text
);

-- +goose StatementBegin
create or replace function follow_feed(uuid, timestamp, timestamp, text, text)
returns setof follow_feed_result as $$
declare
    userID uuid;
    feedID uuid;
    followID uuid;
begin
    select users.id into userID from users where name = $4;

    if userID is null then
        raise exception 'user with username % does not exist', $4;
    end if;

    select feeds.id into feedID from feeds where url = $5;

    if feedID is null then
        raise exception 'feed with url % does not exist', $5;
    end if;

    insert into feed_follows (id, created_at, updated_at, user_id, feed_id) values
    ($1, $2, $3, userID, feedID) on conflict(user_id, feed_id) do nothing returning
    id into followID;

    if followID is null then 
        raise exception 'you already follow this feed';
    end if;

    return query 
        with created_feed_follow as (
            select * from feed_follows where id = followID
        ) select cff.id, cff.created_at, cff.updated_at, users.name, feeds.url from created_feed_follow as cff
        inner join users on cff.user_id = users.id
        inner join feeds on cff.feed_id = feeds.id;
end;
$$
language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
create function feed_follows_for_users(text)
returns setof feed_follows_for_users_result as $$
declare
    userID uuid;
begin

    select users.id into userID from users where users.name = $1;

    if userID is null then
        raise exception 'user with username % does not exist', $1;
    end if;

    return query 
        select users.name as user_name, feeds.name as feed_name, feeds.url as feed_url from feed_follows 
    inner join users on users.id = feed_follows.user_id
    inner join feeds on feeds.id = feed_follows.feed_id where users.id = userID;
end;
$$
language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
create function unfollow_feed(text, text)
returns void as $$
declare
    userID uuid;
    feedID uuid;
begin
    select users.id into userID from users where users.name = $2;

    if userID is null then
        raise exception 'user with username % does not exists', $2;
    end if;
    
    select feeds.id into feedID from feeds where feeds.url = $1;

    if feedID is null then
        raise exception 'feed with url % does not exist', $1;
    end if;

    delete from feed_follows where user_id = userID and feed_id = feedID;
end;
$$
language plpgsql;
-- +goose StatementEnd

-- +goose Down
drop function follow_feed(uuid, timestamp, timestamp, text, text);
drop function feed_follows_for_users(text);
drop function unfollow_feed(text, text);
drop type follow_feed_result;
drop type feed_follows_for_users_result;
