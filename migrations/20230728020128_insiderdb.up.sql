CREATE TABLE IF NOT EXISTS "accounts" (
    "id" BIGINT PRIMARY KEY,
    "steam" BIGINT NOT NULL,
    "discord" BIGINT,
    "steam_name" VARCHAR NOT NULL,
    "discord_name" VARCHAR,
    "created_on" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    "flags" BIGINT NOT NULL DEFAULT 0,
    "accepted_by" BIGINT NOT NULL,
    "accepted_on" TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS "clans" (
    "id" BIGINT PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "marker" VARCHAR,
    -- clans can be unclaimed.
    "owner" BIGINT REFERENCES "accounts"("id"),
    "added_by" BIGINT NOT NULL REFERENCES "accounts"("id"),
    "created_on" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS "clan_members" (
    "id" BIGINT PRIMARY KEY,
    "clan" BIGINT REFERENCES "clans"("id"),
    "member" BIGINT REFERENCES "accounts"("id"),
    "added_on" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    "owner" BOOLEAN NOT NULL DEFAULT FALSE
);