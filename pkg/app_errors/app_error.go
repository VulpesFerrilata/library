package app_errors

type AppError interface {
	error
	WebError
	GrpcError
	WebsocketError
}
