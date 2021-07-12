package utils

import (
	"reflect"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func RemoveDuplicateTags(t []string) []string {
	tmp := removeDuplicateValuesInSlice(t)
	tags := make([]string, len(tmp))
	for i, v := range tmp {
		tags[i] = v.(string)
	}
	return tags
}

func removeDuplicateValuesInSlice(t interface{}) []interface{} {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		unq := make(map[interface{}]bool)
		for i := 0; i < s.Len(); i++ {
			if _, ok := unq[s.Index(i).Interface()]; !ok {
				unq[s.Index(i).Interface()] = true
			}
		}

		keys := make([]interface{}, 0)
		for key := range unq {
			keys = append(keys, key)
		}
		return keys
	default:
		return nil
	}
}
