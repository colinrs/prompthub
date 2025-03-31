drop table if exists `users_like`;
CREATE TABLE IF NOT EXISTS `users_like` (
                                             `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                             `user_id` int(11) unsigned NOT NULL,
                                             `prompts_id` int(11) unsigned NOT NULL,
                                             `created_at` timestamp default CURRENT_TIMESTAMP,
                                             `updated_at` timestamp default CURRENT_TIMESTAMP,
                                             `deleted_at` timestamp DEFAULT NULL,
                                             PRIMARY KEY (`id`),
                                             UNIQUE KEY `idx_user_prompt` (`user_id`, `prompts_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (1, 1, NOW(), NOW());

INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (1, 2, NOW(), NOW());

INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (2, 2, '2023-06-12 09:15:00', '2023-06-12 09:15:00');

INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (3, 3, NOW(), NOW());

INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (4, 4, '2023-08-20 16:45:00', '2023-08-20 16:45:00', '2023-09-05 11:30:00');

INSERT INTO `users_like` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (5, 5, NOW(), NOW());
