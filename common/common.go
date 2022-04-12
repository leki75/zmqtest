package common

import "testing"

type TestFunc func(*testing.B, int) func(*testing.B)

var Msg = []byte("12345678901234567890123456789012345678901234567890")
