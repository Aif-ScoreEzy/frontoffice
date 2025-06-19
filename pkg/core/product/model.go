package product

type productResponseData struct {
	ProductId   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductSlug string `json:"product_slug_name"`
}
