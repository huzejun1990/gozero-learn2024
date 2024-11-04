CREATE TABLE `user`  (
                         `id` BIGINT NOT NULL AUTO_INCREMENT,
                         `username` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                         `password` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                         `register_time` DATETIME NOT NULL,
                         `last_login_time` DATETIME NOT NULL,
                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = INNODB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;