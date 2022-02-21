CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "full_name" varchar(100) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "hash" varchar(255) NOT NULL,
  "salt" varchar(255) NOT NULL,
  "is_active" smallint NOT NULL DEFAULT 0,
  "is_trashed" smallint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "labels" (
  "id" serial PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "user_id" int NOT NULL,
  "is_trashed" smallint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "notes" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "title" varchar(255),
  "color" varchar(10),
  "type" varchar(4) NOT NULL DEFAULT 'note',
  "is_pinned" smallint NOT NULL DEFAULT 0,
  "is_archived" smallint NOT NULL DEFAULT 0,
  "is_trashed" smallint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL
);

CREATE TABLE "notes_items" (
  "id" serial PRIMARY KEY,
  "note_id" int NOT NULL,
  "text" varchar(1000) NOT NULL,
  "is_checked" smallint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notes_labels" (
  "note_id" int,
  "label_id" int,
  PRIMARY KEY ("note_id", "label_id")
);

ALTER TABLE "labels" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "notes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "notes_items" ADD FOREIGN KEY ("note_id") REFERENCES "notes" ("id");

ALTER TABLE "notes_labels" ADD FOREIGN KEY ("note_id") REFERENCES "notes" ("id");

ALTER TABLE "notes_labels" ADD FOREIGN KEY ("label_id") REFERENCES "labels" ("id");


COMMENT ON COLUMN "users"."hash" IS 'password hash';

COMMENT ON COLUMN "users"."salt" IS 'password salt';
