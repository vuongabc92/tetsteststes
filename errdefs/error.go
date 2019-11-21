package errdefs

// Page Not Found Error
type PageNotFoundError struct{}

func (PageNotFoundError) Error() string {
	return "page not found"
}

func (PageNotFoundError) NotFound() {}

// Forbidden Error
type ForbiddenError struct{}

func (ForbiddenError) Error() string {
	return "403 Forbidden"
}

func (ForbiddenError) Forbidden() {}
