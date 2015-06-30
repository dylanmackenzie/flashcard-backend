package word

type Word struct {
	Word string   `json:"word"`
	Pos  string   `json:"pos"`  // Part of speech
	Defs []string `json:"defs"` // Definitions

	// Miscellaneous stuff
	Attr map[string]string `json:"attr"`
}
