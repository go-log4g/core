package impl_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

type callerTestService struct {
}

func (this *callerTestService) PointerMethod() {
}

func (this callerTestService) ValueMethod() {
}

func callerTestFunction() {
}

func TestCallerResolver(test *testing.T) {
	resolver := impl.NewCallerResolver()

	pointerCaller := resolver.Resolve(reflect.ValueOf((*callerTestService).PointerMethod).Pointer())
	require.True(test, strings.HasSuffix(pointerCaller.Logger, "/callerTestService"))
	require.Equal(test, "PointerMethod", pointerCaller.Method)

	valueCaller := resolver.Resolve(reflect.ValueOf(callerTestService.ValueMethod).Pointer())
	require.True(test, strings.HasSuffix(valueCaller.Logger, "/callerTestService"))
	require.Equal(test, "ValueMethod", valueCaller.Method)

	functionCaller := resolver.Resolve(reflect.ValueOf(callerTestFunction).Pointer())
	require.False(test, strings.HasSuffix(functionCaller.Logger, "/callerTestService"))
	require.Equal(test, "callerTestFunction", functionCaller.Method)
}
