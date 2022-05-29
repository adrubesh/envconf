package envconf

import "testing"

type TestConfigStruct struct {
	TestString string `env:"TEST_STRING_VALUE" default:"DEFAULT_STRING_VALUE"`
	TestInt    int    `env:"TEST_INT_VALUE" default:"0"`
	TestBool   bool   `env:"TEST_BOOL_VALUE" default:"true"`
}

func TestLoadConfig(t *testing.T) {
	var testConfig TestConfigStruct

	LoadConfig(&testConfig)

	if testConfig.TestBool != true {
		t.Errorf("failed to parse bool")
	}

	if testConfig.TestString != "DEFAULT_STRING_VALUE" {
		t.Errorf("failed to parse string")
	}

	if testConfig.TestInt != 0 {
		t.Errorf("failed to parse int")
	}
}
