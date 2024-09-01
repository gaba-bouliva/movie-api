package validator

import "regexp"

var (
EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func (v Validator) New() *Validator{
	return &Validator{
		make(map[string]string, 0),
	}
}

func (v Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddErr(key, errMessage string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = errMessage
	}
}

func (v *Validator) Check(ok bool, key, errMessage string) {
	if !ok {
		v.AddErr(key, errMessage)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}

	return false
}

func (v Validator) Unique(values []string) bool {
	uniqeValues := make(map[string]int)

	for i, val := range values {
		if _, exists := uniqeValues[val]; exists {
			return false
		}
			uniqeValues[val] = i
	}

	return true
}

func (v Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}