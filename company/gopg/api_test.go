package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/tullo/crdb/company/gopg/model"
)

func TestCompanyAPI(t *testing.T) {
	// Customer
	c := createCustomer(t, strings.NewReader(`{"name": "John Doe"}`))

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
	}

	// Order
	o := createOrder(t, strings.NewReader(fmt.Sprintf(
		`{"subtotal": 630.08,"customer": {"id": %d},"products": [{"id": %d}]}`, c.ID, p.ID)))

	t.Logf("Order %+v", o)
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
	t.Log("createCustomer", c.ID)

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
	t.Log("getProducts", len(ps))

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
	t.Log("createProduct", p.ID)

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
	t.Log("createOrder", o.ID)

	return &o
}
