package hw09structvalidator

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrNotStruct  = errors.New("not a struct value")
	ErrInvalidTag = errors.New("invalid or unknown tag")

	ErrInvalidValue          = errors.New("invalid value")
	ErrInvalidLen            = errors.New("invalid value length")
	ErrInvalidRegexpNotMatch = errors.New("invalid value regexp not match")
	ErrInvalidExcludeIn      = errors.New("invalid value exclude in")
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	err := errors.New("validation errors")
	for i := range v {
		err = errors.Wrapf(v[i].Err, "%v: [%v]", err, v[i].Field)
	}
	return err.Error()
}

func Validate(v interface{}) error {
	var verrors ValidationErrors
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		var err error
		ftype := t.Field(i)
		fval := value.Field(i)

		if !ftype.IsExported() {
			continue
		}

		ftag := ftype.Tag.Get("validate")
		if ftag == "" {
			continue
		}
		ftags := strings.Split(ftag, "|")

		switch fval.Kind() { //nolint:exhaustive
		case reflect.Int:
			err = validateInt(int(fval.Int()), ftags)
		case reflect.String:
			err = validateString(fval.String(), ftags)
		case reflect.Slice:
			err = validateSlice(fval.Interface(), ftags)
		default:
			continue
		}

		if err != nil {
			verrors = append(verrors, ValidationError{Field: ftype.Name, Err: err})
		}
	}

	if len(verrors) == 0 {
		return nil
	}
	return verrors
}
