package dto

import (
	"encoding/json"
	"regexp"
)

// ContentType content type
type ContentType byte

// Content Type
const (
	TypText ContentType = iota
	TypMarkdown
	TypHTML
)

var typArr = []string{
	"text",
	"markdown",
	"html"}

// ParseContentType parses src to ContentType
func ParseContentType(src string) ContentType {
	for k, v := range typArr {
		if v == src {
			return ContentType(k)
		}
	}
	return TypText
}

func (c ContentType) String() string {
	return typArr[c]
}

// MarshalJSON implements the json.Marshaler interface
func (c ContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// A Content contains Typ and Body
type Content struct {
	Typ  ContentType `json:"typ"`
	Body string      `json:"body"`
}

// Title extracts title from Content
func (c *Content) Title(limit int) string {
	switch c.Typ {
	case TypMarkdown:
		return extractMDTitle(c.Body, limit)
	default:
		return extractTXTTitle(c.Body, limit)
	}
}

func extractTXTTitle(content string, limit int) string {
	if len(content) > limit {
		return content[0:limit]
	}
	return content
}

func extractMDTitle(content string, limit int) string {
	//$matched = preg_match('/# ([^#\n]+)/', $content, $matches);
	re := regexp.MustCompile(`# ([^#\n]+)\n`)
	founds := re.FindStringSubmatch(content)

	//fmt.Println(founds)
	if len(founds) >= 2 {
		return founds[1]
	}

	l := len(content)
	if l > limit {
		return content[0:limit]
	}

	return content
}
