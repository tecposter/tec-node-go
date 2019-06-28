package draft

import (
	"encoding/json"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/uuid"
	"testing"
)

func TestDraft(t *testing.T) {
	id, err := uuid.NewID()
	if err != nil {
		t.Fatal(err)
	}

	drft := newDrft(id, dto.MakeContent("markdown", "body"))
	t.Log(drft.PID.Base58(), drft.Changed, drft.Cont.Typ, drft.Cont.Body)

	drftBytes, err := drft.marshal()
	if err != nil {
		t.Fatal(err)
	}

	var ndf draft
	err = ndf.unmarshal(drftBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ndf.PID.Base58(), ndf.Changed, ndf.Cont.Typ, ndf.Cont.Body)

	if ndf.PID.Base64() != drft.PID.Base64() {
		t.Fatalf("PID expected: %s, received: %s", drft.PID.Base58(), ndf.PID.Base58())
	}

	byts, err := json.Marshal(ndf)
	t.Log(string(byts), err)
}
