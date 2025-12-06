package envutil

import (
	"fmt"
	"strconv"

	"github.com/apotourlyan/ludus-studii/pkg/envutil/envvar"
)

type machineID struct {
	value int64
}

func MachineID(provider Provider) Variable[int64] {
	idstr := provider.Get(envvar.MachineId)
	value := parseMachineId(idstr)
	return &machineID{value}
}

func parseMachineId(s string) int64 {
	if s == "" {
		message := fmt.Sprintf("%q environment var not set", envvar.MachineId)
		panic(message)
	}

	id, err := strconv.ParseUint(s, 10, 10)
	if err != nil {
		message := fmt.Sprintf(
			"%q must be a valid 10-bit integer: %v",
			envvar.MachineId,
			err)
		panic(message)
	}

	return int64(id)
}

func (mid *machineID) Value() int64 {
	return mid.value
}
