-- +goose Up
-- +goose StatementBegin
create function add_feed(uuid, timestamp, timestamp, text, text, text)
returns uuid as $$
declare
    userID uuid;
    feedID uuid;
    resultID uuid;
begin
    select id into userID from users where name = $6;

    if userID is null then
        raise exception 'User with username % not found', $6;
    end if;

    insert into feeds (id, created_at, updated_at, name, url, user_id) values
    ($1, $2, $3, $4, $5, userID) on conflict(url) do nothing returning id into feedID;

    if feedID is null then
        raise exception 'Feed with url % already exists', $5;
    end if;

    insert into feed_follows (id, created_at, updated_at, user_id, feed_id) values (
        gen_random_uuid(), now(), now(), userID, feedID  
    ) on conflict(user_id, feed_id) do nothing returning id into resultID;

    if resultID is null then
        raise exception 'User % already follows this feed', $6;
    end if;

    return resultID;
end;
$$
language plpgsql;
-- +goose StatementEnd

-- +goose Down
drop function add_feed(uuid, timestamp, timestamp, text, text, text);
