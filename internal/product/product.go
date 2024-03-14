package product

import (
	"time"

	"github.com/google/uuid"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type Product struct {
	ID            int
	UUID          uuid.UUID
	Name          string
	ImageURL      string
	Stock         int
	Condition     Condition
	Tags          []string
	IsPurchasable bool
	Price         int
	PurchaseCount int
	User          user.User
	CreatedAt     time.Time
}

type Condition string

const (
	New    Condition = "new"
	Second Condition = "second"
)

var Conditions []interface{} = []interface{}{New, Second}
