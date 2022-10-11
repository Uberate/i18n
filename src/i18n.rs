use std::collections::HashMap;
use std::sync::RwLock;

/// The struct I18n save all message(namespace, code and language). It provide some simple apis to
/// get message. The i18n is thread-safe. And I18n support disable write behavior by settings.
///
/// Note that, because of [I18n] like a const value storage, so all the value move is clone. So for
/// high-performance, the I18n info should build at application bootstrap age. At run time, the I18n
/// should be a read-only system(or tools) to used.
///
/// The disable_change is non reversible operation. For this, we think the i18n resource is a static
/// resource when [I18n] used in a application or system. But the [I18n] tools can ignore this rule.
/// The tools is the application which are used to generate or manage the I18n messages.
///
#[derive(Debug)]
pub struct I18n {
    // messages save the all message info. namespace-code-ln-message
    messages: HashMap<String, HashMap<String, HashMap<String, String>>>,

    default_language: String,

    enable_change: bool,

    rw_lock: RwLock<i8>,
}

impl I18n {
    /// The function message will return the message by specify namespace, code and languages.
    ///
    /// If specify message not found, the [message] will return [None], else return [Some<String>].
    /// And if specify message not found in specify language, will try to search message in
    /// default_language(if default_language not empty).
    ///
    /// ## Example
    /// ```rust
    /// let mut i = I18n::new("en".to_string(), true);
    /// i.register_message(&"test".to_string(),
    ///                    &"test".to_string(),
    ///                    &"test".to_string(),
    ///                    &"test".to_string());
    /// match i.message(&"test".to_string(), &"test".to_string(), &"test".to_string()) {
    ///      Some(t) => { assert_eq!(t, "test".to_string()); }
    ///      _ => { assert_eq!(1, 2); }
    /// }
    /// ```
    pub fn message(&self, namespace: &String, code: &String, language: &String) -> Option<&String> {
        let _ = self.rw_lock.read();
        if let Some(namespace_item) = self.messages.get(namespace) {
            if let Some(code_item) = namespace_item.get(code) {
                if let Some(message) = code_item.get(language) {
                    return Some(message);
                }
            }
        }

        if language != &self.default_language && self.default_language.len() != 0 {
            return self.message(namespace, code, language);
        }

        None
    }

    /// The function register_message_object like [register_message], is used to register a message
    /// info to [I18n].
    pub fn register_message_object(&mut self, object: MessageObject) {
        for x in object.message {
            self.register_message(&object.namespace, &object.code, &x.0, &x.1)
        }
    }

    /// The function register_message_objects will register all message object to [I18n].
    pub fn register_message_objects(&mut self, objects: Vec<MessageObject>) {
        for x in objects {
            self.register_message_object(x);
        }
    }

    /// The function register_message will register a message to [I18n], and if message is empty,
    /// will remove. And the [I18n] always keep less memory used, so it will prune the empty
    /// namespace and code value.
    ///
    /// If the [I18n] is not enable change, the register will do nothing.
    /// ## Example
    /// ```
    /// let mut i = I18n::new("en".to_string(), true);
    /// i.register_message(&"test".to_string(),
    ///                    &"test".to_string(),
    ///                    &"test".to_string(),
    ///                    &"test".to_string());
    /// match i.message(&"test".to_string(), &"test".to_string(), &"test".to_string()) {
    ///      Some(t) => { assert_eq!(t, "test".to_string()); }
    ///      _ => { assert_eq!(1, 2); }
    /// }
    /// ```
    ///
    /// You can use [MessageObject] directly to register. See function [register_message_object].
    pub fn register_message(&mut self,
                            namespace: &String, code: &String, language: &String, message: &String) {
        if !self.enable_change {
            return;
        }
        // write lock
        let _ = self.rw_lock.write().unwrap();

        // get namespace
        if !self.messages.contains_key(namespace) {
            self.messages.insert(namespace.clone(), HashMap::new());
        }
        // get code
        if let Some(namespace_item) = self.messages.get_mut(namespace) {
            if !namespace_item.contains_key(code) {
                namespace_item.insert(code.clone(), HashMap::new());
            }
            if let Some(code_item) = namespace_item.get_mut(code) {
                if message.len() == 0 {
                    // remove specify language value
                    code_item.remove(language);
                } else {
                    code_item.insert(language.clone(), message.clone());
                }

                // if code item is emtpy, remove the code value from the namespace
                if code_item.len() == 0 {
                    namespace_item.remove(code);
                }
            }

            // if namespace is emtpy, remove the namespace value
            if namespace_item.len() == 0 {
                self.messages.remove(namespace);
            }
        }
    }

    /// The disable_change is non reversible function, that means you can't change the changeable
    /// when [I18n] is disable changed.
    pub fn disable_change(&mut self) {
        self.enable_change = false
    }

