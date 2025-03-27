drop table if exists `users_table`;
CREATE TABLE IF NOT EXISTS `users_table` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `user_name` varchar(60) NOT NULL,
                               `password` varchar(256) NOT NULL,
                               `email` varchar(100)  NOT NULL,
                               `user_status`  int(11) unsigned default 1,
                               `description` varchar(500)  default '',
                               `created_at` timestamp default CURRENT_TIMESTAMP,
                               `updated_at` timestamp default CURRENT_TIMESTAMP,
                               `deleted_at` timestamp DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `idx_email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;