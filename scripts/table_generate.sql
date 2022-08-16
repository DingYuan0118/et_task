CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `usr_name` varchar(64) NOT NULL DEFAULT '',
  `usr_nickname` varchar(64) NOT NULL DEFAULT '',
  `usr_password` varchar(64) NOT NULL DEFAULT '',
  `profile_pic_url` varchar(1024) NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modified time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_user_usr_name` (`usr_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8