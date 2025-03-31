drop table if exists `prompts_table`;
CREATE TABLE IF NOT EXISTS `prompts_table` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `title` varchar(60) NOT NULL,
                               `content` text NOT NULL,
                               `category` int(11) unsigned NOT NULL,
                               `prompts_status`  int(11) unsigned default 1,
                               `created_by` int(11) unsigned NOT NULL,
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
                                 UNIQUE `idx_prompts_id` (`prompts_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- prompts_table 数据
INSERT INTO `prompts_table` (`title`, `content`, `category`, `prompts_status`, `created_by`)
VALUES ('高效写作助手', '帮助您快速生成高质量文章的提示词', 1, 1, 1);

INSERT INTO `prompts_table` (`title`, `content`, `category`, `prompts_status`, `created_by`)
VALUES ('代码调试专家', '解决编程问题的AI助手提示词', 2, 1, 2);

INSERT INTO `prompts_table` (`title`, `content`, `category`, `prompts_status`, `created_by`, `deleted_at`)
VALUES ('旧版营销文案', '已下架的营销文案提示词', 3, 1, 3, '2023-11-01 00:00:00');

INSERT INTO `prompts_table` (`title`, `content`, `category`, `prompts_status`, `created_by`)
VALUES ('学习计划生成器', '帮助学生制定学习计划的提示词', 1, 2, 4);

-- prompts_count_table 对应数据
INSERT INTO `prompts_count_table` (`prompts_id`, `like_count`, `review_count`)
VALUES (1, 156, 28);

INSERT INTO `prompts_count_table` (`prompts_id`, `like_count`, `review_count`)
VALUES (2, 89, 15);

INSERT INTO `prompts_count_table` (`prompts_id`, `like_count`, `review_count`)
VALUES (4, 42, 7);
