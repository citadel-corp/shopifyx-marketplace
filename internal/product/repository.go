package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	List(ctx context.Context, filter ListProductPayload) ([]Product, *response.Pagination, error)
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
			product.Name, product.ImageURL, product.Stock, product.Condition, pq.Array(product.Tags), product.IsPurchasable, product.Price, product.User.ID)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (d *DBRepository) List(ctx context.Context, filter ListProductPayload) ([]Product, *response.Pagination, error) {
	var products []Product
	var pagination *response.Pagination

	var (
		selectStatement     string
		whereStatement      string
		query               string
		joinStatement       string
		orderStatement      string
		paginationStatement string
		args                []interface{}
		columnCtr           int = 1
	)

	if filter.UserOnly && filter.UserID != 0 {
		whereStatement = fmt.Sprintf("%s WHERE products.user_id = $%d", whereStatement, columnCtr)
		joinStatement = fmt.Sprintf("%s JOIN users ON users.id = products.user_id", joinStatement)
		args = append(args, filter.UserID)
		columnCtr++
	}

	if len(filter.Tags) > 0 {
		for i := range filter.Tags {
			whereStatement = insertWhereStatement(i > 0, whereStatement)
			whereStatement = fmt.Sprintf("%s $%d = ANY(products.tags)", whereStatement, columnCtr)
			args = append(args, filter.Tags[i])
			columnCtr++
		}
	}

	if filter.Condition != "" {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s products.condition = $%d", whereStatement, columnCtr)
		args = append(args, filter.Condition)
		columnCtr++
	}

	if !filter.ShowEmptyStock {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s products.stock > $%d", whereStatement, columnCtr)
		args = append(args, 0)
		columnCtr++
	}

	if filter.MinPrice > 0 {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s products.price > $%d", whereStatement, columnCtr)
		args = append(args, filter.MinPrice)
		columnCtr++
	}

	if filter.MaxPrice > 0 {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s products.price < $%d", whereStatement, columnCtr)
		args = append(args, filter.MaxPrice)
		columnCtr++
	}

	if filter.Search != "" {
		whereStatement = insertWhereStatement(len(args) > 0, whereStatement)
		whereStatement = fmt.Sprintf("%s lower(products.name) LIKE CONCAT('%%',$%d::text,'%%')", whereStatement, columnCtr)
		args = append(args, strings.ToLower(filter.Search))
		columnCtr++
	}

	var orderBy string
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "asc":
			orderBy = "asc"
		case "dsc":
			orderBy = "desc"
		}
	}

	if filter.SortBy != "" {
		switch filter.SortBy {
		case productSortBy(SortByPrice):
			orderStatement = fmt.Sprintf("%s ORDER BY products.price %s", orderStatement, orderBy)
		case productSortBy(SortByDate):
			orderStatement = fmt.Sprintf("%s ORDER BY products.created_at %s", orderStatement, orderBy)
		}
	}

	var rows *sql.Rows
	var err error
	pagination = &response.Pagination{
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	if filter.Limit != 0 {
		selectStatement = fmt.Sprintf(`
			SELECT COUNT(*) OVER() AS total_count, products.uid as productId, products.name as name, products.image_url as imageUrl, 
				products.stock as stock, products.condition as condition, products.tags as tags, products.is_purchaseable as isPurchasable, 
				products.price as price, products.purchase_count as purchaseCount 
			FROM products 
		%s`, selectStatement)

		paginationStatement = fmt.Sprintf("%s LIMIT $%d", paginationStatement, columnCtr)
		args = append(args, filter.Limit)
		columnCtr++

		paginationStatement = fmt.Sprintf("%s OFFSET $%d", paginationStatement, columnCtr)
		args = append(args, filter.Offset)

		query = fmt.Sprintf("%s %s %s %s %s;", selectStatement, joinStatement, whereStatement, orderStatement, paginationStatement)

		// sanitize query
		query = strings.Replace(query, "\t", "", -1)
		query = strings.Replace(query, "\n", "", -1)

		rows, err = d.db.DB().QueryContext(ctx, query, args...)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			var p Product
			if err := rows.Scan(&pagination.Total, &p.UUID, &p.Name, &p.ImageURL, &p.Stock, &p.Condition,
				pq.Array(&p.Tags), &p.IsPurchasable, &p.Price, &p.PurchaseCount); err != nil {
				return products, nil, err
			}
			products = append(products, p)
		}
	} else {
		selectStatement = fmt.Sprintf(`
			SELECT products.uid as productId, products.name as name, products.image_url as imageUrl, 
				products.stock as stock, products.condition as condition, products.tags as tags, products.is_purchaseable as isPurchasable, 
				products.price as price, products.purchase_count as purchaseCount 
			FROM products 
		%s`, selectStatement)

		query = fmt.Sprintf("%s %s %s %s %s;", selectStatement, joinStatement, whereStatement, orderStatement, paginationStatement)

		// sanitize query
		query = strings.Replace(query, "\t", "", -1)
		query = strings.Replace(query, "\n", "", -1)

		rows, err = d.db.DB().QueryContext(ctx, query, args...)
		if err != nil {
			return nil, nil, err
		}

		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.UUID, &p.Name, &p.ImageURL, &p.Stock, &p.Condition,
				pq.Array(&p.Tags), &p.IsPurchasable, &p.Price, &p.PurchaseCount); err != nil {
				return products, nil, err
			}
			products = append(products, p)
		}

		pagination.Total = len(products)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return products, nil, err
	}

	return products, pagination, nil
}

func insertWhereStatement(condition bool, statement string) string {
	if condition {
		return fmt.Sprintf(`%v AND`, statement)
	}
	return fmt.Sprintf(`%v WHERE`, statement)
}
