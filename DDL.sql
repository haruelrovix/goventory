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
  foreign key(warehouse_id) references warehouses(warehouse_id)
);

CREATE TABLE TransactionTypes(
  code varchar(2) primary key,
  description varchar(50)
);

CREATE TABLE transactions(
  id int primary key,
  amount int,
  note varchar(50),
  price int,
  timestamp DATETIME,
  transaction_code varchar(2),
  transaction_sku varchar(20),
  foreign key(transaction_code) references TransactionTypes(code),
  foreign key(transaction_sku) references items(sku)
);

CREATE TABLE OrderDetails(
  id varchar(18) primary key
);

CREATE TABLE IncomingTransactions(
  transaction_id int,
  order_id varchar(18),
  receipt varchar(14),
  foreign key(transaction_id) references transactions(id),
  foreign key(order_id) references OrderDetails(id)
);
