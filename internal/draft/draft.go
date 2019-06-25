package draft

import (
	"regexp"
	"time"
)

const (
	mdTyp = "markdown"
)

type draft struct {
	PID     string    `json:"pid"`
	Changed time.Time `json:"changed"`
	Cont    content   `json:"cont"`
}

type content struct {
	Typ  string `json:"typ"`
	Body string `json:"body"`
}

type draftItem struct {
	PID     string    `json:"pid"`
	Changed time.Time `json:"changed"`
	Title   string    `json:"title"`
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
