drop table if exists `users_table`;
CREATE TABLE IF NOT EXISTS `users_table` (
                               `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                               `user_name` varchar(60) NOT NULL,
                               `password` varchar(256) NOT NULL,
                               `email` varchar(100)  NOT NULL,
                               `user_status`  int(11) unsigned default 1,
                               `avatar` varchar(256)  default '',
                               `description` varchar(500)  default '',
                               `created_at` timestamp default CURRENT_TIMESTAMP,
                               `updated_at` timestamp default CURRENT_TIMESTAMP,
                               `deleted_at` timestamp DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `idx_email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 基本活跃用户
INSERT INTO `users_table` (`user_name`, `password`, `email`, `user_status`, `avatar`, `description`)
VALUES ('user1', '$2a$10$xJwL5v5zP6zE2NQ1F3YQZ.9Zq1lWb6nX1c', 'user1@example.com', 1, 'avatar1.jpg', '我是用户1');

-- 非活跃用户
INSERT INTO `users_table` (`user_name`, `password`, `email`, `user_status`, `avatar`, `description`)
VALUES ('user2', '$2a$10$yKv7r8s9t0u1v2w3x4y5z', 'user2@test.com', 0, '', '已停用账号');

-- 管理员用户
INSERT INTO `users_table` (`user_name`, `password`, `email`, `user_status`, `avatar`, `description`)
VALUES ('admin', '$2a$10$zA1b2c3d4e5f6g7h8i9j0', 'admin@prompthub.com', 2, 'admin.png', '系统管理员');

-- 已删除用户
INSERT INTO `users_table` (`user_name`, `password`, `email`, `user_status`, `avatar`, `description`, `deleted_at`)
VALUES ('old_user', '$2a$10$q1w2e3r4t5y6u7i8o9p0', 'old@user.com', 1, '', '已删除账号', '2023-10-01 00:00:00');

-- 新注册用户(使用默认值)
INSERT INTO `users_table` (`user_name`, `password`, `email`)
VALUES ('new_user', '$2a$10$a1s2d3f4g5h6j7k8l9z0', 'new@user.com');
