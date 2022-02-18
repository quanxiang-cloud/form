

DROP TABLE IF EXISTS `permit_group`;
CREATE TABLE `permit_group` (
    `id` 		 VARCHAR(64) 	COMMENT 'unique id',
    `app_id` 	 VARCHAR(64) 	NOT NULL COMMENT 'app id',
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
    `id` 		 VARCHAR(64) 	COMMENT 'unique id',
    `permit_id` 		 VARCHAR(64) NOT NULL 	COMMENT 'permit id',
    `owner` 	 VARCHAR(64) 	NOT NULL COMMENT 'owner id',
    `owner_name` VARCHAR(64)     NOT NULL COMMENT 'owner_nam',
    `types`        TINYINT(1)   COMMENT 'types 1 init 2 create',
    UNIQUE KEY `idx_global_name` (`permit_id`, `owner`,`types`),
    PRIMARY KEY  (`id`)

)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `permit_form`;
CREATE TABLE `permit_form` (
    `id` 		 VARCHAR(64) 	COMMENT 'unique id',
    `permit_id`  VARCHAR(64) 	 NOT NULL COMMENT 'permit id',
    `form_id` 	 VARCHAR(64) 	 NOT NULL COMMENT 'form id',
    `form_type`  TINYINT(1)    NOT NULL COMMENT 'form type',
    `authority`  INT   COMMENT 'authority' ,
    `conditions` 		 TEXT 	COMMENT 'conditions',
    `field_json` 	 TEXT	COMMENT 'field json',
    `web_schema` TEXT     COMMENT 'web schema',
    UNIQUE KEY `idx_global_name` (`permit_id`, `form_id`),
    PRIMARY KEY  (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `table`;

CREATE TABLE `table` (
   `id` 		 VARCHAR(64) 	COMMENT ' id',
   `app_id` 	 VARCHAR(64) 	NOT NULL COMMENT 'table is which app',
   `table_id`    VARCHAR(64)    NOT NULL COMMENT 'tableID',
   `schema`      TEXT   COMMENT 'web schema' ,
   `config` 	 TEXT 	COMMENT 'config',
   UNIQUE KEY `idx_global_name` (`app_id`, `table_id`),
   PRIMARY KEY  (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `table_schema`;

CREATE TABLE `table_schema` (
   `id` 		 VARCHAR(64) 	COMMENT ' id',
   `app_id` 	 VARCHAR(64) 	NOT NULL COMMENT 'table is which app',
   `table_id`    VARCHAR(64)    NOT NULL COMMENT 'table id',
   `title`       VARCHAR(32)     COMMENT 'title',
   `field_len`   INT             COMMENT 'field_len' ,
   `description` VARCHAR(100)    COMMENT 'description',
   `source`      TINYINT(1)      COMMENT 'source',
   `created_at`  BIGINT(20) 	 COMMENT 'create time',
   `update_at`   BIGINT(20) 	 COMMENT 'update time',
   `creator_id`    VARCHAR(36) COMMENT 'creator id',
   `creator_name`   VARCHAR(16) COMMENT 'creator name',
   `editor_id`    VARCHAR(36) COMMENT 'editor id',
   `editor_name`   VARCHAR(16) COMMENT 'editor name',
   `schema`      TEXT   COMMENT 'web schema' ,
   UNIQUE KEY `idx_global_name` (`app_id`, `table_id`),
   PRIMARY KEY  (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

