package bye

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"go.szostok.io/bye/internal/multierror"
)

// ShutdownableService represents a service that supports graceful shutdown pattern.
type ShutdownableService interface {
	Shutdown() error
}

// ParentService aggregates services that are Shutdownable.
// Those services are registered in parent and shutdown is cascaded to them.
type ParentService struct {
	// children is a list of services which should be shut down together.
	// Each service is implementing ShutdownableService interface and may be shut-down in parallel.
	children []ShutdownableService

	// timeout defines the max waiting time for all registered children to shut down.
	// If set, we wait with a given timeout. Otherwise, there is no time limit - it's scheduled on app shutdown level,
	// so leaking goroutines is not a problem.
	timeout time.Duration
}

func NewParentService(opts ...Option) *ParentService {
	svc := &ParentService{}

	for _, mutate := range opts {
		mutate(svc)
	}

	return svc
}

// Register is registering dependent service to be shutdown on parent service shutdown.
func (s *ParentService) Register(child ShutdownableService) *ParentService {
	s.children = append(s.children, child)
	return s
}

// Shutdown is called to trigger shutdown of all associated children. It waits for all children.
//
// Child shutdown is considered successful also in cases when context.Cancelled error is returned.
func (s *ParentService) Shutdown() error {
	childShutdownFeedback := make(chan error, len(s.children))

	ctx, cancel := s.getContextWithOptionalTimeout()
	defer cancel()

	wg, _ := errgroup.WithContext(ctx)

	// trigger shutdown
	for _, child := range s.children {
		child := child
		wg.Go(func() error {
			childShutdownFeedback <- child.Shutdown()
			return nil // if we will return error then wg.Wait will be automatically release, but we don't want to.
		})
	}

	// Wait for all children to shut down.
	result := multierror.New()
	result = multierror.Append(result, wg.Wait())

	// At this point we are sure that all children responded.
	close(childShutdownFeedback)

	// produce single result
	for err := range childShutdownFeedback {
		if err == nil || err == context.Canceled {
			continue
		}
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

func (s *ParentService) getContextWithOptionalTimeout() (context.Context, context.CancelFunc) {
	if s.timeout == 0 {
		return context.WithCancel(context.Background())
	}
	return context.WithTimeout(context.Background(), s.timeout)
}
