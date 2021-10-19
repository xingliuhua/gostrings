package gostrs

import "errors"

var NotFoundError = errors.New("data not found")
var allString = make(map[string]map[string]string)
var allStringArray = make(map[string]map[string][]string)

// ShouldGetString if no data is found, it will return "" and NotFoundError.
//
// eg: gostrs.ShouldGetString("en",r.xxx)
func ShouldGetString(language, key string) (string, error) {
	m, exist := allString[language]
	if !exist {
		return "", NotFoundError
	}
	value, exist := m[key]
	if !exist {
		return "", NotFoundError
	}
	return value, nil
}

// GetString if no data is found, it will return "".
//
// eg: gostrs.GetString("en",r.xxx)
func GetString(language, key string) string {
	m, exist := allString[language]
	if !exist {
		return ""
	}
	value, exist := m[key]
	if !exist {
		return ""
	}
	return value
}

// GetStringWithDefault if no data is found, it will return defaultValue.
//
// eg: gostrs.GetStringWithDefault("en",r.xxx,"no data")
func GetStringWithDefault(language, key, defaultValue string) string {
	m, exist := allString[language]
	if !exist {
		return defaultValue
	}
	value, exist := m[key]
	if !exist {
		return defaultValue
	}
	return value
}

// ShouldGetStringArray if no data is found, it will return nil and NotFoundError
//
// eg: gostrs.ShouldGetStringArray("en",r.xxx)
func ShouldGetStringArray(language, key string) ([]string, error) {
	m, exist := allStringArray[language]
	if !exist {
		return nil, NotFoundError
	}
	value, exist := m[key]
	if !exist {
		return nil, NotFoundError
	}
	return value, nil
}
// GetStringArray if no data is found, it will return empty slice
//
// eg: gostrs.GetStringArray("en",r.xxx)
func GetStringArray(language, key string) []string {
	m, exist := allStringArray[language]
	if !exist {
		return make([]string, 0)
	}
	value, exist := m[key]
	if !exist {
		return make([]string, 0)
	}
	return value
}

func SetData(allStringData map[string]map[string]string, allStringArrayData map[string]map[string][]string) {
	allString = allStringData
	allStringArray = allStringArrayData
}
