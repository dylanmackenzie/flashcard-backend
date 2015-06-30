package main

import (
	"encoding/json"
	"log"
	"net/http"

	"dylanmackenzie.com/flashcard/word"

	"github.com/gorilla/mux"
)

type VocabHandler struct {
	scrapers map[string]Scraper
}

func (v *VocabHandler) GetScraper(from, to string) (Scraper, bool) {
	sc, ok := v.scrapers[from+">"+to]
	return sc, ok
}

func (v *VocabHandler) SetScraper(from, to string, s Scraper) {
	v.scrapers[from+">"+to] = s
}

func (v *VocabHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	vars := mux.Vars(r)
	from := vars["from"]
	to := vars["to"]

	w.Header().Set("Access-Control-Allow-Origin", "*")

	sc, ok := v.GetScraper(from, to)
	if !ok {
		w.Write([]byte("No lang found"))
		return
	}

	query := vars["word"]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("No word in query"))
		return
	}

	var words []*word.Word
	var single *word.Word
	var err error
	pos := q.Get("pos")
	if pos == "" {
		words, err = sc.ScrapeWordAllPos(to, query)
	} else {
		single, err = sc.ScrapeWord(to, query, pos)
		words = []*word.Word{single}
	}

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Word not found"))
		return
	}

	json, err := json.Marshal(words)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Invalid definition returned from lookup"))
		return
	}

	w.Write(json)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	vocab := &VocabHandler{
		scrapers: map[string]Scraper{},
	}

	vocab.SetScraper("en", "de", Scraper{root: "http://en.wiktionary.org/wiki/"})

	r := mux.NewRouter()
	r.Handle("/api-v1/card/{from}/{to}/{word}", vocab)
	http.Handle("/", r)
	log.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
