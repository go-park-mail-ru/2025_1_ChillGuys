package models

type Product struct{
	ID int
	Name string 
	Description string
	Count uint
	Image []byte
	Price uint
	ReviewsCount uint
	Rating float64
}

