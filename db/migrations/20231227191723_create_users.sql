-- Create "users" table
CREATE TABLE "users"
(
    "id"              bigint            NOT NULL DEFAULT (generate_id()),
    "create_time"     timestamptz       NOT NULL,
    "update_time"     timestamptz       NOT NULL,
    "email"           character varying NOT NULL,
    "hashed_password" bytea             NOT NULL,
    "first_name"      character varying NOT NULL,
    "last_name"       character varying NOT NULL,
    "phone_number"    character varying NOT NULL,
    "roles"           jsonb             NOT NULL,
    PRIMARY KEY ("id")
);
-- Create index "user_email" to table: "users"
CREATE UNIQUE INDEX "user_email" ON "users" ("email");
