package main

import "testing"

func TestLog(t *testing.T) {
	DebugLogger.Println("debug info")
	InfoLogger.Println("info")
	ErrorLogger.Println("error happened")
}
