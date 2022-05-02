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

type ProductResponse struct {
	Name       string  `json:"name"`
	ProductId  string  `json:"product_id"`
	SalesPrice float32 `json:"sales_price"`
	Status     string  `json:"status"`
}

func ValidateProduct(c *gin.Context) {
	var json Product
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	var url = "http://wiremock:8080/v1/products/" + json.ProductId
	resp, err := http.Get(url)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"err": err})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"err": err, "message": "NOT FOUND"})
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)

	jsonResponse := new(ProductResponse)

	if strings.Contains(string(data), "In Stock") {
		jsonResponse = populateResponse(json, "valid")
	} else {
		jsonResponse = populateResponse(json, "invalid")
	}

	c.JSON(200, gin.H{"body": *jsonResponse})
}

func populateResponse(
	json Product, status string) *ProductResponse {

	response := new(ProductResponse)

	response.Name = json.Name
	response.ProductId = json.ProductId
	response.SalesPrice = json.SalesPrice
	response.Status = status

	return response
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
