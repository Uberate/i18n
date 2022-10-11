use crate::i18n::{I18n, MessageObject};

#[test]
fn test_i18n() {
    let mut i = I18n::new("en".to_string(), true, 2);
    i.register_message(&"test".to_string(),
                       &"test".to_string(),
                       &"test".to_string(),
                       &"test".to_string());
    match i.message(&"test".to_string(), &"test".to_string(), &"test".to_string()) {
        Some(t) => { assert_eq!(t, &"test".to_string()); }
        _ => { assert_eq!(1, 2); }
    }
}

#[test]
fn test_message_objects_to_strings() {
    let mut message_strings: Vec<Vec<String>> = Vec::new();
    let mut message_header: Vec<String> = Vec::new();
    message_header.push("namespace".to_string());
    message_header.push("code".to_string());
    message_header.push("en".to_string());
    message_header.push("zh-cn".to_string());

    message_strings.push(message_header);

    let first_value: Vec<String> = Vec::from(
        ["test".to_string(), "test".to_string(), "en-test".to_string(), "中文".to_string()]
    );
    let second_value: Vec<String> = Vec::from(
        ["test".to_string(), "test2".to_string(), "".to_string(), "空".to_string()]
    );
    let third_value: Vec<String> = Vec::from(
        ["none".to_string(), "none".to_string()]
    );

    message_strings.push(first_value);
    message_strings.push(second_value);
    message_strings.push(third_value);

    let res = MessageObject::from_strings(&message_strings);

    let res = MessageObject::message_objects_to_strings(res);
    message_strings.remove(3);

    assert_eq!(res.len(), message_strings.len());
}

#[test]
fn test_strings_to_message_objects() {
    let mut message_strings: Vec<Vec<String>> = Vec::new();
    let mut message_header: Vec<String> = Vec::new();
    message_header.push("namespace".to_string());
    message_header.push("code".to_string());
    message_header.push("en".to_string());
    message_header.push("zh-cn".to_string());

    message_strings.push(message_header);

    let first_value: Vec<String> = Vec::from(
        ["test".to_string(), "test".to_string(), "en-test".to_string(), "中文".to_string()]
    );
    let second_value: Vec<String> = Vec::from(
        ["test".to_string(), "test2".to_string(), "".to_string(), "空".to_string()]
    );
    let third_value: Vec<String> = Vec::from(
        ["none".to_string(), "none".to_string()]
    );

    message_strings.push(first_value);
    message_strings.push(second_value);
    message_strings.push(third_value);

    let message_strings = message_strings;

    let res = MessageObject::from_strings(&message_strings);

    assert_eq!(res.len(), 2);

    // Build a new i18n instance.
    let mut i18n = I18n::new("zh-cn".to_string(), true, 2);
    for x in res {
        i18n.register_message_object(x);
    }

    let i18n = i18n;

    // First check
    let res = "en-test".to_string();
    let res_opt = Some(&res);
    assert_eq!(i18n.message(&"test".to_string(), &"test".to_string(), &"en".to_string()), res_opt);

    let res = "中文".to_string();
    let res_opt = Some(&res);
    assert_eq!(i18n.message(&"test".to_string(), &"test".to_string(), &"zh-cn".to_string()), res_opt);
}