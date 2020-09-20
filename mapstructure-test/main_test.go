package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

const ISO8601 = "2006-01-02T15:04:05.000Z0700"

func TestMapstructure(t *testing.T) {
	input := map[string]interface{}{
		"name":        123,           // number => string
		"age":         "42",          // string => number
		"id":          "38654707451", // string => number
		"uuidSuccess": "73fa5786-2eb9-4fa7-953c-834cc18bb9cf",
		"uuidFailure": "a6a26705-5915-43c0-9e7d-743b0560abcd",
		"emails":      map[string]interface{}{}, // empty map => empty array
		"time":        "2020-09-20T14:02:39.222Z",
	}

	type TestInput struct {
		Name        string
		Age         int
		ID          int
		UUIDSuccess uuid.UUID `mapstructure:"uuidSuccess"`
		UUIDFailure uuid.UUID `mapstructure:"uuidFailure"`
		Time        time.Time
	}

	var result TestInput
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(ISO8601), // ISO 8601
			StringToUUIDHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           &result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}

	err = decoder.Decode(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", result)
	if result.UUIDSuccess.String() != "73fa5786-2eb9-4fa7-953c-834cc18bb9cf" {
		t.Error("Expected UUID")
	}

	if result.Time.Format(ISO8601) != input["time"].(string) {
		t.Errorf("Time does not match: %s vs %s", result.Time.Format(ISO8601), input["time"].(string))
	}
}

func StringToUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		return uuid.Parse(data.(string))
	}
}
