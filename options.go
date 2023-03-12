package bye

import "time"

type Option func(service *ParentService)

func WithTimeout(timeout time.Duration) Option {
	return func(service *ParentService) {
		service.timeout = timeout
	}
}
