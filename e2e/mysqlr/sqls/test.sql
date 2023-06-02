SELECT
	i.id AS iid,
	i.warehouse_id,
	i.sku_code,
	i.barcode,
	i.quantity_total,
	b.id as bid,
	b.area_id,
	b.code
FROM
	oper_inventory AS i
	INNER JOIN oper_storage_bin AS b ON i.bin_id = b.id
WHERE
	i.sku_code = 'xyz';
