package function

import (
	"testing"

	"github.com/alicebob/miniredis"
)

func TestIncr(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	key := "hello" + separator + "thinking"
	s.Incr(key, 1)
	if got, err := s.Get(key); err != nil || got != "1" {
		t.Error("'key' has the wrong value")
	}
}
