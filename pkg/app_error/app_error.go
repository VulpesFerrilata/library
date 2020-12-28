package app_error

type AppError interface {
	error
	WebError
	GrpcError
	WebsocketError
}
