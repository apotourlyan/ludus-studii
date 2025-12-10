package envutil

import (
	"testing"

	"github.com/apotourlyan/ludus-studii/pkg/testutil"
)

type mockProvider struct {
	value string
}

func (m *mockProvider) Get(varname string) string {
	return m.value
}

func TestMachineID_Get_Success(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		envVal string
		want   int64
	}{
		{
			name:   "valid zero",
			envVal: "0",
			want:   0,
		},
		{
			name:   "valid one",
			envVal: "1",
			want:   1,
		},
		{
			name:   "valid max (1023)",
			envVal: "1023",
			want:   1023,
		},
		{
			name:   "valid mid-range value",
			envVal: "512",
			want:   512,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			mid := MachineID(&mockProvider{value: c.envVal})
			got := mid.Value()
			testutil.GotWant(t, got, c.want)
		})
	}
}

func TestMachineID_NewMachineID_Panics(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		envVal    string
		wantPanic string
	}{
		{
			name:      "empty env var",
			envVal:    "",
			wantPanic: `"MACHINE_ID" environment var not set`,
		},
		{
			name:      "invalid integer",
			envVal:    "abc",
			wantPanic: `"MACHINE_ID" must be a valid 10-bit integer: strconv.ParseUint: parsing "abc": invalid syntax`,
		},
		{
			name:      "negative value",
			envVal:    "-1",
			wantPanic: `"MACHINE_ID" must be a valid 10-bit integer: strconv.ParseUint: parsing "-1": invalid syntax`,
		},
		{
			name:      "overflow value",
			envVal:    "1024",
			wantPanic: `"MACHINE_ID" must be a valid 10-bit integer: strconv.ParseUint: parsing "1024": value out of range`,
		},
		{
			name:      "large overflow value",
			envVal:    "999999",
			wantPanic: `"MACHINE_ID" must be a valid 10-bit integer: strconv.ParseUint: parsing "999999": value out of range`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			testutil.GotWantPanic(t, func() {
				MachineID(&mockProvider{value: c.envVal})
			}, c.wantPanic)
		})
	}
}
