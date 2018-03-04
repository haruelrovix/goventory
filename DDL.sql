CREATE TABLE items(
  sku VARCHAR(20) PRIMARY KEY,
  name VARCHAR(50)
);

CREATE TABLE warehouses(
  id INT PRIMARY KEY,
  description VARCHAR(50)
);

CREATE TABLE stock(
  item_sku VARCHAR(20),
  warehouse_id INT,
  amount INT,
  FOREIGN KEY(item_sku) REFERENCES items(sku),
  FOREIGN KEY(warehouse_id) REFERENCES warehouses(warehouse_id)
);

CREATE TABLE TransactionTypes(
  code VARCHAR(2) PRIMARY KEY,
  description VARCHAR(50)
);

CREATE TABLE transactions(
  id INT PRIMARY KEY,
  amount INT,
  note VARCHAR(50),
  price NUMERIC,
  timestamp DATETIME,
  transaction_code VARCHAR(2),
  transaction_sku VARCHAR(20),
  FOREIGN KEY(transaction_code) REFERENCES TransactionTypes(code),
  FOREIGN KEY(transaction_sku) REFERENCES items(sku)
);

CREATE TABLE OrderDetails(
  id VARCHAR(18) PRIMARY KEY
);

CREATE TABLE IncomingTransactions(
  transaction_id INT,
  booking INT,
  receipt VARCHAR(14),
  FOREIGN KEY(transaction_id) REFERENCES transactions(id)
);

CREATE TABLE OutgoingTransactions(
  transaction_id INT,
  order_id VARCHAR(18),
  FOREIGN KEY(transaction_id) REFERENCES transactions(id),
  FOREIGN KEY(order_id) REFERENCES OrderDetails(order_id)
);
