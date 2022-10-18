package i18n

func NewI18n(standard string) *I18n {
	if len(standard) == 0 {
		standard = Custom
	}
	return &I18n{
		Standard: standard,
		Values:   NewNamespace(),
	}
}

// I18n is core struct of i18n project. It save all MessageValue info.
//
// I18n will drop the language detail info of a LanguageKey. It will save the language string value by I18n.Standard.
// You can change Standard, but the old MessageValue will not be updated. If you want to use different Standard, please
// use different I18n instance to do it.
//
// I18n is not thread-safe, the I18n MessageValue should build at application bootstrap age. After bootstrap, the I18n
// info should not change(the tools of I18n except).
type I18n struct {
	Values          *Namespace  `yaml:"values" json:"values"`
	DefaultLanguage LanguageKey `yaml:"default_language" json:"default_language"`
	Standard        string      `yaml:"standard" json:"standard"`
}

//--------------------------------------------------
// Helper function define start.

type Pusher func(ln LanguageKey, messageValue string)
type IPusher func(ln LanguageKey, messageValue string) IPusher
type PusherByString func(ln string, messageValue string)
type IPusherByString func(ln, messageValue string) IPusherByString

// Helper function define end.
//--------------------------------------------------

//--------------------------------------------------
// Helper struct define start.

type messageBuilder struct {
	push Pusher
}

func (mb *messageBuilder) Push(key LanguageKey, message string) *messageBuilder {
	mb.push(key, message)
	return mb
}

type messageStringBuilder struct {
	push PusherByString
}

func (mb *messageStringBuilder) Push(key, message string) *messageStringBuilder {
	mb.push(key, message)
	return mb
}

// Helper struct define end.
//--------------------------------------------------

// PushMessage will push a MessageValue to I18n instance.
// i.PushMessage(EnglishLn, "test", "namespace", "code") will create a MessageValue with two scopes 'namespace' and 'code'.
// i.PushMessage(EnglishLn, "test", "code") will create a MessageValue with one scope 'code'.
//
// In this example, the 'namespace' and 'code' is the scope, the Message must have one scope at least. But about max
// count, the I18n not limit it, but I18n suggest the count of scope should less than 4. If you really need more scope,
// you can write more scopes. But you should think, why you need so much scopes.
//
// PushMessage will push the MessageValue, and drop some info of LanguageKey. The I18n think, in one system, only one
// Standard should be used. If you have more than one Standard, you should use two instance.
//
// If specify MessageValue already haven value, the new MessageValue will cover it directly. Specify if the input
// MessageValue value is emtpy, the MessageValue will be deleted. See Message.PushMessage.
//
// Note that the PushMessage not thread-safe.
func (i *I18n) PushMessage(ln LanguageKey, messageValue string, scopes ...string) {
	i.PushMessageByString(ln.Lower(i.Standard), messageValue, scopes...)
}

// PushMessageByString like PushMessage, but it receives the string as language key.
func (i *I18n) PushMessageByString(ln string, message string, scopes ...string) {
	i.Values.PushMessage(ln, message, scopes...)
}

// Message return the MessageValue of specify language and scopes. If value not found, return empty and false.
//
// If the I18n Standard changed, the value maybe not found.
func (i *I18n) Message(ln LanguageKey, scopes ...string) (string, bool) {
	return i.MessageByString(ln.Lower(i.Standard), scopes...)
}

// MessageByString like Message, but it receives string as language key.
func (i *I18n) MessageByString(ln string, scopes ...string) (string, bool) {
	return i.Values.Message(ln, scopes...)
}

// Pusher help to quick build I18n MessageValue. It returns a func to add different language MessageValue to specify
// scopes.
func (i *I18n) Pusher(scopes ...string) Pusher {
	return func(ln LanguageKey, messageValue string) {
		i.PushMessage(ln, messageValue, scopes...)
	}
}

// IPusher like Pusher, but it returns IPusher function. The IPusher can invoke more times like it:
// iPusher := i18n.IPusher("test")
// iPusher(EnglishLn, "test")("ChineseLn", "测试")
func (i *I18n) IPusher(scopes ...string) IPusher {
	return func(ln LanguageKey, messageValue string) IPusher {
		i.PushMessage(ln, messageValue, scopes...)
		return i.IPusher(scopes...)
	}
}

// PusherByString like Pusher, but the PusherByString receives the string as language key.
func (i *I18n) PusherByString(scopes ...string) PusherByString {
	return func(ln string, messageValue string) {
		i.PushMessageByString(ln, messageValue, scopes...)
	}
}

// IPusherByString like IPusher, but the IPusherByString receives the string as language key.
func (i *I18n) IPusherByString(scopes ...string) IPusherByString {
	return func(ln, messageValue string) IPusherByString {
		i.PushMessageByString(ln, messageValue)
		return i.IPusherByString(scopes...)
	}
}

