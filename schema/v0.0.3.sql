
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
        `id` 		 VARCHAR(64) 	COMMENT 'unique id',
        `name` 	 VARCHAR(64) 	NOT NULL COMMENT 'name ',
        `description`   VARCHAR(255) COMMENT 'description',
        `created_at`     BIGINT(20) 	    COMMENT 'create time',
        `creator_id`    VARCHAR(36) COMMENT 'creator id',
        `creator_name`   VARCHAR(16) COMMENT 'creator name',
        PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `project_user`;
CREATE TABLE `project_user` (
        `id` 		 VARCHAR(64) 	COMMENT 'unique id',
       `project_id` 	 VARCHAR(64) 	NOT NULL COMMENT 'project id',
       `project_name`   VARCHAR(255) COMMENT 'project',
       `user_id` 	 VARCHAR(64) 	NOT NULL COMMENT 'user id',
       `user_name`    VARCHAR(36) COMMENT 'user name',
        `created_at`     BIGINT(20) 	    COMMENT 'create time',
        `creator_id`    VARCHAR(36) COMMENT 'creator id',
        `creator_name`   VARCHAR(16) COMMENT 'creator name',
       UNIQUE KEY `idx_project_user` (`project_id`, `user_id`),
       PRIMARY KEY  (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;