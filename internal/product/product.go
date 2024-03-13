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

type Condition int64

const (
	New Condition = iota
	Second
)

func (c Condition) String() string {
	switch c {
	case New:
		return "new"
	case Second:
		return "second"
	default:
		return ""
	}
}

func ToCondition(s string) Condition {
	switch s {
	case "new":
		return New
	case "second":
		return Second
	default:
		return 0
	}
}
