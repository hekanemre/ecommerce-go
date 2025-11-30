package infrastructure

import (
	"database/sql"
	"ecommerce-go/domain"
	"errors"
)

// ProductRepository defines CRUD operations for Product
type ProductRepository interface {
	Create(product *domain.Product) (int64, error)
	GetByID(id int64) (*domain.Product, error)
	GetAll() ([]*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id int64) error
}

// productRepo is the concrete implementation
type productRepo struct {
	db Repository
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db Repository) ProductRepository {
	return &productRepo{db: db}
}

// Create inserts a new product into the database
func (r *productRepo) Create(product *domain.Product) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO Product (Name, Price, Description) VALUES (?, ?, ?)",
		product.Name, product.Price, product.Description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetByID retrieves a product by its ID
func (r *productRepo) GetByID(id int64) (*domain.Product, error) {
	row := r.db.QueryRow("SELECT ID, Name, Price, Description FROM Product WHERE ID = ?", id)
	p := &domain.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

// GetAll retrieves all products
func (r *productRepo) GetAll() ([]*domain.Product, error) {
	rows, err := r.db.Query("SELECT ID, Name, Price, Description FROM Product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// Update modifies an existing product
func (r *productRepo) Update(product *domain.Product) error {
	_, err := r.db.Exec(
		"UPDATE Product SET Name = ?, Price = ?, Description = ? WHERE ID = ?",
		product.Name, product.Price, product.Description, product.ID,
	)
	return err
}

// Delete removes a product by its ID
func (r *productRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM Product WHERE ID = ?", id)
	return err
}
