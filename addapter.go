package bye

// The Func type is an adapter to allow the use of ordinary functions as ShutdownableService.
type Func func()

// Shutdown calls wrapped f().
func (f Func) Shutdown() error {
	f()
	return nil
}

// The ErrFunc type is an adapter to allow the use of ordinary functions as ShutdownableService.
type ErrFunc func() error

// Shutdown calls wrapped f().
func (f ErrFunc) Shutdown() error {
	return f()
}
