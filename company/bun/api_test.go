package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/tullo/crdb/company/gobun/model"
)

func TestCompanyAPI(t *testing.T) {
	// Customer
	c := createCustomer(t, strings.NewReader(`{"name": "John Doe"}`))
	t.Logf("Customer %v", *c)

	// Product
	var p *model.Product
	ps := getProducts(t)
	for i := range ps {
		if ps[i].Name == "GopherCon Europe 2021" {
			p = &ps[i]
		}
	}
	if p == nil {
		p = createProduct(t, strings.NewReader(`{"name": "GopherCon Europe 2021", "Price": 315.04}`))
		t.Logf("Product %+v", *p)
	}

	// Order
	o := createOrder(t, strings.NewReader(fmt.Sprintf(
		`{"subtotal": 630.08,"customer": {"id": %d},"products": [{"id": %d}]}`, c.ID, p.ID)))

	t.Logf("Order %+v", *o)

	//p = createProduct(t, strings.NewReader(`{"name": "GopherCon Europe 2022", "Price": 317.07}`))
	//addProductToOrder(t, strconv.Itoa(o.ID), strconv.Itoa(p.ID))
}

func createCustomer(t *testing.T, payload io.Reader) *model.Customer {
	res, err := http.Post("http://api.0.0.0.0.nip.io:6543/customer",
		"application/json; charset=UTF-8", payload)
	if err != nil {
		t.Fatal("error creating customer", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal("error reading body", err)
		}
		t.Fatal(res.StatusCode, string(body))
	}
	var c model.Customer
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatal("error decoding customer", err)
	}

	return &c
}

func getProducts(t *testing.T) []model.Product {
	res, err := http.Get("http://api.0.0.0.0.nip.io:6543/product")
	if err != nil {
		t.Error("error retrieving products", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal("error reading body", err)
		}
		t.Fatal(res.StatusCode, string(body))
	}
	var ps []model.Product
	if err := json.NewDecoder(res.Body).Decode(&ps); err != nil {
		t.Fatal("error decoding products", err)
	}

	return ps
}

func createProduct(t *testing.T, payload io.Reader) *model.Product {
	res, err := http.Post("http://api.0.0.0.0.nip.io:6543/product",
		"application/json; charset=UTF-8", payload)
	if err != nil {
		t.Error("error creating product", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal("error reading body", err)
		}
		t.Fatal(res.StatusCode, string(body))
	}
	var p model.Product
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatal("error decoding product", err)
	}

	return &p
}

func createOrder(t *testing.T, payload io.Reader) *model.Order {
	res, err := http.Post("http://api.0.0.0.0.nip.io:6543/order",
		"application/json; charset=UTF-8", payload)
	if err != nil {
		t.Error("error creating order", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal("error reading body", err)
		}
		t.Fatal(res.StatusCode, string(body))
	}

	var o model.Order
	if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
		t.Fatal("error decoding order", err)
	}

	return &o
}

func addProductToOrder(t *testing.T, oid, pid string) *model.Order {
	res, err := http.Post("http://api.0.0.0.0.nip.io:6543/order/"+oid+"/product?productID="+pid,
		"application/json; charset=UTF-8", nil)
	if err != nil {
		t.Error("error adding product to order", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal("error reading body", err)
		}
		t.Fatal(res.StatusCode, string(body))
	}

	var o model.Order
	if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
		t.Fatal("error decoding order", err)
	}
	t.Log("addProductToOrder", o)

	return &o
}
