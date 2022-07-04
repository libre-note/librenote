CREATE TABLE `users` (
  `id` int PRIMARY KEY,
  `full_name` varchar(100) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `hash` varchar(255) NOT NULL COMMENT 'password hash',
  `is_active` tinyint(1) NOT NULL DEFAULT 0,
  `is_trashed` tinyint(1) NOT NULL DEFAULT 0,
  `list_view_enabled` tinyint(1) NOT NULL DEFAULT 0,
  `dark_mode_enabled` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL
);

CREATE TABLE `labels` (
  `id` int PRIMARY KEY,
  `name` varchar(50) NOT NULL,
  `user_id` int NOT NULL,
  `is_trashed` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL
);

CREATE TABLE `notes` (
  `id` int PRIMARY KEY,
  `user_id` int NOT NULL,
  `title` varchar(255),
  `color` varchar(10),
  `type` varchar(4) NOT NULL DEFAULT "note",
  `is_pinned` tinyint(1) NOT NULL DEFAULT 0,
  `is_archived` tinyint(1) NOT NULL DEFAULT 0,
  `is_trashed` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL
);

CREATE TABLE `notes_items` (
  `id` int PRIMARY KEY,
  `note_id` int NOT NULL,
  `text` varchar(1000) NOT NULL,
  `is_checked` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL
);

CREATE TABLE `notes_labels` (
  `note_id` int,
  `label_id` int,
  PRIMARY KEY (`note_id`, `label_id`)
);

ALTER TABLE `labels` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `notes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `notes_items` ADD FOREIGN KEY (`note_id`) REFERENCES `notes` (`id`);

ALTER TABLE `notes_labels` ADD FOREIGN KEY (`note_id`) REFERENCES `notes` (`id`);

ALTER TABLE `notes_labels` ADD FOREIGN KEY (`label_id`) REFERENCES `labels` (`id`);
