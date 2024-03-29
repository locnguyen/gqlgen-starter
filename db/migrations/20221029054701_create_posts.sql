-- create "posts" table
CREATE TABLE "posts"
(
    "id"          bigint      NOT NULL GENERATED BY DEFAULT AS IDENTITY,
    "create_time" timestamptz NOT NULL,
    "update_time" timestamptz NOT NULL,
    "content"     text        NOT NULL,
    "author_id"   bigint      NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "posts_users_posts" FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE NO ACTION
);
-- modify "sessions" table
ALTER TABLE "sessions"
    ADD COLUMN "create_time" timestamptz NOT NULL,
    ADD COLUMN "update_time" timestamptz NOT NULL,
    ADD COLUMN "user_id"     bigint      NOT NULL,
    DROP CONSTRAINT "sessions_users_sessions",
    ADD CONSTRAINT "sessions_users_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION;
