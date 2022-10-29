package files

import (
	"github.com/uberate/i18n/pkg/provider"
	"github.com/uberate/mocker-utils/files"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var readers = map[string]func(string) (*provider.I18n, error){
	".json": ReadFromJSONFile,
}

// FromFiles will read the files in specify files and load to an i18n instance. By read order, the new instance will
// cover new instance.
func FromFiles(standard string, paths ...string) (*provider.I18n, error) {
	res := provider.NewI18n(standard)

	for _, pathItem := range paths {
		if !files.IsFileExists(pathItem) {
			continue
		}

		fi, err := os.Stat(pathItem)
		if err != nil {
			continue
		}
		if fi.IsDir() {
			var children []string
			_ = filepath.Walk(pathItem, func(pathStr string, info fs.FileInfo, err error) error {
				if pathStr == pathItem {
					return nil
				}
				children = append(children, pathStr)
				return nil
			})
			childI18n, err := FromFiles(standard, children...)
			if err != nil {
				return nil, err
			}
			res.CoveredMessage(childI18n)
		} else {
			if readFunc, ok := readers[strings.ToLower(path.Ext(pathItem))]; ok {
				instance, err := readFunc(pathItem)
				if err != nil {
					return nil, err
				}
				res.CoveredMessage(instance)
			}
		}
	}

	return res, nil
}
