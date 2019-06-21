package scopes

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string `readScope:"read:A" writeScope:"write:A"`
	B string `readScope:"read:B" writeScope:"write:B"`
}

func TestOperationAllowed(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsOperationAllowed(testStruct{}, "A", READ, []string{"read:A"}))
	assert.True(IsOperationAllowed(testStruct{}, "A", WRITE, []string{"write:A"}))
	assert.False(IsOperationAllowed(testStruct{}, "A", READ, []string{"read:B"}))
	assert.False(IsOperationAllowed(testStruct{}, "A", WRITE, []string{"write:B"}))

	output := testStruct{
		A: "A",
		B: "B",
	}

	FilterRead(output, reflect.Indirect(reflect.ValueOf(&output)), []string{"read:A"})
	assert.Zero(output.B, "Outputs B Field Should Be Zero")
	assert.Equal("A", output.A)

}
