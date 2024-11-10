BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "categories" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "name" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "posts" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER,
    "title" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "created_date" TIMESTAMP NOT NULL,
    "amount_of_comments" INTEGER NOT NULL DEFAULT 0,
    "amount_of_watches" INTEGER NOT NULL DEFAULT 0,
    "rate" INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "comments" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "post_id" INTEGER,
    "user_id" INTEGER,
    "content" TEXT NOT NULL,
    "created_date" TIMESTAMP NOT NULL,
    "rate" INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY("post_id") REFERENCES "posts"("id"),
    FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "post_categories" (
    "post_id" INTEGER,
    "category_id" INTEGER,
    PRIMARY KEY("post_id", "category_id"),
    FOREIGN KEY("post_id") REFERENCES "posts"("id"),
    FOREIGN KEY("category_id") REFERENCES "categories"("id")
);
CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "id_github" INTEGER UNIQUE,
    "id_google" INTEGER UNIQUE,
    "username" TEXT UNIQUE NOT NULL,
    "email" TEXT UNIQUE,
    "password_hash" TEXT,
    "created_at" TIMESTAMP NOT NULL,
    "amount_of_posts" INTEGER NOT NULL DEFAULT 0,
    "amount_of_comments" INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS "active_sessions" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "session_id" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    FOREIGN KEY("user_id") REFERENCES "users"("id")
);
CREATE TABLE IF NOT EXISTS "rates" (
    "user_id" INTEGER,
    "post_id" INTEGER,
    "comment_id" INTEGER,
    "status" string NOT NULL,
    "rated_at" TIMESTAMP NOT NULL,
    FOREIGN KEY("user_id") REFERENCES "users"("id"),
    FOREIGN KEY("post_id") REFERENCES "posts"("id"),
    FOREIGN KEY("comment_id") REFERENCES "comments"("id")
);
COMMIT;