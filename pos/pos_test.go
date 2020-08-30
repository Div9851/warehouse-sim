package pos

import "testing"

func TestPos(t *testing.T) {
	pos := New(1, 3)
	if pos.X != 1 {
		t.Fatalf("pos.X should be `1`, but `%v`", pos.X)
	}
	if pos.Y != 3 {
		t.Fatalf("pos.Y should be `3`, but `%v`", pos.Y)
	}
}
