package innererror

// InnerErr is an interface satisfied by the error types which has inner error.
type InnerErr interface {
	GetInnerErr() error
}
