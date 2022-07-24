package hw09structvalidator

import (
	"reflect"
)

func validateSlice(val interface{}, tags []string) error {
	if val == nil {
		return nil
	}
	if reflect.TypeOf(val).Kind() != reflect.Slice {
		return ErrInvalidValue
	}

	var err error
	switch reflect.TypeOf(val).Elem().Kind() { //nolint:exhaustive
	case reflect.Int:
		for i := range val.([]int) {
			err = validateInt(val.([]int)[i], tags)
		}
	case reflect.String:
		for i := range val.([]string) {
			err = validateString(val.([]string)[i], tags)
		}
	default:
		err = ErrInvalidValue
	}

	if err != nil {
		return err
	}

	return nil
}
