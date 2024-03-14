package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	List(ctx context.Context, filter ListProductPayload) ([]Product, error)
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

func (d *DBRepository) List(ctx context.Context, filter ListProductPayload) ([]Product, error) {
	var products []Product

	query := ``
	statement := `SELECT products.uid as productId, products.name as name, products.image_url as imageUrl, 
		products.stock as stock, products.condition as condition, products.tags as tags, products.is_purchaseable as isPurchaseable, 
		products.price as price, products.purchase_count as purchaseCount
		FROM products`
	joinStatement := ``
	orderStatement := ``
	paginationStatement := ``
	args := []interface{}{}

	if filter.UserOnly && filter.UserID != 0 {
		statement = fmt.Sprintf(`%s WHERE products.user_id = ?`, statement)
		joinStatement = fmt.Sprintf(`%s JOIN users ON users.id = products.user_id`, joinStatement)
		args = append(args, filter.UserID)
	}

	if len(filter.Tags) > 0 {
		statement = insertAndToStatement(len(args) > 0, statement)

		for i := range filter.Tags {
			statement = insertAndToStatement(i > 0, statement)
			statement = fmt.Sprintf(`%s WHERE ? = ANY(products.tags)`, statement)
			args = append(args, filter.Tags[i])
		}
	}

	if filter.Condition != "" {
		statement = insertAndToStatement(len(args) > 0, statement)
		statement = fmt.Sprintf(`%s WHERE products.condition = ?`, statement)
		args = append(args, filter.Condition)
	}

	if !filter.ShowEmptyStock {
		statement = insertAndToStatement(len(args) > 0, statement)
		statement = fmt.Sprintf(`%s WHERE products.stock > ?`, statement)
		args = append(args, 0)
	}

	if filter.MinPrice > 0 {
		statement = insertAndToStatement(len(args) > 0, statement)
		statement = fmt.Sprintf(`%s WHERE products.price > ?`, statement)
		args = append(args, filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		statement = insertAndToStatement(len(args) > 0, statement)
		statement = fmt.Sprintf(`%s WHERE products.price < ?`, statement)
		args = append(args, filter.MaxPrice)
	}

	if filter.Search != "" {
		statement = insertAndToStatement(len(args) > 0, statement)
		statement = fmt.Sprintf(`%s WHERE products.name LIKE ?`, statement)
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
	}

	if filter.SortBy != "" {
		orderStatement = fmt.Sprintf(`%s ORDER BY ?`, orderStatement)
		args = append(args, "products."+filter.SortBy)
	}

	if filter.OrderBy != "" {
		orderStatement = fmt.Sprintf(`%s ?`, orderStatement)
		args = append(args, filter.OrderBy)
	}

	if filter.Limit != 0 && filter.Offset != 0 {
		paginationStatement = fmt.Sprintf(`%s LIMIT ?`, paginationStatement)
		args = append(args, filter.Limit)

		paginationStatement = fmt.Sprintf(`%s OFFSET ?`, paginationStatement)
		args = append(args, filter.Offset)
	}

	query = fmt.Sprintf(`%s %s %s %s`, statement, joinStatement, orderStatement, paginationStatement)

	rows, err := d.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.UUID, &p.ImageURL, &p.Stock, &p.Condition, &p.Tags,
			&p.IsPurchasable, &p.Price, &p.PurchaseCount); err != nil {
			return products, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}

	return products, nil
}

func insertAndToStatement(condition bool, statement string) string {
	if condition {
		return fmt.Sprintf(`%v AND`, statement)
	}
	return statement
}
