package openXml

import "testing"

func TestAbs(t *testing.T) {
	got := 1
	want := 0
	if got != want {
		t.Errorf("want %v; go %v", want, got)
	}
}