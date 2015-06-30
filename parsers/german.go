package parsers

import (
	"errors"
	"strings"

	"dylanmackenzie.com/flashcard/word"

	"github.com/PuerkitoBio/goquery"
)

var (
	starts  = []int{4, 6, 16, 18}
	offsets = []int{0, 4, 8, 1, 5, 9}
)

func ParseGermanVerb(doc *goquery.Document, lang, title string) (*word.Word, error) {
	cons := make([]string, 0, 30)
	mwheadline := findHeadline(doc, lang, title)
	h3 := mwheadline.Parent()

	// Get all of the h4's after the 'Verb' heading until we get to the
	// next language heading. One of these will be the 'Conjugation'
	// heading
	conjHead := h3.NextFilteredUntil("h4", "h3")
	conjHead = conjHead.FilterFunction(func(i int, subhead *goquery.Selection) bool {
		mwheadline := subhead.Find(".mw-headline")
		return mwheadline.Text() == "Conjugation"
	})
	if conjHead.Size() != 1 {
		return nil, errors.New("Expected exactly one Conjugation heading")
	}

	tds := conjHead.NextFilteredUntil("div", "h4").Find("td")
	if tds.Size() < 28 {
		return nil, errors.New("Expected at least 28 fields in Conjugation Table")
	}

	tds.Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		if t != "" {
			cons = append(cons, t)
		}
	})

	w := &word.Word{
		Word: cons[0],
		Pos:  "verb",

		Attr: map[string]string{
			"presentParticiple": cons[1],
			"pastParticiple":    cons[2],
			"auxiliary":         cons[3],
			"imperative":        strings.Fields(cons[28])[0],
		},
	}

	tmp := make([]string, 6, 6)
	cts := []string{
		"present",
		"subjunctive1",
		"preterite",
		"subjunctive2",
	}

	for i, ct := range cts {
		for j, off := range offsets {
			phrase := cons[starts[i]+off]
			tmp[j] = strings.Fields(phrase)[1]
		}

		w.Attr[ct] = strings.Join(tmp, ",")
	}

	return w, nil
}
