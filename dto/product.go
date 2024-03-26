package dto

// Product is a data transfer object for product data.
type Product struct {
	Brand    string
	Category string
	Model    string
	Name     string
	SKU      string
	Price    string
	Currency string
	Source   string
	File     string
}

// NewProduct creates a new Product.
func NewProduct(brand, category, model, name, sku, price, currency, source, file string) *Product {
	return &Product{
		Brand:    brand,
		Category: category,
		Model:    model,
		Name:     name,
		SKU:      sku,
		Price:    price,
		Currency: currency,
		Source:   source,
		File:     file,
	}
}

// String returns a string representation of the Product.
func (p *Product) String() []string {

	str := []string{p.Brand, p.Category, p.Model, p.Name, p.SKU, p.Price, p.Currency, p.Source, p.File}

	// brand, category, model, name, sku, price, currency, source
	return str
}
