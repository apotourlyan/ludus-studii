package idutil

import (
	"github.com/apotourlyan/ludus-studii/pkg/panicutil"
	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/timeutil"
)

type Generator interface {
	Next() int64
}

type snowflakeGenerator struct {
	time      timeutil.Provider
	counter   syncutil.Counter
	machineId int64
	last      int64 // Unix milliseconds timestamp of last id generation time
}

func NewGenerator(
	time timeutil.Provider, counter syncutil.Counter, machineId int64,
) Generator {
	var last int64
	return &snowflakeGenerator{time, counter, machineId, last}
}

func (s *snowflakeGenerator) Next() int64 {
	now := s.time.Now().UnixMilli()
	machineID := s.machineId
	sequence := s.counter.Next(func() bool {
		// Will be executed inside the mutex lock of the counter
		shouldReset := now != s.last
		if shouldReset {
			s.last = now
		}
		return shouldReset
	})

	return toSnowflakeID(now, machineID, sequence)
}

func toSnowflakeID(timestamp int64, instance int64, sequence int64) int64 {
	// Validate timestamp (41 bits max)
	panicutil.RequireNonNegative(timestamp, "timestamp")
	panicutil.RequireLessThanOrEqualTo(timestamp, 0x1FFFFFFFFFF, "timestamp")

	// Validate instance (10 bits max)
	panicutil.RequireNonNegative(instance, "instance")
	panicutil.RequireLessThanOrEqualTo(instance, 0x3FF, "instance")

	// Validate sequence (12 bits max)
	panicutil.RequireNonNegative(sequence, "sequence")
	panicutil.RequireLessThanOrEqualTo(sequence, 0xFFF, "sequence")

	// 64-bit number:
	// [1 bit unused/sign][41 bits timestamp][10 bits instance][12 bits sequence]
	//    ↑ no shift           ↑ shift 22        ↑ shift 12
	id := (timestamp << 22) | (instance << 12) | sequence
	return id
}
