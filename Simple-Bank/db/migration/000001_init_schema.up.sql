CREATE TABLE "users" (
    "iin" bigserial PRIMARY KEY,
    "username" varchar UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "name" varchar NOT NULL,
    "surname" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "user_role" varchar NOT NULL DEFAULT 'user'
);

