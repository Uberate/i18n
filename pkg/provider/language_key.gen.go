// GENERATOR, DO NOT EDIT!
package provider

var (
	//--------------------------------------------------
	// language key define

	ChineseLn  = *(NewLanguageKey().Push(Custom, "chinese").Push(ISO6391, "zh").Push(ISO6392B, "chi").Push(ISO6392T, "zho"))
	EnglishLn  = *(NewLanguageKey().Push(Custom, "english").Push(ISO6391, "en").Push(ISO6392B, "en").Push(ISO6392T, "en"))
	JapaneseLn = *(NewLanguageKey().Push(Custom, "japanese").Push(ISO6391, "ja").Push(ISO6392B, "jap").Push(ISO6392T, "jap"))
	//--------------------------------------------------

	//--------------------------------------------------
	// mapper of language

	Mapper = map[string]*LanguageKey{
		ChineseLn.Lower(Custom):  &ChineseLn,
		EnglishLn.Lower(Custom):  &EnglishLn,
		JapaneseLn.Lower(Custom): &JapaneseLn,
	}
	//--------------------------------------------------
)
