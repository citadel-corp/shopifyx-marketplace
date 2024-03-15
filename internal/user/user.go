package user

type User struct {
	ID               uint64
	Username         string
	Name             string
	ProductSoldTotal int
	HashedPassword   string
}
