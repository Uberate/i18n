package i18n

func NewI18n(standard string) *I18n {
	return &I18n{
		standard: standard,
		values:   &Namespace{},
	}
}

type I18n struct {
	values *Namespace

	defaultLanguage LanguageKey
	standard        string
}

type Pusher func(ln LanguageKey, messageValue string)
type IPusher func(ln LanguageKey, messageValue string) IPusher
type PusherByString func(ln string, messageValue string)
type IPusherByString func(ln, messageValue string) IPusherByString

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

// PushMessage will push a message to I18n instance.
// i.PushMessage(EnglishLn, "test", "namespace", "code") will create a message with two scopes 'namespace' and 'code'.
// i.PushMessage(EnglishLn, "test", "code") will create a message with one scope 'code'.
//
// In this example, the 'namespace' and 'code' is the scope, the Message must have one scope at least. But about max
// count, the I18n not limit it, but I18n suggest the count of scope should less than 4. If you really need more scope,
// you can write more scopes. But you should think, why you need so much scopes.
//
// PushMessage will push the message, and drop some info of LanguageKey. The I18n think, in one system, only one
// standard should be used. If you have more than one standard, you should use two instance.
//
// If specify message already haven value, the new message will cover it directly. Specify if the input message value is
// emtpy, the message will be deleted. See Message.PushMessage.
//
// Note that the PushMessage not thread-safe.
func (i *I18n) PushMessage(ln LanguageKey, messageValue string, scopes ...string) {
	i.PushMessageByString(ln.Lower(i.standard), messageValue, scopes...)
}

// PushMessageByString like PushMessage, but it receives the string as language key.
func (i *I18n) PushMessageByString(ln string, message string, scopes ...string) {
	i.values.PushMessage(ln, message, scopes...)
}

// SetStandard will update inner standard, but for message which already in, the ln will not be change. It only effects
// new message info. And the I18n.Message maybe return empty with false.
func (i *I18n) SetStandard(standard string) {
	i.standard = standard
}

func (i *I18n) Standard() string {
	return i.standard
}

// Message return the message of specify language and scopes. If value not found, return empty and false.
//
// If the I18n standard changed, the value maybe not found.
func (i *I18n) Message(ln LanguageKey, scopes ...string) (string, bool) {
	return i.MessageByString(ln.Lower(i.standard), scopes...)
}

// MessageByString like Message, but it receives string as language key.
func (i *I18n) MessageByString(ln string, scopes ...string) (string, bool) {
	return i.values.Message(ln, scopes...)
}

// Pusher help to quick build I18n message. It returns a func to add different language message to specify scopes.
func (i *I18n) Pusher(scopes ...string) Pusher {
	return func(ln LanguageKey, messageValue string) {
		i.PushMessage(ln, messageValue, scopes...)
	}
}

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

func (i *I18n) MessageBuilder(scopes ...string) *messageBuilder {
	return &messageBuilder{
		push: i.Pusher(scopes...),
	}
}

func (i *I18n) MessageStringBuilder(scopes ...string) *messageStringBuilder {
	return &messageStringBuilder{
		push: i.PusherByString(scopes...),
	}
}

func (i *I18n) WalkRecord(f func(languageValue, messageValue string, flags ...string)) {
	i.values.WalkRecord(f)
}

func (i *I18n) WalkMessage(f func(message map[string]string, flags ...string)) {
	i.values.WalkMessage(f)
}

// ---------------------------------------------------------------------------------------------------------------------

const scopeHeaderPrefix = "_"

type Namespace struct {
	// MessageSave the
	Children map[string]*Namespace
	Messages map[string]*Message
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
	for scope, message := range namespace.Messages {
		flags := append(parentFlags, scope)
		f(message.message, flags...)
	}

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

	// not value, return empty value and false
	if len(levelCodes) == 0 {
		return "", false
	}
	currentLevelCode := levelCodes[0]

	// search value in current scope
	if len(levelCodes) == 1 {
		if value, ok := namespace.Messages[currentLevelCode]; ok {
			return value.Message(ln)
		}
	}

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
		// no found the scope
		return
	}

	// Only has one code, save it as Message
	if len(levelCodes) == 1 {
		if namespace.Messages[levelCodes[0]] == nil {
			namespace.Messages[levelCodes[0]] = &Message{
				message: map[string]string{},
			}
		}

		namespace.Messages[levelCodes[0]].PushMessage(ln, messageValue)
		return
	}

	// Try to push to children, if not contain this flag, push a new one.
	if _, ok := namespace.Children[levelCodes[0]]; !ok {
		namespace.Children[levelCodes[0]] = &Namespace{
			Children: map[string]*Namespace{},
			Messages: map[string]*Message{},
		}
	}

	namespace.Children[levelCodes[0]].PushMessage(ln, messageValue, levelCodes[1:]...)
}

type Message struct {
	message map[string]string
}

func (m *Message) Message(ln string) (string, bool) {
	value, ok := m.message[ln]
	return value, ok
}

func (m *Message) PushMessage(ln, messageValue string) {
	if m.message == nil {
		m.message = map[string]string{}
	}

	if len(messageValue) == 0 {
		delete(m.message, ln)
	} else {
		m.message[ln] = messageValue
	}
}
