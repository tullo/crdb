CREATE TABLE order_products
(
   order_id integer REFERENCES orders ON DELETE CASCADE,
   product_id integer REFERENCES products ON DELETE RESTRICT,
   PRIMARY KEY (product_id, order_id)
);
