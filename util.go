package gopo

import (
	"reflect"
	"strings"
	"unicode"
)

func toSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

func toLowerCamelCase(in string) string {
	runes := []rune(in)

	var out []rune
	flag := false
	for i, curr := range runes {
		if (i == 0 && unicode.IsUpper(curr)) || (flag && unicode.IsUpper(curr)) {
			out = append(out, unicode.ToLower(curr))
			flag = true
		} else {
			out = append(out, curr)
			flag = false
		}
	}

	return string(out)
}

// gorm
func parseGormTag2ColumnName(in string) string {
	structTag := reflect.StructTag(strings.Replace(in, "`", "", -1))

	setting := map[string]string{}
	for _, str := range []string{structTag.Get("sql"), structTag.Get("gorm")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	// Even it is ignored, also possible to decode db value into the field
	if value, ok := setting["COLUMN"]; ok {
		return value
	}
	return ""
}

// json
func parseJsonTag2JsonName(in string) string {
	structTag := reflect.StructTag(strings.Replace(in, "`", "", -1))
	jsonTag := structTag.Get("json")
	// json:"tag,hoge"
	if strings.Contains(jsonTag, ",") {
		// json:",hoge"
		if strings.HasPrefix(jsonTag, ",") {
			jsonTag = ""
		} else {
			jsonTag = strings.SplitN(jsonTag, ",", 2)[0]
		}
	}
	if jsonTag != "" {
		return jsonTag
	}
	return ""
}
