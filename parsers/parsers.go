package parsers

import (
	"fmt"
	"regexp"

	"dylanmackenzie.com/flashcard/word"

	"github.com/PuerkitoBio/goquery"
)

type Parser interface {
	Parse(doc *goquery.Document, lang, title string) (*word.Word, error)
}

var List = map[string]map[string]Parser{
	"de": map[string]Parser{
		"noun":        NewRegexpParser(`^\pZ*(?P<word>\pL+)\pZ+(?P<gender>[mfnMFN])\pZ*\(genitive\pZ+(?P<genitive>\pL+)[^,]*,\pZ*(?:plural\pZ+(?P<plural>\pL+)|no\pZ+plural)`),
		"adjective":   NewRegexpParser(`^\pZ*(?P<word>\pL+)\pZ*\(\pZ*comparative\pZ+(?P<comparative>\pL+)[^,]*,\pZ*superlative\pZ+(?P<superlative>[^)]+)`),
		"preposition": NewRegexpParser(`^\pZ*(?P<word>\pL+)\pZ*\(\+\pZ*(?P<case>\pL+)`),
		"adverb":      NewRegexpParser(`^\pZ*(?P<word>\pL+)`),
		"conjunction": NewRegexpParser(`^\pZ*(?P<word>\pL+)`),
		"verb":        &FuncParser{f: ParseGermanVerb},
	},
}

type FuncParser struct {
	f func(doc *goquery.Document, lang, title string) (*word.Word, error)
}

func (p FuncParser) Parse(doc *goquery.Document, lang, title string) (*word.Word, error) {
	h := findHeadline(doc, lang, title)
	if h.Size() == 0 {
		return nil, fmt.Errorf("No headline found for part of speech: %s\n", title)
	}

	w, err := p.f(doc, lang, title)
	if err != nil {
		return nil, err
	}
	w.Defs = scrapeDefs(h)

	return w, nil
}

type RegexpParser struct {
	re    *regexp.Regexp
	names []string
}

func NewRegexpParser(re string) Parser {
	p := new(RegexpParser)
	p.re = regexp.MustCompile(re)
	p.names = p.re.SubexpNames()

	return p
}

func (rp RegexpParser) Parse(doc *goquery.Document, lang, title string) (*word.Word, error) {
	h := findHeadline(doc, lang, title)
	if h.Size() == 0 {
		return nil, fmt.Errorf("No headline found for part of speech: %s\n", title)
	}

	s := h.Parent().Next().Text()
	m := rp.re.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("Regex failed to match %s", s)
	}

	w := new(word.Word)
	w.Defs = scrapeDefs(h)
	w.Attr = make(map[string]string)

	names := rp.re.SubexpNames()
	for i, match := range m {
		if i == 0 {
			continue
		}

		switch names[i] {
		case "word":
			w.Word = match
		default:
			w.Attr[names[i]] = match
		}
	}

	return w, nil
}

// Helper functions for DOM traversal
func scrapeDefs(headline *goquery.Selection) []string {
	defs := make([]string, 0)

	ol := headline.Parent().NextFilteredUntil("ol", "hr").First()
	if len(ol.Nodes) != 1 {
		return defs
	}

	ol.ChildrenFiltered("li").Each(func(i int, s *goquery.Selection) {
		defs = append(defs, s.Text())
	})

	return defs
}

func findHeadline(doc *goquery.Document, lang, title string) *goquery.Selection {
	langHeading := doc.Find("#" + lang).Parent()
	sectionHeadings := langHeading.NextUntil("hr").Find(".mw-headline")

	h := sectionHeadings.FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == title
	})

	return h
}
