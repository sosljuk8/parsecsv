package dto

// PCard is a data transfer object for pcard data.
type PCard struct {
	Brand       string
	Category    string
	Model       string
	Name        string
	SKU         string
	Price       string
	Currency    string
	Source      string
	Img         string
	Properties  string
	Description string
	File        string
}

// NewPCard creates a new PCard.
func NewPCard() *PCard {
	return &PCard{}
}

// String returns a string representation of the PCard.
func (p *PCard) String() []string {

	str := []string{p.Brand, p.Category, p.Model, p.Name, p.SKU, p.Price, p.Currency, p.Source, p.Img, p.Properties, p.Description, p.File}

	// brand, category, model, name, sku, price, currency, source
	return str
}
