drop table if exists `users_token_table`;
CREATE TABLE IF NOT EXISTS `users_token_table` (
                                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned ,
    `pre_refresh_token` varchar(256) default '',
    `refresh_token` varchar(256)  NOT NULL,
    `created_at` timestamp default CURRENT_TIMESTAMP,
    `updated_at` timestamp default CURRENT_TIMESTAMP,
    `deleted_at` timestamp DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `users_token_table` (`user_id`, `pre_refresh_token`, `refresh_token`, `created_at`, `updated_at`)
VALUES (1, '', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9', NOW(), NOW());

INSERT INTO `users_token_table` (`user_id`, `pre_refresh_token`, `refresh_token`, `created_at`, `updated_at`)
VALUES (2, 'old_token_123', 'new_token_456', '2023-01-15 10:00:00', '2023-01-15 10:00:00');

INSERT INTO `users_token_table` (`user_id`, `pre_refresh_token`, `refresh_token`, `created_at`, `updated_at`)
VALUES (3, 'expired_789', 'valid_abc123', NOW(), NOW());

INSERT INTO `users_token_table` (`user_id`, `pre_refresh_token`, `refresh_token`, `created_at`, `updated_at`, `deleted_at`)
VALUES (4, 'revoked_xyz', 'active_987', '2023-03-20 14:30:00', '2023-03-20 14:30:00', '2023-04-01 09:15:00');