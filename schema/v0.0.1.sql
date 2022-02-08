

DROP TABLE IF EXISTS `permit_group`;
CREATE TABLE `permit_group` (
    `id` 		 VARCHAR(64) 	COMMENT 'unique id',
    `app_id` 	 VARCHAR(64) 	COMMENT 'app id',
    `name`          VARCHAR(64)     NOT NULL COMMENT 'permit group name',
    `description`   VARCHAR(255) COMMENT 'description',
    `created_at`     BIGINT(20) 	    COMMENT 'create time',
    `update_at`     BIGINT(20)  	COMMENT 'update time',
    `creator_id`    VARCHAR(36) COMMENT 'creator id',
    `creator_name`   VARCHAR(16) COMMENT 'creator name',
    `types`        TINYINT(1)   COMMENT 'types 1 init 2 create',

    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `permit_grant`;
CREATE TABLE `permit_grant` (
    `permit_id` 		 VARCHAR(64) 	COMMENT 'permit id',
    `owner` 	 VARCHAR(64) 	COMMENT 'owner id',
    `owner_name` VARCHAR(64)     NOT NULL COMMENT 'owner_nam',
    `types`        TINYINT(1)   COMMENT 'types 1 init 2 create',
    PRIMARY KEY (`permit_id`,`owner`,`types`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `permit_form`;
CREATE TABLE `permit_form` (
    `permit_id` 		 VARCHAR(64) 	COMMENT 'permit id',
    `form_id` 	 VARCHAR(64) 	COMMENT 'form id',
    `form_type`  TINYINT(1)    NOT NULL COMMENT 'form type',
    `authority`  INT   COMMENT 'authority' ,
    `conditions` 		 TEXT 	COMMENT 'conditions',
    `field_json` 	 TEXT	COMMENT 'field json',
    `web_schema` TEXT     COMMENT 'web schema',
    PRIMARY KEY (`permit_id`,`form_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

