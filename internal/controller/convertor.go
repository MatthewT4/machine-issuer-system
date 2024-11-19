package controller

import (
	"machineIssuerSystem/api"
	"machineIssuerSystem/internal/model"
)

func convertProductToAPI(product model.Product) api.Product {
	return api.Product{
		Id:          &product.ID,
		Title:       &product.Title,
		Tags:        &product.Tags,
		Description: &product.Description,
		ImageUrls:   &product.ImageURLs,
	}
}
