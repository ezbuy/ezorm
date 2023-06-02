SELECT
	oper_inventory.id AS iid,
	oper_inventory.warehouse_id,
	oper_inventory.sku_code,
	oper_inventory.barcode,
	oper_inventory.quantity_total,
	oper_storage_bin.id AS bid,
	oper_storage_bin.area_id,
	oper_storage_bin.code
FROM
	oper_inventory AS i
	INNER JOIN oper_storage_bin AS b ON i.bin_id = b.id
WHERE
	oper_inventory.sku_code = '1234567890';
