package i18n

import "strings"

const (
	CustomStandard = "Custom standard" // None value or unknown standard
	ISO6391        = "ISO 639-1"       // The ISO 639-1, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
	ISO6392B       = "ISO 639-2 B"     // The ISO 639-2 B, is a type of ISO 639-2, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes
	ISO6392T       = "ISO 639-2 T"     // The ISO 639-2 T, is a type of ISO 639-2, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes
	ISO6393        = "ISO 639-3"
)

// NewLanguageKey return a *LanguageKey.
func NewLanguageKey() *LanguageKey {
	return &LanguageKey{}
}

// LanguageKey represent a language value in I18n system. It can convert to different standard.
type LanguageKey struct {
	DefaultStandard string            `json:"default_standard" yaml:"default_standard" mapstructure:"default_standard"`
	Keys            map[string]string `json:"Keys" yaml:"Keys" mapstructure:"Keys"`
}

// Upper return an upper language key value by standard. If current LanguageKey doesn't have this standard value, try to
// find value from LanguageKey.DefaultStandard. If default standard not has value too, return empty value.
func (lk *LanguageKey) Upper(standard string) string {
	return strings.ToUpper(lk.Key(standard))
}

// Lower return a lower language key value by standard. If current LanguageKey doesn't have this standard value, try to
// find value from LanguageKey.DefaultStandard. If default standard not has value too, return empty value.
func (lk *LanguageKey) Lower(standard string) string {
	return strings.ToLower(lk.Key(standard))
}

// Key return the language key value by standard. If current LanguageKey doesn't have this standard value, try to
// find value from LanguageKey.DefaultStandard. If default standard not has value too, return empty value.
func (lk *LanguageKey) Key(standard string) string {
	if value, ok := lk.Keys[standard]; ok {
		return value
	}

	// If not found specify standard and default standard not nil, try to return default standard key.
	// If input is default, return empty string directly.
	if len(lk.DefaultStandard) != 0 && lk.DefaultStandard != standard {
		return lk.Key(lk.DefaultStandard)
	}

	return ""
}

// SetDefaultStandard will update the default standard.
func (lk *LanguageKey) SetDefaultStandard(defaultStandard string) *LanguageKey {
	lk.DefaultStandard = defaultStandard
	return lk
}

// Push will push a standard value to current LanguageKey. And if LanguageKey.DefaultStandard is emtpy, set the value
// to this standard. The value can be changed in anywhere.
func (lk *LanguageKey) Push(standard, value string) *LanguageKey {

	if len(lk.DefaultStandard) == 0 {

	}

	if lk.Keys == nil {
		lk.Keys = map[string]string{}
	}

	lk.Keys[standard] = value
	return lk
}

var (
	// NoneLn is a specify language, it means no language.
	NoneLn = NewLanguageKey().SetDefaultStandard(ISO6391).Push(ISO6391, "none")

	//--------------------------------------------------
	// todo: should generate it by generator.

	ChineseLn  = NewLanguageKey().Push(ISO6391, "zh").Push(ISO6392B, "chi").Push(ISO6392T, "zho")
	EnglishLn  = NewLanguageKey().Push(ISO6391, "en").Push(ISO6392B, "en").Push(ISO6392T, "en")
	JapaneseLn = NewLanguageKey().Push(ISO6391, "ja").Push(ISO6392B, "jpn").Push(ISO6392T, "jpn")

	Mapper = map[string]*LanguageKey{
		"chinese":  ChineseLn,
		"english":  EnglishLn,
		"japanese": JapaneseLn,
	}
	//--------------------------------------------------
)

func GetLanguageKey(keyName string) *LanguageKey {
	if value, ok := Mapper[keyName]; ok {
		return value
	}
	return NoneLn
}