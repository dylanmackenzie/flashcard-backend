package main

import (
	"reflect"
	"testing"

	"dylanmackenzie.com/flashcard/word"

	"github.com/tonnerre/golang-pretty"
)

var tests = []word.Word{
	word.Word{
		Word: "rennen",
		Pos:  "verb",
		Attr: map[string]string{
			"presentParticiple": "rennend",
			"pastParticiple":    "gerannt",
			"auxiliary":         "haben or sein",
			"imperative":        "renn",

			"present":      "renne,rennst,rennt,rennen,rennt,rennen",
			"preterite":    "rannte,ranntest,rannte,rannten,ranntet,rannten",
			"subjunctive1": "renne,rennest,renne,rennen,rennet,rennen",
			"subjunctive2": "rennte,renntest,rennte,rennten,renntet,rennten",
		},
	},
	word.Word{
		Word: "gehen",
		Pos:  "verb",
		Attr: map[string]string{
			"presentParticiple": "gehend",
			"pastParticiple":    "gegangen",
			"auxiliary":         "sein",
			"imperative":        "geh",

			"present":      "gehe,gehst,geht,gehen,geht,gehen",
			"preterite":    "ging,gingst,ging,gingen,gingt,gingen",
			"subjunctive1": "gehe,gehest,gehe,gehen,gehet,gehen",
			"subjunctive2": "ginge,gingest,ginge,gingen,ginget,gingen",
		},
	},
	word.Word{
		Word: "fallen",
		Pos:  "verb",
		Attr: map[string]string{
			"presentParticiple": "fallend",
			"pastParticiple":    "gefallen",
			"auxiliary":         "sein",
			"imperative":        "fall",

			"present":      "falle,fällst,fällt,fallen,fallt,fallen",
			"preterite":    "fiel,fielst,fiel,fielen,fielt,fielen",
			"subjunctive1": "falle,fallest,falle,fallen,fallet,fallen",
			"subjunctive2": "fiele,fielest,fiele,fielen,fielet,fielen",
		},
	},
	word.Word{
		Word: "Gott",
		Pos:  "noun",
		Attr: map[string]string{
			"gender":   "m",
			"plural":   "Götter",
			"genitive": "Gottes",
		},
	},
	word.Word{
		Word: "stark",
		Pos:  "adjective",
		Attr: map[string]string{
			"comparative": "stärker",
			"superlative": "am stärksten",
		},
	},
	word.Word{
		Word: "durch",
		Pos:  "preposition",
		Attr: map[string]string{
			"case": "accusative",
		},
	},
}

func TestWordFromWiki(t *testing.T) {
	t.Parallel()
	sc := &Scraper{root: "http://en.wiktionary.org/wiki/"}
	for _, expected := range tests {
		w, err := sc.ScrapeWord("de", expected.Word, expected.Pos)
		if err != nil {
			t.Errorf("Error parsing %s: %s\n", sc.root+expected.Word, err)
			continue
		}

		w.Defs = nil
		if !reflect.DeepEqual(*w, expected) {
			t.Errorf("Not Equal %s: %v", sc.root+expected.Word, pretty.Diff(*w, expected))
		}
	}
}
