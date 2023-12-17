package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cosmovid21/Billing-system/controllers"
	controllers "github.com/cosmovid21/Billing-system/database"
	"github.com/cosmovid21/Billing-system/middleware"
	"github.com/cosmovid21/Billing-system/routes"
	"github.com/gin-gonic/gin"
)

type Item struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

func applyTax(item *Item) {
	if item.Price > 1000 && item.Price <= 5000 {
		item.Tax = item.Price * 0.12
	} else if item.Price > 5000 {
		item.Tax = item.Price * 0.18
	} else {
		item.Tax = 0
	}
}

func applyServiceTax(item *Item) {
	if item.Price > 1000 && item.Price <= 8000 {
		item.Tax = item.Price * 0.10
	} else if item.Price > 8000 {
		item.Tax = item.Price * 0.15
	} else {
		item.Tax = 0
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductorServiceData(database.Client, "Products and services"), database.Userdata(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	r.POST("/calculateTax", func(c *gin.Context) {
		var item Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if item.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product/Service name cannot be empty"})
			return
		}

		if item.Price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
			return
		}

		if item.Price > 0 {
			if item.Name == "Product" {
				applyTax(&item)
			} else if item.Name == "Service" {
				applyServiceTax(&item)
			} else {
				item.Tax = 0
			}
		}

		c.JSON(http.StatusOK, gin.H{"item": item})
	})

	r.Run(":8000")

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyfromCart())

	log.Fatal(router.Run(":" + port))
}
