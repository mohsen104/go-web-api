package dto

type CreateUpdateCountryRequest struct {
	Name string `json:"name" binding:"required,alpha,min=3,max=20"`
}

type CountryResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
