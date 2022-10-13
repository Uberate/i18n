use i18n::file::{build_from_csv_file, write_to_csv_file};
use i18n::i18n::{I18n, MessageObject};

#[test]
fn test_build_from_csv() {
    let res = build_from_csv_file("./target/test_output_i18n.csv");
    match res {
        Ok(_) => {
            ()
        }
        Err(e) => {
            panic!("{}", e)
        }
    }
}

#[test]
fn test_write_to_csv() {
    let i18n_value = build_test_i18n();
    let res = write_to_csv_file(&i18n_value, "./target/test_output_i18n.csv");
    match res {
        Err(e) => {
            panic!("{}", e)
        }
        _ => ()
    }
}

fn build_test_i18n() -> I18n {
    let mut message_strings: Vec<Vec<String>> = Vec::new();
    let mut message_header: Vec<String> = Vec::new();
    message_header.push("namespace".to_string());
    message_header.push("code".to_string());
    message_header.push("en".to_string());
    message_header.push("zh-cn".to_string());

    message_strings.push(message_header);

    let value1: Vec<String> = Vec::from(
        ["test".to_string(), "test".to_string(), "en-test".to_string(), "中文".to_string()]
    );
    let value2: Vec<String> = Vec::from(
        ["test".to_string(), "test2".to_string(), "123".to_string(), "".to_string()]
    );
    let value3: Vec<String> = Vec::from(
        ["test2".to_string(), "test2".to_string(), "".to_string(), "v2".to_string()]
    );
    let value4: Vec<String> = Vec::from(
        ["test2".to_string(), "test".to_string(), "ttt".to_string(), "a1".to_string()]
    );
    let value5: Vec<String> = Vec::from(
        ["none".to_string(), "none".to_string()]
    );

    message_strings.push(value4);
    message_strings.push(value1);
    message_strings.push(value3);
    message_strings.push(value5);
    message_strings.push(value2);


    let message_strings = message_strings;

    let res = MessageObject::from_strings(&message_strings);


    // Build a new i18n instance.
    let mut i18n = I18n::new("zh-cn".to_string(), true, 2);
    for x in res {
        i18n.register_message_object(x);
    }

    i18n
}