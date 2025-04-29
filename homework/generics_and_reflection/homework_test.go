package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	typeOfPerson := reflect.TypeOf(person)
	valueOfPerson := reflect.ValueOf(person)
	numberOfFields := typeOfPerson.NumField()
	var stringBuilder strings.Builder
	for i := 0; i < numberOfFields; i++ {
		value, present := typeOfPerson.Field(i).Tag.Lookup("properties")
		if present {
			values := strings.Split(value, ",")
			isFieldEmpty := valueOfPerson.Field(i).IsZero()
			switch {
			case isFieldEmpty && !containsOmitempty(values):
				zeroValueForField := reflect.Zero(typeOfPerson.Field(i).Type)
				stringBuilder.WriteString(fmt.Sprintf("%s=%v", values[0], zeroValueForField))
			case !isFieldEmpty:
				stringBuilder.WriteString(fmt.Sprintf("%s=%v", values[0], valueOfPerson.Field(i)))
			}
			if i != numberOfFields-1 && !(isFieldEmpty && containsOmitempty(values)) {
				stringBuilder.WriteString("\n")
			}
		}
	}
	return stringBuilder.String()
}

func containsOmitempty(values []string) bool {
	for _, value := range values {
		if value == "omitempty" {
			return true
		}
	}
	return false
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
