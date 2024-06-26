package main

import (
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/price"
	"github.com/stripe/stripe-go/v78/product"
	"os"
)

type Product struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	PriceID     string          `json:"price_id"`
	Price       float64         `json:"price"`
	Currency    stripe.Currency `json:"currency"`
}

type Price struct {
	ID       string          `json:"id"`
	Price    float64         `json:"price"`
	Currency stripe.Currency `json:"currency"`
}

type Checkout struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func GetProducts() []Product {
	productParams := &stripe.ProductListParams{}
	productParams.Filters.AddFilter("active", "", "true")

	iter := product.List(productParams)
	products := []Product{}
	for iter.Next() {
		productData := iter.Product()
		if productData.DefaultPrice == nil {
			continue
		}

		priceData, err := GetProductPrice(productData.DefaultPrice.ID)
		if err != nil {
			continue
		}

		products = append(products, Product{
			ID:          productData.ID,
			Name:        productData.Name,
			Description: productData.Description,
			PriceID:     priceData.ID,
			Price:       priceData.Price,
			Currency:    priceData.Currency,
		})
	}

	return products
}

func GetProduct(productID string) (Product, error) {
	productParams := &stripe.ProductParams{}
	productData, err := product.Get(productID, productParams)
	if err != nil {
		return Product{}, err
	}

	priceData, err := GetProductPrice(productData.DefaultPrice.ID)
	if err != nil {
		return Product{}, err
	}

	return Product{
		ID:          productData.ID,
		Name:        productData.Name,
		Description: productData.Description,
		PriceID:     priceData.ID,
		Price:       priceData.Price,
		Currency:    priceData.Currency,
	}, nil
}

func GetProductPrice(priceID string) (Price, error) {
	priceParams := &stripe.PriceParams{}
	priceData, err := price.Get(priceID, priceParams)
	if err != nil {
		return Price{}, err
	}

	return Price{
		ID:       priceData.ID,
		Price:    float64(priceData.UnitAmount) / 100,
		Currency: priceData.Currency,
	}, nil
}

func CreateCheckoutSession(productID string, priceID string, nickname string) (Checkout, error) {
	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"blik",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(os.Getenv("CHECKOUT_SUCCESS_URL") + "?nickname=" + nickname),
		CancelURL:  stripe.String(os.Getenv("CHECKOUT_CANCEL_URL")),
		Metadata: map[string]string{
			"nickname":   nickname,
			"product_id": productID,
		},
	}

	checkoutSession, err := session.New(sessionParams)
	if err != nil {
		return Checkout{}, err
	}

	return Checkout{
		ID:  checkoutSession.ID,
		URL: checkoutSession.URL,
	}, nil
}
