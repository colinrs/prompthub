drop table if exists `users_token_table`;
CREATE TABLE IF NOT EXISTS `users_token_table` (
                                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned ,
    `pre_efresh_token` varchar(256) default '',
    `refresh_token` varchar(256)  NOT NULL,
    `created_at` timestamp default CURRENT_TIMESTAMP,
    `updated_at` timestamp default CURRENT_TIMESTAMP,
    `deleted_at` timestamp DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;