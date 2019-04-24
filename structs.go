package thesaurus

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Results struct {
	ID       string   `json:"id"`
	Metadata Metadata `json:"metadata"`
	Results  []Result `json:"results"`
	Word     string   `json:"word"`
}

type Metadata struct {
	Operation string `json:"operation"`
	Provider  string `json:"provider"`
	Schema    string `json:"schema"`
}

type Result struct {
	ID             string         `json:"id"`
	Language       string         `json:"language"`
	LexicalEntries []LexicalEntry `json:"lexicalEntries"`
	Type           string         `json:"type"`
	Word           string         `json:"word"`
}

type LexicalEntry struct {
	Entries         []Entry         `json:"entries"`
	Language        string          `json:"language"`
	LexicalCategory LexicalCategory `json:"lexicalCategory"`
	Text            string          `json:"text"`
}

func (entry *LexicalEntry) RenderLexicalCategory() string {
	return strings.ToUpper(entry.LexicalCategory.Text)
}

type Entry struct {
	Senses []Sense `json:"senses"`
}

type Sense struct {
	Antonyms  []Onym            `json:"antonyms"`
	Examples  []Example         `json:"examples"`
	ID        string            `json:"id"`
	Registers []LexicalCategory `json:"registers"`
	Subsenses []Subsense        `json:"subsenses"`
	Synonyms  []Onym            `json:"synonyms"`
}

func (s Sense) RenderExamples() string {
	var buf bytes.Buffer
	for _, ex := range s.Examples {
		fmt.Fprintf(&buf, "%s\n", ex.Render())
	}
	return buf.String()
}

func (s Sense) HasSynonyms() bool {
	return len(s.Synonyms) > 0
}

func (s Sense) HasAntonyms() bool {
	return len(s.Antonyms) > 0
}

func (s Sense) RenderTags() string {
	tags := make([]string, 0)

	for _, t := range s.Registers {
		tags = append(tags, t.Text)
	}

	if len(tags) == 0 {
		return ""
	}
	return color.GreenString(fmt.Sprintf("%s ", strings.Join(tags, ", ")))
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

func (s *Sense) RenderAntonyms() string {
	antonyms := make([]string, len(s.Antonyms))
	for i, syn := range s.Antonyms {
		antonyms[i] = syn.Text
	}

	bold := color.New(color.Bold, color.FgBlue)

	antonyms[0] = bold.Sprintf("%s", antonyms[0])

	return strings.Join(antonyms, ", ")
}

type Onym struct {
	ID       *string  `json:"id,omitempty"`
	Language Language `json:"language"`
	Text     string   `json:"text"`
}

type Example struct {
	Text string `json:"text"`
}

func (ex Example) Render() string {
	return ex.Text
}

type LexicalCategory struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Subsense struct {
	ID        string            `json:"id"`
	Synonyms  []Onym            `json:"synonyms"`
	Registers []LexicalCategory `json:"registers"`
	Regions   []LexicalCategory `json:"regions"`
}

func (s Subsense) HasSynonyms() bool {
	return len(s.Synonyms) > 0
}

func (s Subsense) RenderTags() string {
	tags := make([]string, 0)

	for _, t := range s.Regions {
		tags = append(tags, cleanTag(t.Text))
	}

	for _, t := range s.Registers {
		tags = append(tags, cleanTag(t.Text))
	}

	if len(tags) == 0 {
		return ""
	}
	return color.GreenString(fmt.Sprintf("%s ", strings.Join(tags, ", ")))
}

func cleanTag(s string) (output string) {
	output = s
	output = strings.Replace(output, "_", " ", -1)
	return
}

func (s *Subsense) RenderSynonyms() string {
	synonyms := make([]string, len(s.Synonyms))
	for i, syn := range s.Synonyms {
		synonyms[i] = syn.Text
	}

	bold := color.New(color.Bold, color.FgBlue)

	synonyms[0] = bold.Sprintf("%s", synonyms[0])

	return strings.Join(synonyms, ", ")
}

type Language string
