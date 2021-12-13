package model

type Customer struct {
	ID   int    `json:"id,omitempty" bun:",pk"`
	Name string `json:"name" bun:",notnull"`
}

type Order struct {
	ID       int     `json:"id,omitempty" bun:",pk"`
	Subtotal float64 `json:"subtotal" bun:"type:decimal(18,2)"`

	// Order belongs to Customer.
	Customer   Customer `json:"customer" bun:"rel:belongs-to,join:customer_id=id"`
	CustomerID int      `json:"-" bun:",notnull"`

	Products []Product `json:"products" bun:"m2m:order_to_products,join:Order=Product"`
}

type Product struct {
	ID    int     `json:"id,omitempty" bun:",pk"`
	Name  string  `json:"name" bun:",notnull,unique"`
	Price float64 `json:"price" bun:"type:decimal(18,2)"`
}

type OrderToProduct struct {
	OrderID int    `bun:",pk"`
	Order   *Order `bun:"rel:belongs-to,join:order_id=id"`

	ProductID int      `json:"-" bun:",pk"`
	Product   *Product `json:"product" bun:"rel:belongs-to,join:product_id=id"`
}
