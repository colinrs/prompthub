drop table if exists `prompts_table`;
CREATE TABLE IF NOT EXISTS `prompts_table` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `title` varchar(60) NOT NULL,
                               `content` text NOT NULL,
                               `category` varchar(100)  NOT NULL,
                               `prompts_status`  int(11) unsigned default 1,
                               `created_by` varchar(60) NOT NULL,
                               `created_at` timestamp default CURRENT_TIMESTAMP,
                               `updated_at` timestamp default CURRENT_TIMESTAMP,
                               `deleted_at` timestamp DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `idx_prompts_status` (`prompts_status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

drop table if exists `prompts_count_table`;
CREATE TABLE IF NOT EXISTS `prompts_count_table` (
                                 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                 `prompts_id` int(11) unsigned  NOT NULL,
                                 `like_count` int(64) unsigned  NOT NULL,
                                 `review_count` int(64) unsigned  NOT NULL,
                                 `created_at` timestamp default CURRENT_TIMESTAMP,
                                 `updated_at` timestamp default CURRENT_TIMESTAMP,
                                 `deleted_at` timestamp DEFAULT NULL,
                                 PRIMARY KEY (`id`),
                                 KEY `idx_prompts_id` (`prompts_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;