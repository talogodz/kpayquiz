package merchant

import (
	"kpayquiz/rand"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/binding"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

func getMerchant(c *gin.Context) Merchant {
	val, found := c.Get("merchant")
	if !found {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Merchant not found",
		})
		return Merchant{}
	}
	return val.(Merchant)
}
func RegisterEndPoint(c *gin.Context) {
	var req struct {
		BankAccount string `bson:"bank_account" json:"bank_account"  binding:"required"`
		Name        string `bson:"name" json:"name"  binding:"required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	exists, err := IsBankAccountExists(req.BankAccount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bank Account Is Already Used",
		})
		return
	}

	merchant := Merchant{
		ID:          bson.NewObjectId(),
		Name:        req.Name,
		BankAccount: req.BankAccount,
		Username:    rand.RandomHex(20),
		Password:    rand.RandomHex(20),
	}

	err = Insert(merchant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	type response struct {
		ID       bson.ObjectId `json="id"`
		Username string        `json="username"`
		Password string        `json="password"`
	}
	res := response{
		merchant.ID, merchant.Username, merchant.Password,
	}
	c.JSON(http.StatusOK, res)
}

func InfoEndPoint(c *gin.Context) {
	merchant := getMerchant(c)
	c.JSON(http.StatusOK, merchant)
}

func UpdateEndPoint(c *gin.Context) {
	var req struct {
		Name string `bson="name" json="name" binding:"required" `
	}
	merchant := getMerchant(c)
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if req.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "info is missing",
		})
		return
	}
	merchant.Name = req.Name
	err = Update(merchant)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

func ProductsEndPoint(c *gin.Context) {
	merchant := getMerchant(c)
	products, err := FindProductsFromMerchantID(merchant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if products == nil {
		products = make([]Product, 0)
	}
	c.JSON(http.StatusOK, products)
}

func AddProductEndPoint(c *gin.Context) {
	var req struct {
		Name   string  `bson:"name" json:"name" binding="required"`
		Amount float64 `bson:"amount" json:"amount"  binding="required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	merchant := getMerchant(c)
	products, err := FindProductsFromMerchantID(merchant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(products) >= 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "You can have only 5 products",
		})
		return
	}
	product := Product{
		MerchantID:   merchant.ID,
		ID:           bson.NewObjectId(),
		Name:         req.Name,
		Amount:       req.Amount,
		AmountChange: req.Amount,
	}

	err = InsertProduct(product)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func UpdateProductEndPoint(c *gin.Context) {
	product_id := c.Param("product_id")
	merchant := getMerchant(c)
	product, err := FindProductById(product_id, merchant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req struct {
		Amount float64 `bson:"amount" json:"amount"  binding="required"`
	}
	err = c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("Product Amount = %f", product.Amount)
	product.AmountChange = product.Amount
	product.Amount = req.Amount
	err = UpdateProduct(product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

func BuyProductEndPoint(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id" binding="required"`
		Volume    int    `json:"volume"  binding="required"`
	}
	err := c.Bind(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	product, err := FindSingleProduct(req.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product.SellingVolume += req.Volume
	err = UpdateProduct(product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

func RemoveProductEndPoint(c *gin.Context) {
	product_id := c.Param("product_id")
	merchant := getMerchant(c)
	product, err := FindProductById(product_id, merchant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if product.SellingVolume > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Cannot product that have been sold",
		})
		return
	}
	err = DeleteProduct(product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func SellReportEndPoint(c *gin.Context) {
	merchant := getMerchant(c)
	products, err := FindProductsFromMerchantID(merchant.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if products == nil || len(products) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Products not found",
		})
		return
	}
	type responseProduct struct {
		Name          string `json="name"`
		SellingVolume int    `json="selling_volume"`
	}
	var res struct {
		Date       time.Time         `json="date"`
		Products   []responseProduct `json="products"`
		Accumulate float64           `json="accumulate"`
	}
	var sellList []responseProduct
	accum := 0.0
	for _, product := range products {
		if product.SellingVolume > 0 {
			sellList = append(sellList, responseProduct{
				product.Name, product.SellingVolume,
			})
			accum += float64(product.SellingVolume) * product.AmountChange
		}
	}
	if sellList == nil || len(sellList) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No one buy your product yet",
		})
		return
	}
	res.Date = time.Now()
	res.Products = sellList
	res.Accumulate = accum
	c.JSON(http.StatusOK, res)
}

func BasicAuth(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "id missing",
		})
		return
	}

	var req struct {
		Username string `json="username" binding="required"`
		Password string `json="password" binding="required"`
	}
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	merchant, err := FindById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if merchant.Password != req.Password || merchant.Username != req.Username {
		c.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, gin.H{
			"error": "Username & Password  missing or incorrect",
		})
		return
	}
	c.Set("merchant", merchant)
	c.Next()
}
