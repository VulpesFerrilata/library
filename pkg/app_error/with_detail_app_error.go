package app_error

type WithDetailAppError interface {
	error
	AppError
	AddDetailError(detailErr DetailError)
}
