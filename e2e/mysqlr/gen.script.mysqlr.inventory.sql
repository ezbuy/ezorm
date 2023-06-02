

USE `drwms`;

CREATE TABLE `oper_inventory` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '库存明细表主键',
	`warehouse_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '仓库ID',
	`bin_id` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '库位表',
	`sku_code` VARCHAR(500) NOT NULL DEFAULT '' COMMENT 'sku号',
	`barcode` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '条形码',
	`quantity_total` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '总库存',
	`quantity_available` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '可用库存',
	`quantity_locked` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '锁定库存',
	`create_at` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '创建时间',
	`update_at` BIGINT(20) NOT NULL DEFAULT '0' COMMENT '更新时间',
	`create_by` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '创建人',
	`update_by` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '更新人',
	PRIMARY KEY(`id`)
	KEY `warehouse_id_bin_id_barcode_of_inventory_idx` (`warehouse_idbin_idbarcode`),
	KEY `update_at_of_inventory_idx` (`update_at`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '库存明细表';

