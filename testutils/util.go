package testutils

// WhenHandler is used in mocks and it is associated to a function name. If registered, when a
// function by a given name is triggered, all arguments received by the function will be relayed
// to the handler. The handler should respond as many arguments as the function has.
type WhenHandler func(args ...interface{}) []interface{}

// MakeReturn is a test utility function useful to convert multiple return values to an array
// of interface{}.
//
//   Usage:  Given the following function
//
//      SomeFunction() (int, string, bool, error) {
//          return 1, "hello", true, error.New("some error")
//      }
//
//   Doing the following:
//
//      result := MakeReturn(SomeFunction())
//
//   The value of `result` will be:
//
//       result = []interface{}{ 1, "hello", true, error.New("some error") }
//
func MakeReturn(retArgs ...interface{}) []interface{} {
	return retArgs
}
