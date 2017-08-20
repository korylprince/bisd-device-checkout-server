CREATE TABLE devices (
  id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  Inventory_Number varchar(255) DEFAULT NULL,
  Serial_Number varchar(255) DEFAULT NULL,
  Bag_Tag varchar(255) DEFAULT NULL,
  status varchar(255) DEFAULT NULL,
  User varchar(255) DEFAULT NULL,
  user_type varchar(255) DEFAULT NULL,
  device_type varchar(255) DEFAULT NULL,
  manufacturer varchar(255) DEFAULT NULL,
  model varchar(255) DEFAULT NULL,
  Campus varchar(255) DEFAULT NULL,
  Room varchar(255) DEFAULT NULL,
  on_network char(1) DEFAULT NULL,
  notes longtext,
  po_number varchar(255) DEFAULT NULL,
  UNIQUE KEY Inventory Number (Inventory_Number)
)

CREATE TABLE charges (
  id INTEGER UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  Inventory_Number varchar(255) DEFAULT NULL,
  User varchar(255) DEFAULT NULL,
  Amount_Paid double DEFAULT NULL,
  Charges longtext,
  Notes longtext
)
