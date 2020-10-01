package badgerdb

import (
	"testing"
)

func TestFmtDBKey(t *testing.T) {
	expected := "testPrefix:testID"
	actual := string(fmtDBKey("testPrefix", "testID"))
	if actual != expected {
		t.Errorf("fmtDBKey() = %v, want %v", actual, expected)
	}
}
