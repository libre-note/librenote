CREATE TABLE `users` (
  `id` INTEGER PRIMARY KEY,
  `full_name` TEXT NOT NULL,
  `email` TEXT NOT NULL,
  `hash` TEXT NOT NULL,
  `salt` TEXT NOT NULL,
  `is_active` INTEGER NOT NULL DEFAULT 0,
  `is_trashed` INTEGER NOT NULL DEFAULT 0,
  `created_at` TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TEXT NOT NULL
);
