package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	ValidationErrors map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.ValidationErrors) == 0
}

func (v *Validator) AddValidationError(key, message string) {
	if v.ValidationErrors == nil {
		v.ValidationErrors = make(map[string]string)
	}
	if _, ok := v.ValidationErrors[key]; !ok {
		v.ValidationErrors[key] = message
	}
}

func (v *Validator) CheckValue(ok bool, key, message string) {
	if !ok {
		v.AddValidationError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, choices ...int) bool {
	for i := range choices {
		if value == choices[i] {
			return true
		}
	}
	return false
}
