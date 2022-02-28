CREATE TABLE `users` (
  `id` INTEGER NOT NULL,
  `full_name` TEXT NOT NULL,
  `email` TEXT NOT NULL,
  `hash` TEXT NOT NULL,
  `salt` TEXT NOT NULL,
  `is_active` INTEGER NOT NULL DEFAULT 0,
  `is_trashed` INTEGER NOT NULL DEFAULT 0,
  `list_view_enabled` INTEGER NOT NULL DEFAULT 0,
  `dark_mode_enabled` INTEGER NOT NULL DEFAULT 0,
  `created_at` TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` TEXT NOT NULL,
  CONSTRAINT user_PK PRIMARY KEY(id),
  CONSTRAINT user_email_UNIQUE UNIQUE(email)
);

CREATE TABLE `labels` (
  `id` INTEGER NOT NULL,
  `name` TEXT NOT NULL,
  `user_id` INTEGER NOT NULL,
  `is_trashed` INTEGER NOT NULL DEFAULT 0,
  `created_at` TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` TEXT NOT NULL,
   CONSTRAINT label_PK PRIMARY KEY(id),
   CONSTRAINT user_id_FK FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE `notes` (
  `id` INTEGER NOT NULL,
  `user_id` INTEGER NOT NULL,
  `title` TEXT,
  `color` TEXT,
  `type` TEXT NOT NULL DEFAULT "note",
  `is_pinned` INTEGER NOT NULL DEFAULT 0,
  `is_archived` INTEGER NOT NULL DEFAULT 0,
  `is_trashed` INTEGER NOT NULL DEFAULT 0,
  `created_at` TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` TEXT NOT NULL,
  CONSTRAINT note_PK PRIMARY KEY(id)
  CONSTRAINT user_id_FK FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE `notes_items` (
  `id` INTEGER NOT NULL,
  `note_id` INTEGER NOT NULL,
  `text` TEXT NOT NULL,
  `is_checked` INTEGER NOT NULL DEFAULT 0,
  `created_at` TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
   CONSTRAINT notes_item_PK PRIMARY KEY(id),
   CONSTRAINT note_id_FK FOREIGN KEY(note_id) REFERENCES notes(id)
);

CREATE TABLE `notes_labels` (
  `note_id` INTEGER NOT NULL,
  `label_id` INTEGER NOT NULL,
  CONSTRAINT notes_labels_PK PRIMARY KEY(note_id, label_id),
  CONSTRAINT note_id_FK FOREIGN KEY(note_id) REFERENCES notes(id),
  CONSTRAINT label_id_FK FOREIGN KEY(label_id) REFERENCES labels(id)
);
