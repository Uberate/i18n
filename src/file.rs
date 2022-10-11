//! The mod file provide file operator for i18n tools. The project can ready or operate the i18n
//! message for file.

use crate::i18n::I18n;

/// The function build_from_csv_file will build a I18n instance from a csv file. The csv file must
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
pub fn build_from_csv_file() -> Result<I18n, String> {
    Ok(I18n::new("en".to_string(), false, 10))
}

pub fn write_to_csv_file(i18n: &I18n) {}