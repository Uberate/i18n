// GENERATOR, DO NOT EDIT!
package i18n

var (
    //--------------------------------------------------
    // language key define

{{ range $languageName, $item := .}}
    {{$languageName}} = *(NewLanguageKey()
    {{- if $item.Custom}}.Push(Custom, "{{$item.Custom}}"){{- end}}
    {{- if $item.ISO6391}}.Push(ISO6391, "{{$item.ISO6391}}"){{- end}}
    {{- if $item.ISO6392B}}.Push(ISO6392B, "{{$item.ISO6392B}}"){{- end}}
    {{- if $item.ISO6392T}}.Push(ISO6392T, "{{$item.ISO6392T}}"){{- end -}}
    )
{{- end}}
    //--------------------------------------------------

    //--------------------------------------------------
    // mapper of language

    Mapper = map[string]*LanguageKey{
        {{- range $languageName, $item := .}}
            {{$languageName}}.Lower(Custom): &{{$languageName}},
        {{- end }}
    }
    //--------------------------------------------------
)