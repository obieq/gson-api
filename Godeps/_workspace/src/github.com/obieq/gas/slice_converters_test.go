package gas_test

import (
	"testing"

	"github.com/obieq/gas"
)

func TestInterfacesToStrings(t *testing.T) {
	interfacesToStringsTestHelper(t, []interface{}{float64(64), float32(32), 6, int32(32), int64(64), "string1", true, false})
}

func TestInterfacesToInts(t *testing.T) {
	interfacesToIntsTestHelper(t, []interface{}{float64(5), 6, "4"})
}

func TestIntsToInterfaces(t *testing.T) {
	var test gas.Ints = []int{2, 4, 6}

	// convert to interface slice
	interfaces, err := test.ToInterfaces()

	if err != nil {
		t.Error("Error should not have occurred:", err)
	}

	if l := len(interfaces); l != len(test) {
		t.Error("interfaces length is incorrect:", l)
	}

	// convert back to int slice
	interfacesToIntsTestHelper(t, interfaces)
}

func interfacesToStringsTestHelper(t *testing.T, interfaces []interface{}) {
	var test gas.Interfaces = interfaces

	ints, err := test.ToStrings()

	if err != nil {
		t.Error("Error should not have occurred:", err)
	}

	if l := len(ints); l != len(interfaces) {
		t.Error("ints length is incorrect:", l)
	}
}

func interfacesToIntsTestHelper(t *testing.T, interfaces []interface{}) {
	var test gas.Interfaces = interfaces

	ints, err := test.ToInts()

	if err != nil {
		t.Error("Error should not have occurred:", err)
	}

	if l := len(ints); l != len(interfaces) {
		t.Error("ints length is incorrect:", l)
	}
}
