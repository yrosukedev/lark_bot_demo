package main

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestHttpHeaderBuild(t *testing.T) {
	header := make(http.Header)
	header.Set("key 1", "value1, value2")
	values := header.Values("key 1")
	expected := []string{"value1", "value2"}
	if !reflect.DeepEqual(values, expected) {
		t.Fail()
	}
}

func TestStringSplit(t *testing.T) {
	if !reflect.DeepEqual(strings.Split("value1, value 2", ","), []string{"value1", " value 2"}) {
		t.Fatal("Split(\"value1, value 2\", \",\") != [\"value1\", \" value 2\"]")
	}

	if !reflect.DeepEqual(strings.Split("value 3", ","), []string{"value 3"}) {
		t.Fatal("Split(\"value 3\", \",\") != [\"value 3\"]")
	}
}
