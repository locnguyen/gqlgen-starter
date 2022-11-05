-- modify "sessions" table
ALTER TABLE "sessions" ADD COLUMN "token" character varying NOT NULL, ADD COLUMN "data" bytea NOT NULL, DROP CONSTRAINT "sessions_users_sessions";
ALTER TABLE "sessions" DROP COLUMN "sid";
ALTER TABLE "sessions" DROP COLUMN "type";
ALTER TABLE "sessions" DROP COLUMN "create_time";
ALTER TABLE "sessions" DROP COLUMN "update_time";
ALTER TABLE "sessions" DROP COLUMN "user_id";
-- create index "sessions_token_key" to table: "sessions"
CREATE UNIQUE INDEX "sessions_token_key" ON "sessions" ("token");
-- create index "session_expiry" to table: "sessions"
CREATE INDEX "session_expiry" ON "sessions" ("expiry");
