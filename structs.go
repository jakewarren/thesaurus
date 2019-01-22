package thesaurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

// A CrossReference is a reference to a separate dictionary word's entry and how
// it relates to the current result (Type).
type CrossReference struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

type Derivative struct {
	Domains   []string `json:"domains,omitempty"`
	ID        string   `json:"id,omitempty"`
	Language  string   `json:"language,omitempty"`
	Regions   []string `json:"regions,omitempty"`
	Registers []string `json:"registers,omitempty"`
	Text      string   `json:"text,omitempty"`
}

func (d *Derivative) String() string {
	return d.Text
}

type Entry struct {
	Etymologies         []string             `json:"etymologies,omitempty"`
	GrammaticalFeatures []GrammaticalFeature `json:"grammaticalFeatures,omitempty"`
	HomographNumber     string               `json:"homographNumber,omitempty"`
	Notes               []interface{}        `json:"notes,omitempty"`
	Senses              []Sense              `json:"senses,omitempty"`
	VariantForms        []VariantForm        `json:"variantForms,omitempty"`
}

func (e Entry) RenderVariantForms() string {
	variants := make([]string, len(e.VariantForms))
	for i, variant := range e.VariantForms {
		variants[i] = variant.Text
	}
	return fmt.Sprintf("Alternatively: %s", strings.Join(variants, "', '"))
}

type Example struct {
	Registers []string `json:"registers,omitempty"`
	Text      string   `json:"text,omitempty" mapstructure:"text"`
}

// Render renders an example, optionally in a short form that omits registers
// and other information.
func (e *Example) Render() string {
	var buf bytes.Buffer

	if len(e.Registers) > 0 {
		regs := make([]string, len(e.Registers))
		for i, reg := range e.Registers {
			regs[i] = reg
		}
		fmt.Fprintf(&buf, "`%s` ", strings.Join(regs, "`, `"))
	}
	fmt.Fprintf(&buf, "'%s'", e.Text)
	return buf.String()
}

type GrammaticalFeature struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

type HTMLOutput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// A LexicalEntry is the grouping of senses and definitions that a word posesses
// for a given lexical category (noun, verb, etc.).
type LexicalEntry struct {
	DerivativeOf    []Derivative    `json:"derivativeOf,omitempty"`
	Derivatives     []Derivative    `json:"derivatives,omitempty"`
	Entries         []Entry         `json:"entries,omitempty"`
	LexicalCategory string          `json:"lexicalCategory,omitempty"`
	Pronunciations  []Pronunciation `json:"pronunciations,omitempty"`
}

// IsDerivative checks if a LexicalEntry is a derivative of another word.
func (l *LexicalEntry) IsDerivative() bool {
	return len(l.DerivativeOf) > 0
}

func (l *LexicalEntry) ShortLexicalCategory() string {
	switch strings.ToLower(l.LexicalCategory) {
	case "noun":
		return "n."
	case "adjective":
		return "adj."
	case "verb":
		return "v."
	default:
		return ""
	}
}

// RenderPronunciation returns a stringified pronuncation string of a provided notation.
func (l LexicalEntry) RenderPronunciation(notation string) string {
	data, err := json.Marshal(l)
	if err != nil {
		return ""
	}
	res := gjson.Get(string(data), fmt.Sprintf("pronunciations.#[phoneticNotation==\"%s\"].phoneticSpelling", notation))
	return fmt.Sprintf("/%s/", res.String())
}

func (l *LexicalEntry) RenderLexicalCategory() string {
	return strings.ToUpper(l.LexicalCategory)
}

type Pronunciation struct {
	AudioFile        string   `json:"audioFile,omitempty"`
	Dialects         []string `json:"dialects,omitempty"`
	PhoneticNotation string   `json:"phoneticNotation,omitempty"`
	PhoneticSpelling string   `json:"phoneticSpelling,omitempty"`
}

