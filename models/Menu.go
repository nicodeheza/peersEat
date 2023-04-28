package models

type DishOptions struct {
	Name        string
	Description string
	Price       float32
}

type Dish struct {
	Name        string
	Description string
	Price       float32
	ImageUrl    string
	Options     []DishOptions `bson:"options,omitempty" json:"options,omitempty"`
}

type MenuSection struct {
	Name   string
	Dishes []Dish
}

type Menu struct {
	Sections []MenuSection
}
