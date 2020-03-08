CREATE DATABASE IF NOT EXISTS `gophr`;
USE `gophr`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`(
  `id` int(36) NOT NULL AUTO_INCREMENT,
  `username` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `email` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

LOCK TABLES `user` WRITE;
INSERT INTO `user` VALUES
  (1, 'luffy.monkey', 'luffy.monkey@gmail.com', 'secretpass', NULL, NULL, NULL);
UNLOCK TABLES;