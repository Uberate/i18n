package i18n

import "fmt"

// FillStrings will make strings to specify length
func FillStrings(v []string, length int) []string {
	// The v is bigger or equals of length, return it directly.
	if len(v) >= length {
		return v
	}

	res := make([]string, length)

	for i := 0; i < len(v); i++ {
		res[i] = v[i]
	}
	for i := len(v); i < len(res); i++ {
		res[i] = ""
	}

	return res
}

// StringsMapperCopy will copy specify value from origin to target, and the target must bigger than mapper value.
func StringsMapperCopy(origin, target []string, mapper map[int]int, originOffset, targetOffset int) ([]string, error) {
	if mapper == nil {
		return target, nil
	}
	for originIndex, targetIndex := range mapper {
		if len(target) <= targetIndex+targetOffset {
			return nil, fmt.Errorf("The mapper of (target index + targetOffset)[%d] is out of range"+
				" len(target)[%d] ",
				targetIndex+targetOffset, len(target))
		}
		if len(origin) <= originIndex+originOffset {
			return nil, fmt.Errorf("The mapper of (origin index + originOffset)[%d] is out of range"+
				" len(origin)[%d] ",
				originIndex+originOffset, len(origin))
		}

		target[targetIndex+targetOffset] = origin[originIndex+originOffset]
	}

	return target, nil
}
