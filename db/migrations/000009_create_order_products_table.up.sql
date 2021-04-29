CREATE TABLE order_products
(
   order_id integer REFERENCES orders,
   product_id integer REFERENCES products,
   PRIMARY KEY (product_id, order_id)
);