func (p *Pronunciation) String() string {
	return fmt.Sprintf("/%s/", p.PhoneticSpelling)
}

// A Result is a response returned by the Oxford Dictionaries API. A result
// should be synonymous with a word's definition in a dictionary, with all
// associated usages.
type Result struct {
	ID             string         `json:"id,omitempty"`
	Language       string         `json:"language,omitempty"`
	LexicalEntries []LexicalEntry `json:"lexicalEntries,omitempty"`
	Type           string         `json:"type,omitempty"`
	Word           string         `json:"word"`
}

// RenderTitle renders a Markdown title for an entry.
func (r Result) RenderTitle() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s\n", r.Word)
	// BUG(travis-g): words should be treated as []rune
	for i := 0; i < len(r.Word); i++ {
		fmt.Fprintf(&buf, "▬")
	}
	return buf.String()
}

// FilterCategory removes all entries of a Result that do not match a given
// lexical category (noun, verb, adjective, etc.). The Result object itself is
// not modified.
func (r *Result) FilterCategory(category string) []LexicalEntry {
	matches := make([]LexicalEntry, 0)
	for _, x := range r.LexicalEntries {
		if strings.ToLower(x.LexicalCategory) == strings.ToLower(category) {
			matches = append(matches, x)
		}
	}
	return matches
}

// Results is a wrapper for the full response object sent by the API.
type Results struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Results  []Result               `json:"results,omitempty"`
}

// SplitResultsByHomograph takes the result returned by the API and splits it
// into sets of definitions, ex. ace[1] and ace[2].
//
// TODO(tssde71): this needs to have a lot of work done on it
func (r *Results) SplitResultsByHomograph() {
	if len(r.Results) == 0 {
		return
	}
	master := r.Results[0]
	homographs := make([]int, 0)
	for _, lexicalEntry := range master.LexicalEntries {
		for _, entry := range lexicalEntry.Entries {
			number, err := strconv.Atoi(entry.HomographNumber)
			if err != nil {
				panic(fmt.Sprintf("error parsing homograph ID: %s", err))
			}
			homographs = append(homographs, number/100)
		}
	}
	// uniques only
	homographs = Ints(homographs)

	// make an array of results the same length as the number of homographs
	results := make([]Result, len(homographs))

	// create complete clones except for entries:
	for i, _ := range results {
		result := Result{
			ID:             master.ID,
			Language:       master.Language,
			LexicalEntries: make([]LexicalEntry, 0),
			Type:           master.Type,
			Word:           master.Word,
		}
		results[i] = result
	}

	fmt.Println(homographs, len(homographs), Min(homographs))
}

type Sense struct {
	Antonyms              []ThesaurusEntry `json:"antonyms,omitempty" mapstructure:"antonyms"`
	CrossReferenceMarkers []string         `json:"crossReferenceMarkers,omitempty"`
	CrossReferences       []CrossReference `json:"crossReferences,omitempty"`
	Definitions           []string         `json:"definitions,omitempty"`
	Domains               []string         `json:"domains,omitempty"`
	Examples              []Example        `json:"examples,omitempty" mapstructure:"examples"`
	ID                    string           `json:"id,omitempty"`
	Notes                 []interface{}    `json:"notes,omitempty"`
	Regions               []string         `json:"regions,omitempty"`
	Registers             []string         `json:"registers,omitempty"`
	ShortDefinitions      []string         `json:"short_definitions,omitempty"`
	Subsenses             []Sense          `json:"subsenses,omitempty" mapstructure:"subsenses"`
	Synonyms              []ThesaurusEntry `json:"synonyms,omitempty" mapstructure:"synonyms"`
	ThesaurusLinks        []ThesaurusLink  `json:"thesaurusLinks,omitempty"`
	VariantForms          []VariantForm    `json:"variantForms,omitempty"`
}

func (s *Sense) ListExamples() []string {
	examples := make([]string, len(s.Examples))
	for i, ex := range s.Examples {
		examples[i] = ex.Text
	}
	return examples
}

