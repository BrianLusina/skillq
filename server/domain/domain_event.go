package sharedkernel

import (
	"time"
)

type (
	DomainEvent interface {
		CreatedAt() time.Time
		Identity() string
	}
)
