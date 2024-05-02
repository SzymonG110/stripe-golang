package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DBPayment struct {
	Email       string    `bson:"email,omitempty"`
	Nickname    string    `bson:"nickname,omitempty"`
	ProductName string    `bson:"productName,omitempty"`
	Date        time.Time `bson:"date,omitempty"`
}

func InsertPayment(collection *mongo.Collection, item DBPayment) error {
	_, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}

	return nil
}

func GetLast(collection *mongo.Collection, limit uint) ([]DBPayment, error) {
	payments := []DBPayment{}
	productsOptions := options.Find().SetSort(bson.D{{"date", -1}}).SetLimit(int64(limit))
	productsCursor, err := collection.Find(context.TODO(), bson.D{}, productsOptions)
	if err != nil {
		return payments, err
	}

	for productsCursor.Next(context.TODO()) {
		var payment DBPayment
		if err := productsCursor.Decode(&payment); err != nil {
			return payments, err
		}
		payments = append(payments, DBPayment{
			Email:       "not to leak XD",
			Nickname:    payment.Nickname,
			ProductName: payment.ProductName,
			Date:        payment.Date,
		})
	}

	if err := productsCursor.Err(); err != nil {
		return payments, err
	}

	return payments, nil
}
