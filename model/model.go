package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// getProducts returns all products from the database
func GetProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT id,name,quantity,price FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// getProduct returns a single product from the database
func GetProduct(db *sql.DB, id int) (Product, error) {
	query := "SELECT id,name,quantity,price FROM products WHERE id=?"
	row := db.QueryRow(query, id)
	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return p, err
	}
	return p, nil
}

// p * Product is a pointer to a Product struct which is a receiver
func (p *Product) GetProductt(db *sql.DB) error {
	query := fmt.Sprintf("SELECT id,name,quantity,price FROM products WHERE id=%d", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

// createProduct creates a new product in the database
func (p *Product) CreateProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products (name,quantity,price) VALUES('%s',%d,%f)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (p *Product) UpdateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products SET name='%s',quantity=%d,price=%f WHERE id=%d", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	log.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No such row exists")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) DeleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id=%d", p.ID)
	result, err := db.Exec(query)
	log.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No such row exists")
	}
	if err != nil {
		return err
	}
	return nil
}
