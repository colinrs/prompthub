drop table if exists `users_save`;
CREATE TABLE IF NOT EXISTS `users_save` (
                                            `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned NOT NULL,
    `prompts_id` int(11) unsigned NOT NULL,
    `created_at` timestamp default CURRENT_TIMESTAMP,
    `updated_at` timestamp default CURRENT_TIMESTAMP,
    `deleted_at` timestamp DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_prompt` (`user_id`, `prompts_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (1, 1, NOW(), NOW());

INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (2, 2, '2023-05-10 08:30:00', '2023-05-10 08:30:00');



INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (2, 1, '2023-05-10 08:30:00', '2023-05-10 08:30:00');

INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (3, 3, NOW(), NOW());

INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`, `deleted_at`)
VALUES (4, 4, '2023-07-15 14:20:00', '2023-07-15 14:20:00', '2023-08-01 10:00:00');

INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (5, 5, NOW(), NOW());


INSERT INTO `users_save` (`user_id`, `prompts_id`, `created_at`, `updated_at`)
VALUES (5, 4, NOW(), NOW());