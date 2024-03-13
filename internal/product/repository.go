package product

import (
	"context"
	"database/sql"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
}

type DBRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (d *DBRepository) Create(ctx context.Context, product *Product) error {
	err := d.db.StartTx(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO products (
				name, image_url, stock, condition, tags, is_purchaseable, price, user_id
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8
			)`,
			product.Name, product.ImageURL, product.Stock, product.Condition.String(), pq.Array(product.Tags), product.IsPurchasable, product.Price, product.User.ID)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
