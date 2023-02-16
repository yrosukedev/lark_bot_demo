package main

import (
	"context"
	"testing"
)

func TestLarkSendMessage(t *testing.T) {

	ctx := context.Background()
	
	if err := sendMessageToLark(ctx, "Hello, I'm Harry."); err != nil {
		t.Error(err)
	}
}
