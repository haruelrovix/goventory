CREATE TABLE items(
  sku varchar(20) primary key,
  name varchar(50)
);

CREATE TABLE warehouses(
  id int primary key,
  description varchar(50)
);

CREATE TABLE stock(
  item_sku varchar(20),
  warehouse_id int,
  amount int,
  foreign key(item_sku) references items(sku),
  foreign key(warehouse_id) references items(warehouse_id)
);
