CREATE TABLE IF NOT EXISTS "reports" (
    "id" BIGINT PRIMARY KEY,
    -- optional. if the config var require_user_to_be_logged_in_to_submit_reports is true then
    -- it can be assumed that this will be here.
    "reporter" BIGINT REFERENCES "accounts"("id"),
    "created_on" TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    "target_steam_link" VARCHAR NOT NULL,
    "insider_sid64" BIGINT NOT NULL,
    "insider_sid32" VARCHAR NOT NULL,
    "insider_sid3" VARCHAR NOT NULL,
    "insider_sid" VARCHAR NOT NULL,
    "server_reported_on" VARCHAR NOT NULL,
    "server_reported_on_ip" INET,
    "accepted_by" BIGINT REFERENCES "accounts"("id"),
    "accepted_on" TIMESTAMP WITHOUT TIME ZONE,
    "flags" BIGINT NOT NULL DEFAULT 0
);