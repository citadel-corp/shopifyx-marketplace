package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
)

type Repository interface {
	Create(product *Product) error
}

type DBRepository struct {
	db *db.DB
}

func (d *DBRepository) Create(product *Product) error {
	err := d.db.StartTx(context.TODO(), func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO products (image_url, stock, conditon, tags, is_purchaseable, price, user_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			product.ImageURL, product.Stock, product.Condition.String(), product.Tags, product.IsPurchasable, product.Price, product.User.ID, time.Now())
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
