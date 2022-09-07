package UuidUtils

import "testing"

func TestNewUUID(t *testing.T) {
	t.Log(NewUUID())
	t.Log(NewSpaceUUID("1"))
	t.Log(NewSpaceUUID("1"))
	t.Log(NewSpaceUUID("2"))
}
