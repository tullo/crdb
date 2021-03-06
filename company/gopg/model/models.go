package model

// Customer is a model in the "customers" table.
type Customer struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name" pg:",notnull"`
}

// Order is a model in the "orders" table.
type Order struct {
	ID       int     `json:"id,omitempty"`
	Subtotal float64 `json:"subtotal" pg:"type:'decimal(18,2)'"`

	Customer   Customer `json:"customer" pg:"fk:customer_id,rel:has-one"`
	CustomerID int      `json:"-"`

	Products []Product `json:"products" pg:"many2many:order_products"`
}

// Product is a model in the "products" table.
type Product struct {
	ID    int     `json:"id,omitempty"`
	Name  string  `json:"name" pg:",notnull,unique"`
	Price float64 `json:"price" pg:"type:'decimal(18,2)'"`
}

// OrderProduct is a model in the "order_products" table.
type OrderProduct struct {
	Order   Order `json:"order" pg:"fk:order_id,on_delete:CASCADE,rel:has-one"`
	OrderID int   `json:"-"`

	Product   Product `json:"product" pg:"fk:product_id,rel:has-one"`
	ProductID int     `json:"-"`
}
