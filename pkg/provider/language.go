package provider

//go:generate go run ./tools/code_gen/language_gen.go ./tools/code_gen/language_gen.tpl ./language_key.gen.go
import "strings"

const (
	Custom   = "Custom"      // None value for the default value.
	ISO6391  = "ISO 639-1"   // The ISO 639-1, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
	ISO6392B = "ISO 639-2 B" // The ISO 639-2 B, is a type of ISO 639-2, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes
	ISO6392T = "ISO 639-2 T" // The ISO 639-2 T, is a type of ISO 639-2, the WIKI: https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes
	ISO6393  = "ISO 639-3"
)

// NewLanguageKey return a *LanguageKey.
func NewLanguageKey() *LanguageKey {
	return &LanguageKey{}
}

// LanguageKey represent a language value in AbsI18n system. It can convert to different Standard.
type LanguageKey struct {
	DefaultStandard string            `json:"default_standard" yaml:"default_standard" mapstructure:"default_standard"`
	Keys            map[string]string `json:"Keys" yaml:"Keys" mapstructure:"Keys"`
}

// Upper return an upper language key value by Standard. If current LanguageKey doesn't have this Standard value, try to
// find value from LanguageKey.DefaultStandard. If default Standard not has value too, return empty value.
func (lk *LanguageKey) Upper(standard string) string {
	return strings.ToUpper(lk.Key(standard))
}

// Lower return a lower language key value by Standard. If current LanguageKey doesn't have this Standard value, try to
// find value from LanguageKey.DefaultStandard. If default Standard not has value too, return empty value.
func (lk *LanguageKey) Lower(standard string) string {
	return strings.ToLower(lk.Key(standard))
}

// Key return the language key value by Standard. If current LanguageKey doesn't have this Standard value, try to
// find value from LanguageKey.DefaultStandard. If default Standard not has value too, return empty value.
func (lk *LanguageKey) Key(standard string) string {
	if value, ok := lk.Keys[standard]; ok {
		return value
	}

	// If not found specify Standard and default Standard not nil, try to return default Standard key.
	// If input is default, return empty string directly.
	if len(lk.DefaultStandard) != 0 && lk.DefaultStandard != standard {
		return lk.Key(lk.DefaultStandard)
	}

	return ""
}

// SetDefaultStandard will update the default Standard.
func (lk *LanguageKey) SetDefaultStandard(defaultStandard string) *LanguageKey {
	lk.DefaultStandard = defaultStandard
	return lk
}

// Push will push a Standard value to current LanguageKey. And if LanguageKey.DefaultStandard is emtpy, set the value
// to this Standard. The value can be changed in anywhere.
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
)

func GetLanguageKey(keyName string) *LanguageKey {
	if value, ok := Mapper[keyName]; ok {
		return value
	}
	return NoneLn
}
