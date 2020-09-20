package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/graphql-go/relay"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const ISO8601 = "2006-01-02T15:04:05.000Z0700"

type RelayResolvedGlobalIntID struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
}

func TestMapstructure(t *testing.T) {
	input := map[string]interface{}{
		"relayID":     "UXVlc3Rpb246Nzk0OTg=",
		"relayIntID":  "UXVlc3Rpb246Nzk0OTg=",
		"name":        123,           // number => string
		"age":         "42",          // string => number
		"id":          "38654707451", // string => number
		"uuidSuccess": "73fa5786-2eb9-4fa7-953c-834cc18bb9cf",
		"uuidFailure": "a6a26705-5915-43c0-9e7d-743b0560abcd",
		"emails":      map[string]interface{}{}, // empty map => empty array
		"time":        "2020-09-20T14:02:39.222Z",
	}

	type TestInput struct {
		RelayID     *relay.ResolvedGlobalID
		RelayIntID  *RelayResolvedGlobalIntID
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
			StringToRelayResolvedGlobalIDHookFunc(),
			StringToRelayResolvedGlobalIntIDHookFunc(),
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
	assert := assert.New(t)

	assert.Equal(result.UUIDSuccess.String(), input["uuidSuccess"])
	assert.Equal(result.Time.Format(ISO8601), input["time"])
	assert.Equal(result.RelayID.Type, "Question")
	assert.Equal(result.RelayIntID.Type, "Question")
	assert.Equal(result.RelayIntID.ID, 79498)
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

func StringToRelayResolvedGlobalIDHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(relay.ResolvedGlobalID{}) {
			return data, nil
		}

		resolvedID := relay.FromGlobalID(data.(string))
		if resolvedID == nil {
			return nil, fmt.Errorf("Failed to resolve: %s", data.(string))
		}
		return resolvedID, nil
	}
}

func StringToRelayResolvedGlobalIntIDHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(RelayResolvedGlobalIntID{}) {
			return data, nil
		}

		resolvedID := relay.FromGlobalID(data.(string))
		if resolvedID == nil {
			return nil, fmt.Errorf("Failed to resolve: %s", data.(string))
		}

		id, err := strconv.Atoi(resolvedID.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse id: %s", resolvedID)
		}

		return &RelayResolvedGlobalIntID{
			Type: resolvedID.Type,
			ID:   id,
		}, nil
	}
}
