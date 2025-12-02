package idutil

import (
	"testing"
	"time"

	"github.com/apotourlyan/ludus-studii/pkg/syncutil"
	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

type mockTimeProvider struct {
	timestamp int64
}

func (m *mockTimeProvider) Now() time.Time {
	return time.Unix(0, m.timestamp*int64(time.Millisecond))
}

func TestSnowflakeID_Next_ContinuosSequence(t *testing.T) {
	snowflake := NewGenerator(
		&mockTimeProvider{timestamp: 1700000000000},
		syncutil.NewCounter(),
		64,
	)

	initial := (1700000000000 << 22) | (64 << 12) | int64(0)

	for i := range 5 {
		got := snowflake.Next()
		testutil.GotWant(t, got, initial+int64(i))
	}
}

func TestSnowflakeID_Next_SequenceReset(t *testing.T) {
	provider := &mockTimeProvider{timestamp: 1700000000000}
	snowflake := NewGenerator(
		provider,
		syncutil.NewCounter(),
		64,
	)

	initial := (1700000000000 << 22) | (64 << 12) | int64(0)

	for i := range 5 {
		got := snowflake.Next()
		testutil.GotWant(t, got, initial+int64(i))
	}

	provider.timestamp = 1700000000001
	initial = (1700000000001 << 22) | (64 << 12) | int64(0)

	for i := range 5 {
		got := snowflake.Next()
		testutil.GotWant(t, got, initial+int64(i))
	}
}

func TestToSnowflakeID_Success(t *testing.T) {
	cases := []struct {
		name      string
		timestamp int64
		instance  int64
		sequence  int64
		want      int64
	}{
		{
			name:      "all zeros",
			timestamp: 0,
			instance:  0,
			sequence:  0,
			want:      0,
		},
		{
			name:      "only timestamp",
			timestamp: 1,
			instance:  0,
			sequence:  0,
			want:      1 << 22,
		},
		{
			name:      "only instance",
			timestamp: 0,
			instance:  1,
			sequence:  0,
			want:      1 << 12,
		},
		{
			name:      "only sequence",
			timestamp: 0,
			instance:  0,
			sequence:  1,
			want:      1,
		},
		{
			name:      "all ones in each field",
			timestamp: 1,
			instance:  1,
			sequence:  1,
			want:      (1 << 22) | (1 << 12) | 1,
		},
		{
			name:      "max timestamp (41 bits)",
			timestamp: 0x1FFFFFFFFFF,
			instance:  0,
			sequence:  0,
			want:      0x1FFFFFFFFFF << 22,
		},
		{
			name:      "max instance (10 bits)",
			timestamp: 0,
			instance:  0x3FF,
			sequence:  0,
			want:      0x3FF << 12,
		},
		{
			name:      "max sequence (12 bits)",
			timestamp: 0,
			instance:  0,
			sequence:  0xFFF,
			want:      0xFFF,
		},
		{
			name:      "realistic id",
			timestamp: 1700000000000, // realistic Unix milliseconds
			instance:  512,
			sequence:  100,
			want:      (1700000000000 << 22) | (512 << 12) | 100,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := toSnowflakeID(c.timestamp, c.instance, c.sequence)
			testutil.GotWant(t, got, c.want)
		})
	}
}

func TestToSnowflakeID_NoOverlap(t *testing.T) {
	// Ensure fields don't overlap by checking that OR-ing max values
	// produces the expected result
	timestamp := int64(0x1FFFFFFFFFF) // 41 bits
	instance := int64(0x3FF)          // 10 bits
	sequence := int64(0xFFF)          // 12 bits

	result := toSnowflakeID(timestamp, instance, sequence)

	// Extract each field back
	xSequence := result & 0xFFF
	xInstance := (result >> 12) & 0x3FF
	xTimestamp := result >> 22

	testutil.GotWant(t, xSequence, sequence)
	testutil.GotWant(t, xInstance, instance)
	testutil.GotWant(t, xTimestamp, timestamp)
}

func TestToSnowflakeID_Panics(t *testing.T) {
	cases := []struct {
		name      string
		timestamp int64
		instance  int64
		sequence  int64
		wantPanic string
	}{
		{
			name:      "negative timestamp",
			timestamp: -1,
			instance:  0,
			sequence:  0,
			wantPanic: `"timestamp" must be >= 0, got -1`,
		},
		{
			name:      "timestamp overflow",
			timestamp: 0x1FFFFFFFFFF + 1,
			instance:  0,
			sequence:  0,
			wantPanic: `"timestamp" must be <= 2199023255551, got 2199023255552`,
		},
		{
			name:      "negative instance",
			timestamp: 0,
			instance:  -1,
			sequence:  0,
			wantPanic: `"instance" must be >= 0, got -1`,
		},
		{
			name:      "instance overflow",
			timestamp: 0,
			instance:  0x3FF + 1,
			sequence:  0,
			wantPanic: `"instance" must be <= 1023, got 1024`,
		},
		{
			name:      "negative sequence",
			timestamp: 0,
			instance:  0,
			sequence:  -1,
			wantPanic: `"sequence" must be >= 0, got -1`,
		},
		{
			name:      "sequence overflow",
			timestamp: 0,
			instance:  0,
			sequence:  0xFFF + 1,
			wantPanic: `"sequence" must be <= 4095, got 4096`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testutil.GotWantPanic(t, func() {
				toSnowflakeID(c.timestamp, c.instance, c.sequence)
			}, c.wantPanic)
		})
	}
}
