package errdefs

// ErrNotFound signals that the requested object doesn't exist
type ErrNotFound interface {
	NotFound()
}

// ErrForbidden signals that the requested action cannot be performed under any circumstances.
// When a ErrForbidden is returned, the caller should never retry the action.
type ErrForbidden interface {
	Forbidden()
}
