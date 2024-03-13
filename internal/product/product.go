package product

import (
	"time"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/twinj/uuid"
)

type Product struct {
	ID            int
	UUID          uuid.UUID
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

var Conditions []Condition = []Condition{New, Second}

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