    pub fn is_enable_change(&self) -> bool {
        return self.enable_change;
    }

    pub fn set_default_language(&mut self, default_language: String) {
        self.default_language = default_language
    }

    pub fn new(default_language: String, enable_change: bool, mut read_thread_count: i8) -> Self {
        // if read_thread_count <= 0, reset to default '10'
        if read_thread_count <= 0 {
            read_thread_count = 10
        }
        I18n {
            messages: HashMap::new(),
            default_language,
            enable_change,
            rw_lock: RwLock::new(read_thread_count),
        }
    }
}

/// The struct MessageObject is contain a message info for all languages.
pub struct MessageObject {
    namespace: String,
    code: String,
    message: HashMap<String, String>,
}

impl MessageObject {
    /// The function from_strings will build a Vec<MessageObject> from values. The values must has
    /// a header at first line. The first line always define the language info.
    ///
    /// ## Like this:
    /// ```
    /// ------------------------------------
    /// | namespace |  code  |  en  |  zh  |  <- Header at first line.
    /// ------------------------------------
    /// |   test    |  test  | test | test |  <- Value start at second line.
    /// |   test2   |  test  |      | test |
    /// |   test    |  test2 | test |      |  <- If specify language not set value,
    ///      ^           ^      ^                 used emtpy value.
    ///       \           \      \______________ The language info always start
    ///        \           \                      at third column.
    ///         \___________\___________________ The first and second column is
    ///                                           namespace and code.
    /// ```
    ///
    /// The function will try to read header, if header can't be parsed, return emtpy Array.
    ///
    /// For less memory used, it will ignore emtpy info.
    pub fn from_strings(values: &Vec<Vec<String>>) -> Vec<MessageObject> {
        let mut res = Vec::new();

        let mut find_header = false;
        let mut header = Vec::new();
        for x in values {
            // fist is header
            if !find_header {

                // If header.len() <= 2, it means that not define any language. Return emtpy
                // array directly.
                if x.len() <= 2 {
                    return res;
                }

                // To get all language value form second column.
                for index in 2..x.len() {
                    if let Some(ln) = x.get(index) {
                        header.push(ln.clone())
                    }
                }
                find_header = true;
                continue;
            }

            //--------------------------------------------------
            // not the header

            if x.len() <= 2 {
                // The message has no value, skip it.
                continue;
            }

            // get namespace(at index 0) and code(at index 1). If not found, skip this value.
            let namespace = match x.get(0) {
                Some(namespace) => namespace,
                _ => continue
            };
            let code = match x.get(1) {
                Some(code) => code,
                _ => continue
            };

            //--------------------------------------------------
            // Build a MessageObject.

            // Create a new MessageObject
            let mut mo = MessageObject {
                namespace: namespace.clone(),
                code: code.clone(),
                message: HashMap::new(),
            };

            // Set language message info.
            for index in 2..x.len() {
                if let Some(ln_message) = x.get(index) {

                    // Get language name from header.
                    if let Some(ln_name) = header.get(index - 2) {
                        mo.message.insert(ln_name.clone(), ln_message.clone());
                    }
                }
            }

            // if message has no message, skip this message.
            if mo.message.len() != 0 {
                // Add this message object.
                res.push(mo);
            }
        };

        res
    }


    pub fn message_objects_to_strings(message_objects: Vec<MessageObject>) -> Vec<Vec<String>> {
        let mut res = Vec::new();

        // To cache the header index.
        let mut header: HashMap<String, usize> = HashMap::new();
        let mut current_header_index: usize = 0;

        // 'header_string' will insert to res[0].
        let mut header_string: Vec<String> = Vec::from(["namespace".to_string(), "code".to_string()]);

        // For each all message_objects
        for x in message_objects {

            // Get namespace and code and push them to res.
            let namespace = &x.namespace;
            let code = &x.code;
            let mut temp_string: Vec<String> = Vec::new();
            temp_string.push(namespace.clone());
            temp_string.push(code.clone());

            // Get all language info.
            for (language, message) in x.message {

                // If header not contains the language, insert it.
                if !header.contains_key(&language) {
                    header.insert(language.clone(), current_header_index);
                    current_header_index = current_header_index + 1;
                    header_string.push(language.clone())
                }

                // Get the header index.
                if let Some(index) = header.get(&language) {
                    while temp_string.len() <= (index.clone() + 2) {
                        temp_string.push("".to_string())
                    }

                    // Replace specify language value.
                    temp_string.remove(index.clone() + 2);
                    temp_string.insert(index.clone() + 2, message);
                }
            }
            // Insert the value.
            res.push(temp_string);
        };

        // Complete all empty language for each line of message end.
        for x in res.iter_mut() {
            if x.len() < header_string.len() {
                x.push("".to_string());
            }
        }
        res.insert(0, header_string);
        res
    }
}

