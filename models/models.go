package models

type Category struct {
	Id       int
	Name     string
	Link     string
	ParentId int
	Level    int
}

type Product struct {
	Id                string
	CategoryId        int
	Link              string
	Title             string
	Name              string
	PriceOrg          string
	PriceDisc         string
	PriceDiscDesc     string
	PriceDiscStamp    string
	ColorVariantCount string
	Sizes             []string
	Images            []Image
	Colors            []string
}

type Image struct {
	BaseImage bool
	IsStamp   bool
	Link      string
}

type Brand struct {
	Name string
	Logo string
}
