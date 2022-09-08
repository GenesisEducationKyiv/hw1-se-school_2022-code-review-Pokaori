package main

import (
	"bitcoin-service/pkg/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env")
	code := m.Run()
	os.Exit(code)
}
