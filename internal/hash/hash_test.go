package hash

import (
	"testing"
)

func TestGetContentId(t *testing.T) {
	content := "hello world"
	expected := "DULfJyE3WQqNxy3ymuhAChyNR3yufT88pmqvAazKFMG4"

	cid := GetContentId(content)

	if cid != expected {
		t.Errorf("expected: %s, got: %s\n", expected, cid)
	}
}
