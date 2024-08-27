package store

import (
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/opplieam/bb-product-server/.gen/buy-better-core/public/model"
	. "github.com/opplieam/bb-product-server/.gen/buy-better-core/public/table"
)

type ProductStore struct {
	DB *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{
		DB: db,
	}
}

type ImageProductCustom struct {
	ID       int32  `sql:"primary_key" alias:"image_product.id"`
	ImageURL string `alias:"image_product.image_url"`
}

type PriceNowCustom struct {
	ID        int32     `sql:"primary_key" alias:"price_now.id"`
	Price     float64   `alias:"price_now.price"`
	Currency  string    `alias:"price_now.currency"`
	CreatedAt time.Time `alias:"price_now.created_at"`
}

type ResultGetAllProducts struct {
	ID   int32  `sql:"primary_key" alias:"products.id"`
	Name string `alias:"products.name"`
	URL  string `alias:"products.url"`
	model.Sellers
	ImageProduct []ImageProductCustom
	PriceNow     []PriceNowCustom
}

func (s *ProductStore) GetAllProducts(userID uint32) ([]ResultGetAllProducts, error) {
	var dest []ResultGetAllProducts
	stmt := SELECT(
		Products.AllColumns, Sellers.AllColumns,
		ImageProduct.AllColumns.Except(ImageProduct.ProductID),
		PriceNow.AllColumns.Except(PriceNow.ProductID),
	).FROM(UserSubProduct.
		INNER_JOIN(GroupProduct, GroupProduct.ID.EQ(UserSubProduct.GroupProductID)).
		INNER_JOIN(MatchProductGroup, MatchProductGroup.GroupID.EQ(GroupProduct.ID)).
		INNER_JOIN(Products, Products.ID.EQ(MatchProductGroup.ProductID)).
		INNER_JOIN(Sellers, Sellers.ID.EQ(Products.SellerID)).
		INNER_JOIN(ImageProduct, ImageProduct.ProductID.EQ(Products.ID)).
		INNER_JOIN(PriceNow, PriceNow.ProductID.EQ(Products.ID)),
	).
		WHERE(UserSubProduct.UserID.EQ(Int32(int32(userID))))
	err := stmt.Query(s.DB, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
