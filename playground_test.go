package main

import (
	"encoding/json"
	"testing"
)

func TestMarshal(t *testing.T) {

	body := map[string]string{
		"msg": "123",
	}
	bodyData, err := json.Marshal(body)
	if err != nil {
		t.Fatal()
	}

	if string(bodyData) != "{\"msg\":\"123\"}" {
		t.Fatal()
	}
	
}
