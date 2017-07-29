package ep

// Waiter is a wait-only interface to sync.WaitGroup
type Waiter interface {
	Wait()
}
