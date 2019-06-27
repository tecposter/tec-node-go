package draft

import (
	"encoding/json"
	"fmt"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"regexp"
	"time"
)

const (
	mdTyp = "markdown"
)

type draft struct {
	PID     dto.ID    `json:"pid"`
	Changed time.Time `json:"changed"`
	Cont    content   `json:"cont"`
}

type content struct {
	Typ  string `json:"typ"`
	Body string `json:"body"`
}

type draftItem struct {
	PID     dto.ID    `json:"pid"`
	Changed time.Time `json:"changed"`
	Title   string    `json:"title"`
}

const (
	idSize   = 16
	timeSize = 15
)

func newDrft(id dto.ID, typ string, body string) *draft {
	return &draft{
		PID:     id,
		Changed: time.Now(),
		Cont: content{
			Typ:  typ,
			Body: body}}
}

func (d *draft) Marshal() ([]byte, error) {
	fmt.Println("Marshal:", d.PID.Bytes())

	id := d.PID.Bytes()
	changed, err := d.Changed.MarshalBinary()
	if err != nil {
		return nil, err
	}

	cont, err := json.Marshal([]string{
		d.Cont.Typ,
		d.Cont.Body})
	if err != nil {
		return nil, err
	}

	fmt.Printf("id: %d, changed: %d\n", len(id), len(changed))

	return append(id, append(changed, cont...)...), err
}

func (d *draft) Unmarshal(src []byte) error {
	id := dto.ID(src[0:idSize])
	var changed time.Time
	err := changed.UnmarshalBinary(src[idSize : idSize+timeSize])
	if err != nil {
		return err
	}
	var arr [2]string
	err = json.Unmarshal(src[idSize+timeSize:], &arr)
	if err != nil {
		return err
	}
	d.PID = id
	d.Changed = changed
	d.Cont.Typ = arr[0]
	d.Cont.Body = arr[1]

	return nil
}

func (d *draft) Title() string {
	if d.Cont.Typ == mdTyp {
		return extractMdTitle(d.Cont.Body)
	}

	if len(d.Cont.Body) > 100 {
		return d.Cont.Body[0:100]
	}
	return d.Cont.Body
}

func extractMdTitle(content string) string {
	//$matched = preg_match('/# ([^#\n]+)/', $content, $matches);
	re := regexp.MustCompile(`# ([^#\n]+)\n`)
	founds := re.FindStringSubmatch(content)

	//fmt.Println(founds)
	if len(founds) >= 2 {
		return founds[1]
	}

	l := len(content)
	if l > 100 {
		return content[0:100]
	}

	return content
}
