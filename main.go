package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type Product struct {
	Name       string  `json:"name"`
	ProductId  string  `json:"product_id"`
	SalesPrice float32 `json:"sales_price"`
}

func ValidateProduct(c *gin.Context) {
	var json Product
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var url = "localhost:8080/products/" + json.ProductId
	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse struct {
		Name       string  `json:"name"`
		ProductId  string  `json:"product_id"`
		SalesPrice float32 `json:"sales_price"`
		Status     string  `json:"status"`
	}

	if strings.Contains(string(data), "status: \"valid\"") {
		populateResponse(jsonResponse, json, "valid")
	} else {
		populateResponse(jsonResponse, json, "invalid")
	}
}

func populateResponse(
	jsonResponse struct {
		Name       string  `json:"name"`
		ProductId  string  `json:"product_id"`
		SalesPrice float32 `json:"sales_price"`
		Status     string  `json:"status"`
	},
	json Product, status string) {

	jsonResponse.Name = json.Name
	jsonResponse.ProductId = json.ProductId
	jsonResponse.SalesPrice = json.SalesPrice
	jsonResponse.Status = status
}

func main() {
	fmt.Println("Hello World")
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		v1.POST("products/validate", ValidateProduct)
	}

	r.Run(":8000")
}
