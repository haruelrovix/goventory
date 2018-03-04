INSERT INTO items (sku, name)
VALUES ( 'SSI-D00791015-LL-BWH', 'Zalekia Plain Casual Blouse (L,Broken White)' ),
( 'SSI-D00791077-MM-BWH', 'Zalekia Plain Casual Blouse (M,Broken White)' ),
( 'SSI-D00791091-XL-BWH', 'Zalekia Plain Casual Blouse (XL,Broken White)' ),
( 'SSI-D00864612-LL-NAV', 'Deklia Plain Casual Blouse (L,Navy)' ),
( 'SSI-D00864614-XL-NAV', 'Deklia Plain Casual Blouse (XL,Navy)' ),
( 'SSI-D00864652-SS-NAV', 'Deklia Plain Casual Blouse (S,Navy)' ),
( 'SSI-D00864661-MM-NAV', 'Deklia Plain Casual Blouse (M,Navy)' ),
( 'SSI-D01037807-X3-BWH', 'Dellaya Plain Loose Big Blouse (XXXL,Broken White)' ),
( 'SSI-D01037812-X3-BLA', 'Dellaya Plain Loose Big Blouse (XXXL,Black)' ),
( 'SSI-D01037822-XX-BLA', 'Dellaya Plain Loose Big Blouse (XXL,Black)' ),
( 'SSI-D01220307-XL-SAL', 'Devibav Plain Trump Blouse (XL,Salem)' ),
( 'SSI-D01220322-MM-YEL', 'Devibav Plain Trump Blouse (M,Yellow)' ),
( 'SSI-D01220334-XL-YEL', 'Devibav Plain Trump Blouse (XL,Yellow)' ),
( 'SSI-D01220338-XX-SAL', 'Devibav Plain Trump Blouse (XXL,Salem)' ),
( 'SSI-D01220346-LL-SAL', 'Devibav Plain Trump Blouse (L,Salem)' ),
( 'SSI-D01220349-LL-YEL', 'Devibav Plain Trump Blouse (L,Yellow)' ),
( 'SSI-D01220355-XX-YEL', 'Devibav Plain Trump Blouse (XXL,Yellow)' ),
( 'SSI-D01220357-SS-YEL', 'Devibav Plain Trump Blouse (S,Yellow)' ),
( 'SSI-D01220388-MM-SAL', 'Devibav Plain Trump Blouse (M,Salem)' ),
( 'SSI-D01322234-LL-WHI', 'Thafqya Plain Raglan Blouse (L,White)' ),
( 'SSI-D01322275-XL-WHI', 'Thafqya Plain Raglan Blouse (XL,White)' ),
( 'SSI-D01326201-XL-KHA', 'Siunfhi Ethnic Trump Blouse (XL,Khaki)' ),
( 'SSI-D01326205-MM-NAV', 'Siunfhi Ethnic Trump Blouse (M,Navy)' ),
( 'SSI-D01326223-MM-KHA', 'Siunfhi Ethnic Trump Blouse (M,Khaki)' ),
( 'SSI-D01326286-LL-KHA', 'Siunfhi Ethnic Trump Blouse (L,Khaki)' ),
( 'SSI-D01326299-LL-NAV', 'Siunfhi Ethnic Trump Blouse (L,Navy)' ),
( 'SSI-D01401050-MM-RED', 'Zeomila Zipper Casual Blouse (M,Red)' ),
( 'SSI-D01401064-XL-RED', 'Zeomila Zipper Casual Blouse (XL,Red)' ),
( 'SSI-D01401071-LL-RED', 'Zeomila Zipper Casual Blouse (L,Red)' ),
( 'SSI-D01466013-XX-BLA', 'Salyara Plain Casual Big Blouse (XXL,Black)' ),
( 'SSI-D01466064-X3-BLA', 'Salyara Plain Casual Big Blouse (XXXL,Black)' );

INSERT INTO warehouses (id, description)
VALUES ( 0, 'Warehouse at Kaliurang' ),
( 1, 'Warehouse at Paskal Hyper Square' );

INSERT INTO stock (item_sku, warehouse_id, amount)
VALUES ( 'SSI-D00791077-MM-BWH', 0, 100 ),
( 'SSI-D00791077-MM-BWH', 1, 54 ),
( 'SSI-D00864612-LL-NAV', 0, 47 ),
( 'SSI-D00864612-LL-NAV', 1, 90 ),
( 'SSI-D01220357-SS-YEL', 0, 74 );

INSERT INTO TransactionTypes (code, description)
VALUES ( 'BM', 'Barang Masuk' ),
( 'BK', 'Barang Keluar' );

INSERT INTO transactions (id, timestamp, sku )
VALUES ( 'BM', 'Barang Masuk' ),
( 'BK', 'Barang Keluar' );
