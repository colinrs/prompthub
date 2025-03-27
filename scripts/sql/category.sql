drop table if exists `category_table`;
CREATE TABLE IF NOT EXISTS `category_table` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `category_name` varchar(60) NOT NULL,
                               `color` varchar(256) NOT NULL,
                               `category_status`  int(11) unsigned default 1,
                               `description` varchar(500)  default '',
                               `created_by` varchar(60) NOT NULL,
                               `created_at` timestamp default CURRENT_TIMESTAMP,
                               `updated_at` timestamp default CURRENT_TIMESTAMP,
                               `deleted_at` timestamp DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `idx_category_status` (`category_status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;