package util

func pointer[T any](x T) *T {
	return &x
}
