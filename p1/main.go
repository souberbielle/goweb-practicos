package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	CodeValue   string  `json:"code_value" validate:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}

func main() {
	router := gin.Default()
	jsonFile, err := os.Open("products.json")

	if err != nil {
		return
	}
	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return
	}

	var products []Product
	if err = json.Unmarshal(byteValue, &products); err != nil {
		return
	}

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, &products)
	})

	router.GET("/products/:id", func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error de formato en el ID.",
			})
		}
		var productFound bool = false
		for _, product := range products {
			if product.Id == id {
				productFound = true
				ctx.JSON(200, product)
				break
			}
		}
		if !productFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found."})
		}
	})
	router.GET("/products/search", func(ctx *gin.Context) {
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Error amigo!",
			})
		}
		var productFound bool = false
		for _, product := range products {
			if product.Price > price {
				productFound = true
				ctx.JSON(200, product)
			}
		}
		if !productFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Product not found.",
			})
		}
	})

	router.Run(":8082")
	defer jsonFile.Close()

}
