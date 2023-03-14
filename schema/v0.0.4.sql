
alter table project
    add `serial_number` varchar(32) comment '项目编号',
    add `start_at` BIGINT(20) 	 COMMENT '项目开始日期',
    add `end_at` BIGINT(20) comment '项目结束日期',
    add `status` varchar(32)  comment '状态',
    add `remark` varchar(64) comment '备注';
