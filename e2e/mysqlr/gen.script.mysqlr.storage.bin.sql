

USE `drwms`;

CREATE TABLE `oper_storage_bin` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '库位表主键',
	`type_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '库位类型ID',
	`code` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '库位代码（自动生成）',
	`state` INT(11) NOT NULL DEFAULT '0' COMMENT '库位状态',
	`warehouse_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '仓库ID',
	`region_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '库域ID',
	`area_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '库区ID',
	`aisle` INT(11) NOT NULL DEFAULT '0' COMMENT '巷道',
	`rack_row` INT(11) NOT NULL DEFAULT '0' COMMENT '货架行数',
	`rack_col` INT(11) NOT NULL DEFAULT '0' COMMENT '货架列数',
	`capacity_preempted` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '预占库容（单位：g）',
	`capacity_occupied` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '实占库容（单位：g）',
	`capacity_lock` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '锁定库容（单位：g）',
	`create_at` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
	`update_at` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '更新时间',
	`create_by` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '创建人',
	`update_by` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '更新人',
	PRIMARY KEY(`id`)
	KEY `update_at_of_storage_bin_idx` (`update_at`),
	KEY `type_id_of_storage_bin_idx` (`type_id`),
	KEY `rack_row_of_storage_bin_idx` (`rack_row`),
	KEY `rack_col_of_storage_bin_idx` (`rack_col`),
	KEY `code_of_storage_bin_idx` (`code`),
	KEY `aisle_of_storage_bin_idx` (`aisle`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '基础库位表';

