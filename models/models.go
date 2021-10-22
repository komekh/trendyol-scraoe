package models

type Category struct {
	Id       int
	Name     string
	Link     string
	ParentId int
	Level    int
}

type Product struct {
	Id             int
	CategoryId     int
	Link           string
	Title          string
	Name           string
	PriceOrg       string
	PriceDisc      string
	PriceDiscDesc  string
	PriceDiscStamp string
	Sizes          []string
}
