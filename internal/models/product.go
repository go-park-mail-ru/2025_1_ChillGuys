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

type BriefProduct struct{
	ID int
	Name string
	Image []byte
	Price uint
	ReviewsCount uint
	Rating float64
}

func ConvertToBriefProduct(product *Product) BriefProduct{
	return BriefProduct{
		ID:           product.ID,
		Name:         product.Name,
		Image:        product.Image,
		Price:        product.Price,
		ReviewsCount: product.ReviewsCount,
		Rating:       product.Rating,
	}
}