package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/tullo/crdb/company/gobun/model"
	"github.com/uptrace/bun"
)

// Server is a http server that handles REST requests.
type Server struct {
	db *bun.DB
}

// NewServer creates a new instance of Server.
func NewServer(db *bun.DB) *Server {
	return &Server{db: db}
}

// RegisterRoutes maps routes to handlers.
func (s *Server) RegisterRouter(router *httprouter.Router) {
	router.GET("/ping", s.ping)

	router.GET("/customer", s.getCustomers)
	router.POST("/customer", s.createCustomer)
	router.GET("/customer/:customerID", s.getCustomer)
	router.PUT("/customer/:customerID", s.updateCustomer)
	router.DELETE("/customer/:customerID", s.deleteCustomer)

	router.GET("/product", s.getProducts)
	router.POST("/product", s.createProduct)
	router.GET("/product/:productID", s.getProduct)
	router.PUT("/product/:productID", s.updateProduct)
	router.DELETE("/product/:productID", s.deleteProduct)

	router.GET("/order", s.getOrders)
	router.POST("/order", s.createOrder)
	router.GET("/order/:orderID", s.getOrder)
	router.PUT("/order/:orderID", s.updateOrder)
	router.DELETE("/order/:orderID", s.deleteOrder)
	router.POST("/order/:orderID/product", s.addProductToOrder)
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	writeTextResult(w, "go/gopg")
}

func (s *Server) getCustomers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var customers []model.Customer
	if err := s.db.NewSelect().Model(&customers).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, customers)
	}
}

func (s *Server) createCustomer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	if _, err := s.db.NewInsert().Model(&customer).Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, customer)
	}
}

func (s *Server) getCustomer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	customerID, err := strconv.Atoi(ps.ByName("customerID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	customer := model.Customer{
		ID: customerID,
	}
	if err := s.db.NewSelect().Model(&customer).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, customer)
	}
}

func (s *Server) updateCustomer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	customerID, err := strconv.Atoi(ps.ByName("customerID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	customer.ID = customerID
	if _, err := s.db.NewUpdate().Model(&customer).WherePK().Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, customer)
	}
}

func (s *Server) deleteCustomer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	customerID, err := strconv.Atoi(ps.ByName("customerID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	customer := model.Customer{
		ID: customerID,
	}
	res, err := s.db.NewDelete().Model(&customer).WherePK().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "", http.StatusNotFound)
	} else {
		writeTextResult(w, "ok")
	}
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var products []model.Product
	if err := s.db.NewSelect().Model(&products).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, products)
	}
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	if _, err := s.db.NewInsert().Model(&product).Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, product)
	}
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	productID, err := strconv.Atoi(ps.ByName("productID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	product := model.Product{
		ID: productID,
	}
	if err := s.db.NewSelect().Model(&product).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, product)
	}
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	productID, err := strconv.Atoi(ps.ByName("productID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	product.ID = productID
	if _, err := s.db.NewUpdate().Model(&product).WherePK().Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, product)
	}
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	productID, err := strconv.Atoi(ps.ByName("productID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	product := model.Product{
		ID: productID,
	}
	res, err := s.db.NewDelete().Model(&product).WherePK().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "", http.StatusNotFound)
	} else {
		writeTextResult(w, "ok")
	}
}

func (s *Server) getOrders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var orders []model.Order
	if err := s.db.NewSelect().Model(&orders).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, orders)
	}
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	if order.Customer.ID == 0 {
		http.Error(w, "must specify user", http.StatusBadRequest)
		return
	}
	if err := s.db.NewSelect().Model(&order.Customer).Where("id = ?", order.Customer.ID).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}
	order.CustomerID = order.Customer.ID

	for i, product := range order.Products {
		if product.ID == 0 {
			http.Error(w, "must specify a product ID", http.StatusBadRequest)
			return
		}
		if err := s.db.NewSelect().Model(&order.Products[i]).Where("id = ?", product.ID).Scan(r.Context()); err != nil {
			http.Error(w, err.Error(), errToStatusCode(err))
			return
		}
	}

	if _, err := s.db.NewInsert().Model(&order).Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		for _, p := range order.Products {
			var rel model.OrderToProduct
			rel.Order = &order
			rel.OrderID = order.ID
			rel.Product = &p
			rel.ProductID = p.ID
			if _, err := s.db.NewInsert().Model(&rel).Exec(r.Context()); err != nil {
				http.Error(w, err.Error(), errToStatusCode(err))
			}
		}
		writeJSONResult(w, order)
	}
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID, err := strconv.Atoi(ps.ByName("orderID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	order := model.Order{
		ID: orderID,
	}
	if err := s.db.NewSelect().Model(&order).Scan(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, order)
	}
}

func (s *Server) updateOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var order model.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	orderID, err := strconv.Atoi(ps.ByName("orderID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	order.ID = orderID
	if _, err := s.db.NewUpdate().Model(&order).WherePK().Exec(r.Context()); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		writeJSONResult(w, order)
	}
}

func (s *Server) deleteOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	orderID, err := strconv.Atoi(ps.ByName("orderID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	order := model.Order{
		ID: orderID,
	}
	res, err := s.db.NewDelete().Model(&order).WherePK().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "", http.StatusNotFound)
	} else {
		writeTextResult(w, "ok")
	}
}

func (s *Server) addProductToOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}

	orderID, err := strconv.Atoi(ps.ByName("orderID"))
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	order := model.Order{
		ID: orderID,
	}
	if err := s.db.NewSelect().Model(&order).Scan(r.Context()); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	const productIDParam = "productID"
	productIDString := r.URL.Query().Get(productIDParam)
	if productIDString == "" {
		tx.Rollback()
		writeMissingParamError(w, productIDParam)
		return
	}

	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	}
	addedProduct := model.Product{
		ID: productID,
	}
	if err := s.db.NewSelect().Model(&addedProduct).WherePK().Scan(r.Context()); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	orderProduct := model.OrderToProduct{
		Order:     &order,
		OrderID:   order.ID,
		Product:   &addedProduct,
		ProductID: addedProduct.ID,
	}
	if _, err := tx.NewInsert().Model(&orderProduct).Exec(r.Context()); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), errToStatusCode(err))
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), errToStatusCode(err))
	} else {
		err = s.db.NewSelect().Model(&order).Where("customer_id = ?", orderProduct.OrderID).Scan(r.Context())
		if err != nil {
			http.Error(w, err.Error(), errToStatusCode(err))
		}
		writeJSONResult(w, order)
	}
}

func writeTextResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, res)
}

func writeJSONResult(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func writeMissingParamError(w http.ResponseWriter, paramName string) {
	http.Error(w, fmt.Sprintf("missing query param %q", paramName), http.StatusBadRequest)
}

func errToStatusCode(err error) int {
	switch err {
	case pg.ErrNoRows:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
