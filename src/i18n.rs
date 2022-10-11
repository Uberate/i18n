use std::collections::HashMap;
use std::sync::RwLock;

/// The struct I18n save all message(namespace, code and language). It provide some simple apis to
/// get message. The i18n is thread-safe. And I18n support disable write behavior by settings.
///
/// Note that, because of [I18n] like a const value storage, so all the value move is clone. So for
/// high-performance, the I18n info should build at application bootstrap age. At run time, the I18n
/// should be a read-only system(or tools) to used.
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
    pub fn message(&self, namespace: &String, code: &String, language: &String) -> Option<String> {
        let _ = self.rw_lock.read();
        if let Some(namespace_item) = self.messages.get(namespace) {
            if let Some(code_item) = namespace_item.get(code) {
                if let Some(message) = code_item.get(language) {
                    return Some(message.clone());
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