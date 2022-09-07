package dateutils

import "testing"

func TestDate(t *testing.T) {
	t.Log(Date())
}

func TestDateTime(t *testing.T) {
	t.Log(DateTime())
}

func TestBeforeTime(t *testing.T) {
	t.Log(BeforeTime(1))
}
