package models

type Product struct {
	Id          int
	Price       float64
	PhotoUrl    string
	Composition []string
	Name        string
}
