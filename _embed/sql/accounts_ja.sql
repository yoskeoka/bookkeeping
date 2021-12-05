-- SQLite3

insert into accounts(code, name, is_bs, is_left)
values
/* B/S科目 */
-- 資産
(1110, '現金及び預金', TRUE, TRUE),
(1120, '売掛金', TRUE, TRUE),
(1130, '商品', TRUE, TRUE),
(1210, '有形固定資産', TRUE, TRUE),
(1211, '機械装置', TRUE, TRUE),
-- 負債
(2100, '買掛金', TRUE, FALSE),
(2101, '短期借入金', TRUE, FALSE),
(2102, '未払い法人税等', TRUE, FALSE),
(2103, '預り金', TRUE, FALSE),
(2200, '長期借入金', TRUE, FALSE),
-- 純資産
(3100, '資本金', TRUE, FALSE),
(3200, '資本剰余金', TRUE, FALSE),

/* P/L科目 */
-- 売上高
(4100, '商品売上高', FALSE, FALSE),
-- 売上原価
(5100, '期首商品棚卸高', FALSE, TRUE),
(5200, '商品仕入高', FALSE, TRUE),
(5300, '期末商品棚卸高', FALSE, TRUE),
-- 販売費及び一般管理費
(7200, '給与・賞与', FALSE, TRUE),
(7300, '経費', FALSE, TRUE),
-- 営業外損益・特別損益
(8100, '営業外収益', FALSE, FALSE),
(8200, '営業外費用', FALSE, TRUE),
(8300, '特別利益', FALSE, FALSE),
(8400, '特別損失', FALSE, TRUE),
-- 法人税等
(9000, '法人税等', FALSE, TRUE)
;
