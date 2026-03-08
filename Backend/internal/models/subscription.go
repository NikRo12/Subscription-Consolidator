package models

type Period string
type Category string

const (
	Monthly Period = "monthly"
	Yearly  Period = "yearly"
)

const (
	Entertainment Category = "entertainment"
	Work          Category = "work"
	Clouds        Category = "clouds"
	Food          Category = "food"
	Health        Category = "health"
	Other         Category = "other"
)

type Subscription struct {
	ID          int
	Title       string
	Currency    string
	Category    Category
	IconURL     string
	BrandColor  string
	Description string
}
