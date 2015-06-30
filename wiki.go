package main

import (
	"fmt"
	"strings"

	"dylanmackenzie.com/flashcard/parsers"
	"dylanmackenzie.com/flashcard/word"

	"github.com/PuerkitoBio/goquery"
)

var langCodes = map[string]string{
	"de": "German",
	"en": "English",
}

type Scraper struct {
	root string
}

func (sc Scraper) ScrapeWordAllPos(lang, vocab string) ([]*word.Word, error) {
	words := make([]*word.Word, 0)
	doc, err := sc.getPage(vocab)
	if err != nil {
		return nil, err
	}

	for pos, _ := range parsers.List[lang] {
		w, err := scrapeWord(doc, lang, vocab, pos)
		if err != nil {
			continue
		}

		words = append(words, w)
	}

	if len(words) == 0 {
		return nil, fmt.Errorf("No parts of speech found for %s", vocab)
	}

	return words, nil
}

func (sc Scraper) ScrapeWord(lang, vocab, pos string) (*word.Word, error) {
	doc, err := sc.getPage(vocab)
	if err != nil {
		return nil, err
	}

	return scrapeWord(doc, lang, vocab, pos)
}

func (sc Scraper) getPage(vocab string) (*goquery.Document, error) {
	url := sc.root + vocab
	return goquery.NewDocument(url)
}

func scrapeWord(doc *goquery.Document, lang, vocab, pos string) (*word.Word, error) {
	fullLang, _ := langCodes[lang]
	titlePos := strings.Title(pos)

	p, ok := parsers.List[lang][pos]
	if !ok {
		return nil, fmt.Errorf("Parser not found for %s:%s", lang, pos)
	}

	w, err := p.Parse(doc, fullLang, titlePos)
	if err != nil {
		return nil, err
	}
	w.Pos = pos

	return w, nil
}
