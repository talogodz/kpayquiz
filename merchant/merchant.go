package merchant

import (
	"errors"
	"kpayquiz/database"

	"gopkg.in/mgo.v2/bson"
)

const (
	C_MERCHANT = "merchant"
	C_PRODUCT  = "product"
)

type Merchant struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Username    string        `bson:"username" json:"-"`
	Password    string        `bson:"password" json:"-"`
	BankAccount string        `bson:"bank_account" json:"bank_account"`
}

type Product struct {
	MerchantID    bson.ObjectId `bson:"merchant_id" json:"merchant_id"`
	ID            bson.ObjectId `bson:"_id" json:"id"`
	Name          string        `bson:"name" json:"name"`
	Amount        float64       `bson:"amount" json:"amount"`
	AmountChange  float64       `bson:"amount_change" json:"amount_change"`
	SellingVolume int           `bson:"selling_volume" json:"selling_volume"`
}

func FindById(id string) (Merchant, error) {
	var merchant Merchant
	if bson.IsObjectIdHex(id) {
		err := database.DB.C(C_MERCHANT).FindId(bson.ObjectIdHex(id)).One(&merchant)
		return merchant, err
	}
	return merchant, errors.New("Invalid IdObject")
}

func IsBankAccountExists(bankaccount string) (bool, error) {
	count, err := database.DB.C(C_MERCHANT).Find(bson.M{"bank_account": bankaccount}).Count()
	return count > 0, err
}

func Insert(merchant Merchant) error {
	err := database.DB.C(C_MERCHANT).Insert(&merchant)
	return err
}

func Delete(merchant Merchant) error {
	err := database.DB.C(C_MERCHANT).Remove(&merchant)
	return err
}

func Update(merchant Merchant) error {
	err := database.DB.C(C_MERCHANT).UpdateId(merchant.ID, &merchant)
	return err
}

/// PRODUCT
func FindProductById(id string, merchantID bson.ObjectId) (Product, error) {
	var product Product
	//err := database.DB.C(C_PRODUCT).FindId(bson.ObjectIdHex(id)).Select(bson.M{"merchant_id": merchantID}).One(&product)
	err := database.DB.C(C_PRODUCT).Find(bson.M{"_id": bson.ObjectIdHex(id), "merchant_id": merchantID}).One(&product)

	return product, err
}

func FindSingleProduct(id string) (Product, error) {
	var product Product
	if bson.IsObjectIdHex(id) {
		err := database.DB.C(C_PRODUCT).FindId(bson.ObjectIdHex(id)).One(&product)
		return product, err
	}
	return product, errors.New("Invalid IdObject")
}

func InsertProduct(product Product) error {
	err := database.DB.C(C_PRODUCT).Insert(&product)
	return err
}

func DeleteProduct(product Product) error {
	err := database.DB.C(C_PRODUCT).Remove(&product)
	return err
}

func UpdateProduct(product Product) error {
	err := database.DB.C(C_PRODUCT).UpdateId(product.ID, &product)
	return err
}

func FindProductsFromMerchantID(id bson.ObjectId) ([]Product, error) {
	var products []Product
	err := database.DB.C(C_PRODUCT).Find(bson.M{"merchant_id": id}).All(&products)
	return products, err
}
