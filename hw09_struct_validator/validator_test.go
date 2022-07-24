//go:build !codeanalysis

package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	//nolint:structcheck,unused
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	//nolint:lll
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "nil input",
			in:          nil,
			expectedErr: ErrNotStruct,
		},
		{
			name:        "not a struct input",
			in:          10,
			expectedErr: ErrNotStruct,
		},

		// User struct tests
		{
			name:        "User struct empty",
			in:          User{},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "User struct invalid ID, age, role and email (id len)",
			in:          User{ID: "abc", Age: 8, Name: "admin", Email: "admin@emailcom"},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "User struct invalid ID, age, role and email (age value)",
			in:          User{ID: "abc", Age: 8, Name: "admin", Email: "admin@emailcom"},
			expectedErr: ErrInvalidValue,
		},
		{
			name:        "User struct invalid ID, age, role and email (role value as default value exclude in tag)",
			in:          User{ID: "abc", Age: 8, Name: "admin", Email: "admin@emailcom"},
			expectedErr: ErrInvalidExcludeIn,
		},
		{
			name:        "User struct invalid ID, age, role and email (email regexp)",
			in:          User{ID: "abc", Age: 8, Name: "admin", Email: "admin@emailcom"},
			expectedErr: ErrInvalidRegexpNotMatch,
		},
		{
			name:        "User struct all valid",
			in:          User{ID: "aaaaaaaaaabbbbbbbbbbccccccccccdddddd", Age: 18, Name: "admin", Email: "admin@email.com", Role: "admin", Phones: []string{"89001234567", "89001234568"}},
			expectedErr: nil,
		},

		// App struct tests
		{
			name:        "App struct empty",
			in:          App{},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "App struct Version empty",
			in:          App{Version: ""},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "App struct Version short",
			in:          App{Version: "123"},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "App struct Version long",
			in:          App{Version: "123456"},
			expectedErr: ErrInvalidLen,
		},
		{
			name:        "App struct Version valid",
			in:          App{Version: "12345"},
			expectedErr: nil,
		},

		// Token struct tests
		{
			name:        "Token struct empty",
			in:          Token{},
			expectedErr: nil,
		},
		{
			name:        "Token struct Header valid",
			in:          Token{Header: []byte{1, 2, 3, 4}},
			expectedErr: nil,
		},

		// Response struct tests
		{
			name:        "Response struct empty",
			in:          Response{},
			expectedErr: ErrInvalidExcludeIn,
		},
		{
			name:        "Response struct Header invalid",
			in:          Response{Code: 10, Body: "abc"},
			expectedErr: ErrInvalidExcludeIn,
		},
		{
			name:        "Response struct Header valid 200",
			in:          Response{Code: 200, Body: "abc"},
			expectedErr: nil,
		},
		{
			name:        "Response struct Header valid 404",
			in:          Response{Code: 404, Body: "abc"},
			expectedErr: nil,
		},
		{
			name:        "Response struct Header valid 500",
			in:          Response{Code: 500, Body: "abc"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tin := tt.in
			tex := tt.expectedErr
			err := Validate(tin)
		out:
			switch verrs := err.(type) { //nolint:errorlint
			case ValidationErrors:
				for _, verr := range verrs {
					if errors.Is(verr.Err, tex) {
						break out
					}
				}
				t.Logf("invalid return value: in=%v expected=%v get=%v", tin, tex, verrs)
				t.Fail()
			case error:
				if !errors.Is(err, tex) {
					t.Logf("invalid return value: in=%v expected='%v' get='%v'", tin, tex, err)
					t.Fail()
				}
			case nil:
				if tex != nil {
					t.Logf("invalid return value: in=%v expected='%v' get='%v'", tin, tex, err)
					t.Fail()
				}
			default:
				t.Logf("unexpected return value")
				t.Fail()
			}
		})
	}
}
