-- Create "posts" table
CREATE TABLE "posts"
(
    "id"          bigint      NOT NULL DEFAULT (generate_id()),
    "create_time" timestamptz NOT NULL,
    "update_time" timestamptz NOT NULL,
    "content"     text        NOT NULL,
    "author_id"   bigint      NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "posts_users_posts" FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE NO ACTION
);