// MessageBuilder return a builder for specify scopes.
func (i *I18n) MessageBuilder(scopes ...string) *messageBuilder {
	return &messageBuilder{
		push: i.Pusher(scopes...),
	}
}

// MessageStringBuilder return a builder for specify scopes.
func (i *I18n) MessageStringBuilder(scopes ...string) *messageStringBuilder {
	return &messageStringBuilder{
		push: i.PusherByString(scopes...),
	}
}

// WalkRecord will for-each all MessageValue language-value.
func (i *I18n) WalkRecord(f func(languageValue, messageValue string, flags ...string)) {
	i.Values.WalkRecord(f)
}

// WalkMessage will for-each all MessageValue value.
func (i *I18n) WalkMessage(f func(message map[string]string, flags ...string)) {
	i.Values.WalkMessage(f)
}

// IsMessageEquals return true when i.message == b.message. And if a == b == nil, return true
// Else if (i == b && b != nil) || (a != nil || b == nil) return false.
func (i *I18n) IsMessageEquals(b *I18n) bool {
	if i == b {
		return true
	}
	if i == nil || b == nil {
		return false
	}
	return i.IsSub(b) && b.IsSub(i)
}

// IsSub return true when current instance all infos can find from b. And if a == b == nil, return true
// Else if (i == b && b != nil) || (a != nil || b == nil) return false.
func (i *I18n) IsSub(b *I18n) bool {
	if i == b {
		return true
	}
	if i == nil || b == nil {
		return false
	}
	res := true
	i.WalkRecord(func(languageValue, messageValue string, flags ...string) {
		if res {
			if currentValue, ok := b.MessageByString(languageValue, flags...); !ok || currentValue != messageValue {
				res = false
			}
		}
	})
	return res
}

// ---------------------------------------------------------------------------------------------------------------------

const scopeHeaderPrefix = "_"

func NewNamespace() *Namespace {
	return &Namespace{
		Children: map[string]*Namespace{},
		Messages: NewMessage(),
	}
}

type Namespace struct {
	// MessageSave the
	Children map[string]*Namespace `yaml:"children" json:"children"`
	Messages *Message              `yaml:"messages" json:"messages"`
}

func (namespace *Namespace) WalkRecord(f func(ln, messageValue string, flags ...string)) {
	namespace.WalkMessage(func(message map[string]string, flags ...string) {
		for ln, message := range message {
			f(ln, message, flags...)
		}
	})
}

func (namespace *Namespace) WalkMessage(f func(message map[string]string, flags ...string)) {
	namespace.walkMessage(f)
}

func (namespace *Namespace) walkMessage(f func(message map[string]string, flags ...string), parentFlags ...string) {
	if parentFlags == nil {
		parentFlags = []string{}
	}
	f(namespace.Messages.MessageValue, parentFlags...)

	for scope, child := range namespace.Children {
		newScope := append(parentFlags, scope)
		child.walkMessage(f, newScope...)
	}
}

// Pusher is a specify iterator implements. It used to register value.
func (namespace *Namespace) Pusher(scopes ...string) func(string, string) {
	return func(ln, messageValue string) {
		namespace.PushMessage(ln, messageValue, scopes...)
	}
}

func (namespace *Namespace) Message(ln string, levelCodes ...string) (string, bool) {

	if len(levelCodes) == 0 {
		return namespace.Messages.Message(ln)
	}
	currentLevelCode := levelCodes[0]

	if value, ok := namespace.Children[currentLevelCode]; ok {

		return value.Message(ln, levelCodes[1:]...)
	}

	// not found, return emtpy value with false.
	return "", false
}

func (namespace *Namespace) PushMessage(ln, messageValue string, levelCodes ...string) {
	if namespace.Children == nil {
		namespace.Children = map[string]*Namespace{}
	}

	// not value, return empty value and false
	if len(levelCodes) == 0 {
		namespace.Messages.PushMessage(ln, messageValue)
		return
	}

	// Try to push to children, if not contain this flag, push a new one.
	if _, ok := namespace.Children[levelCodes[0]]; !ok {
		namespace.Children[levelCodes[0]] = NewNamespace()
	}

	namespace.Children[levelCodes[0]].PushMessage(ln, messageValue, levelCodes[1:]...)
}

type Message struct {
	MessageValue map[string]string `json:"message_value" yaml:"message_value"`
}

func NewMessage() *Message {
	return &Message{
		MessageValue: map[string]string{},
	}
}

func (m *Message) Message(ln string) (string, bool) {
	value, ok := m.MessageValue[ln]
	return value, ok
}

func (m *Message) PushMessage(ln, messageValue string) {
	if m.MessageValue == nil {
		m.MessageValue = map[string]string{}
	}

	if len(messageValue) == 0 {
		delete(m.MessageValue, ln)
	} else {
		m.MessageValue[ln] = messageValue
	}
}
