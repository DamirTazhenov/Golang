package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func productToJSON(p Product) (string, error) {
	jsonProduct, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(jsonProduct), nil
}

func JSONToProduct(data string) (Product, error) {
	var p Product
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

func testJson() {
	prod := Product{Name: "Laptop", Price: 1234.56, Quantity: 100}

	jsonString, err := productToJSON(prod)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println("Encoded JSON:", jsonString)

	decodedProduct, err := JSONToProduct(jsonString)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println("Decoded Product:", decodedProduct)
}
