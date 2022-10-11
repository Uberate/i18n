use crate::i18n::I18n;

#[test]
fn test_i18n() {
    let mut i = I18n::new("en".to_string(), true, 2);
    i.register_message(&"test".to_string(),
                       &"test".to_string(),
                       &"test".to_string(),
                       &"test".to_string());
    match i.message(&"test".to_string(), &"test".to_string(), &"test".to_string()) {
        Some(t) => { assert_eq!(t, "test".to_string()); }
        _ => { assert_eq!(1, 2); }
    }
}