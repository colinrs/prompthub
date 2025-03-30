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

-- 活跃分类
INSERT INTO `category_table` (`category_name`, `color`, `category_status`, `description`, `created_by`)
VALUES ('写作助手', '#FF6B6B', 1, '帮助提升写作效率的分类', 'admin');

INSERT INTO `category_table` (`category_name`, `color`, `category_status`, `description`, `created_by`)
VALUES ('编程开发', '#4ECDC4', 1, '编程相关的提示词分类', 'admin');

-- 非活跃分类
INSERT INTO `category_table` (`category_name`, `color`, `category_status`, `description`, `created_by`)
VALUES ('营销文案', '#FFE66D', 0, '已下架的营销类提示词', 'editor');

-- 已删除分类
INSERT INTO `category_table` (`category_name`, `color`, `category_status`, `description`, `created_by`, `deleted_at`)
VALUES ('生活技巧', '#A5FFD6', 1, '已删除的生活类提示词', 'editor', '2023-12-01 00:00:00');

-- 默认状态分类
INSERT INTO `category_table` (`category_name`, `color`, `description`, `created_by`)
VALUES ('学习辅导', '#B388FF', '教育学习类提示词', 'teacher');
