//! The mod file provide file operator for i18n tools. The project can ready or operate the i18n
//! message for file.

use std::cmp::Ordering;
use std::error::Error;

use crate::i18n::{I18n, MessageObject};

/// The function from_csv_file will build a I18n instance from a csv file. The csv file must
/// has a header to define the column language info. The csv file cloud be like this:
///
/// ```csv
/// namespace,code,en,zh-cn
/// test,test,en-value,中文
/// test,test2,,只有中文
/// test,test3,en-only,,
/// ```
///
/// The csv file's header must has 3 column at least, and first is namespace, second is code. The
/// language is defined start at third column. If the csv file header column less than 3, the I18n
/// is empty.
///
/// If some option is empty but with some 'space', it will not be skipped. The I18n builder will
/// drop empty line. And if some error occur at build time, the error will be returned. Such like
/// file read error, or csv value error.
pub fn from_csv_file(path: &str) -> Result<I18n, Box<dyn Error>> {
    let mut rdr = csv::Reader::from_path(path)?;
    let mut strings: Vec<Vec<String>> = Vec::new();
    for result in rdr.records() {
        let values = result?;
        let mut line: Vec<String> = Vec::new();
        for value in values.iter() {
            line.push(value.to_string())
        }
        strings.push(line)
    }
    let strings = strings;
    let mos = MessageObject::from_strings(&strings);

    let mut i18n = I18n::default();
    i18n.register_message_objects(mos);

    Ok(i18n)
}


/// The function to_csv_file will write all message from an [I18n] instance. And for help
/// human read, the message will order by namespace and code.
///
/// About the format, see the doc of [from_csv_file] function. If specify file is exists, the
/// file value will be covered. If the dir path not found, return the err.
pub fn to_csv_file(i18n: &I18n, path: &str) -> Result<(), Box<dyn Error>> {
    let mut message_objects = i18n.to_message_objects();
    let mut wtr = csv::Writer::from_path(path)?;
    if message_objects.len() == 0 {
        return Ok(());
    }

    message_objects.sort_unstable_by(|a, b| {
        let namespace_cmp_res = a.namespace().cmp(b.namespace());
        if namespace_cmp_res != Ordering::Equal {
            return namespace_cmp_res;
        }
        return a.code().cmp(b.code());
    });

    let strings_value = MessageObject::message_objects_to_strings(message_objects);

    for lines in strings_value {
        wtr.write_record(&lines)?
    }
    Ok(())
}