func (s *Sense) HasAntonyms() bool {
	for _, entry := range s.Antonyms {
		if entry.Text != "" {
			return true
		}
	}
	return false
}

func (s *Sense) HasCrossReferences() bool {
	return len(s.CrossReferences) != 0
}

// HasDefinition returns true if a sense has a definition, or false if not.l
func (s Sense) HasDefinition() bool {
	return len(s.Definitions) != 0
}

func (s *Sense) HasSynonyms() bool {
	for _, entry := range s.Synonyms {
		if entry.Text != "" {
			return true
		}
	}
	return false
}

func (s *Sense) RenderAntonyms() string {
	antonyms := make([]string, len(s.Antonyms))
	for i, ant := range s.Antonyms {
		antonyms[i] = ant.Text
	}
	bold := color.New(color.Bold, color.FgBlue)

	antonyms[0] = bold.Sprintf("%s", antonyms[0])

	return strings.Join(antonyms, ", ")
}

func (s *Sense) RenderCrossReferences() string {
	markers := strings.Join(s.CrossReferenceMarkers, "; ")
	return fmt.Sprintf("[%s]", markers)
}

// RenderDefitions renders what is displayed as the definition of a word. This
// function falls back to displaying cross-references in the event that another
// word's definition is being referenced.
func (s Sense) RenderDefinitions() string {
	var buf bytes.Buffer
	if s.HasDefinition() {
		fmt.Fprint(&buf, capitalizeFirst(strings.Join(s.Definitions, "; ")))
	} else if s.HasCrossReferences() {
		fmt.Fprint(&buf, s.RenderCrossReferences())
	}
	return buf.String()
}

// RenderDomains creates a tag list of Domain information
func (s *Sense) RenderDomains() string {
	if len(s.Domains) == 0 {
		return ""
	}
	return fmt.Sprintf("`%s`", strings.Join(s.Domains, "`, `"))
}

func (s Sense) RenderExamples() string {
	var buf bytes.Buffer
	for _, ex := range s.Examples {
		fmt.Fprintf(&buf, "%s\n", ex.Render())
	}
	return buf.String()
}

func (s Sense) RenderRegisters() string {
	if len(s.Registers) == 0 {
		return ""
	}
	return fmt.Sprintf("`%s`", strings.Join(s.Registers, "`, `"))
}

func (s *Sense) RenderSynonyms() string {
	synonyms := make([]string, len(s.Synonyms))
	for i, syn := range s.Synonyms {
		synonyms[i] = syn.Text
	}

	bold := color.New(color.Bold, color.FgBlue)

	synonyms[0] = bold.Sprintf("%s", synonyms[0])

	return strings.Join(synonyms, ", ")
}

func (s Sense) RenderTags(filters ...string) string {
	tags := s.Tags(filters...)
	if len(tags) == 0 {
		return ""
	}

	return color.GreenString(fmt.Sprintf("%s ", strings.Join(tags, ", ")))
}

func (s *Sense) Tags(filters ...string) []string {
	var tags []string
	tags = append(tags, s.Regions...)
	tags = append(tags, s.Domains...)
	tags = append(tags, s.Registers...)
	for _, filter := range filters {
		tags = Filter(tags, func(s string) bool {
			return s != filter
		})
	}
	return tags
}

// A ThesaurusEntry is either a synonym or antonym of a word.
type ThesaurusEntry struct {
	ID       string `json:"id,omitempty"`
	Language string `json:"language,omitempty"`
	Text     string `json:"text,omitempty" mapstructure:"text"`
}

type ThesaurusLink struct {
	EntryID string `json:"entry_id,omitempty"`
	SenseID string `json:"sense_id,omitempty"`
}

type VariantForm struct {
	Text string `json:"text,omitempty"`
}

func (v *VariantForm) String() string {
	return fmt.Sprintf("'%s'", v.Text)
}
