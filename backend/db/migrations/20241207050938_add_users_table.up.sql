CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL,
    "is_admin" boolean NOT NULL DEFAULT false
)