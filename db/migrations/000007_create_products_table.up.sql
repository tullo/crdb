CREATE TABLE IF NOT EXISTS products
(
   id INT PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL,
   price NUMERIC DEFAULT 99.9 NOT NULL
);