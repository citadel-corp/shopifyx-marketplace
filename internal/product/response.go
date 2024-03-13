package product

type Response struct {
	Code    int
	Message string
	Data    any
}

var (
	SuccessResponse = Response{Code: 200, Message: "Product created successfully"}
)
