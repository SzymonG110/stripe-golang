package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	stripe "github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
	"log"
	"net/http"
	"os"
)

func main() {
	// fmt.Printf("%+v\n", JSON)
	// ^^ print json XD
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
	r := gin.Default()

	r.GET("/products", func(c *gin.Context) {
		products := GetProducts()
		c.JSON(http.StatusOK, products)
	})

	r.POST("/create-checkout-session", func(c *gin.Context) {
		type RequestBody struct {
			ProductID string `json:"product_id"`
			Nickname  string `json:"nickname"`
		}

		var reqBody RequestBody
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		productData, err := GetProduct(reqBody.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
			return
		}

		session, err := CreateCheckoutSession(reqBody.ProductID, productData.PriceID, reqBody.Nickname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": session.ID, "url": session.URL})
	})

	r.POST("/webhook", func(c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		payload, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		event := stripe.Event{}
		event, err = webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to verify webhook signature"})
			return
		}

		if err := json.Unmarshal(payload, &event); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing event data"})
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			var checkoutSession stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &checkoutSession)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error parsing event data"})
				return
			}

			nickname := checkoutSession.Metadata["nickname"]
			productId := checkoutSession.Metadata["product_id"]
			if nickname == "" || productId == "" {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "nickname or product_id not provided"})
				return
			}

			result, err := GetProduct(productId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Incorrect product_id"})
				return
			}

			SendWebhook(nickname + " kupił " + result.Name)

		}

		c.JSON(http.StatusOK, gin.H{"received": true})
	})

	r.Run(":" + os.Getenv("API_PORT"))
}
