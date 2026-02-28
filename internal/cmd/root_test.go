package cmd

import (
	"testing"
)

func TestExecuteVersion(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"version"})
	if err != nil {
		t.Errorf("Execute(version) = %v, want nil", err)
	}
}

func TestExecuteHelp(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"--help"})
	if err != nil {
		t.Errorf("Execute(--help) = %v, want nil", err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	t.Parallel()

	err := Execute([]string{"nonexistent"})
	if err == nil {
		t.Error("Execute(nonexistent) = nil, want error")
	}
}

func TestBoolString(t *testing.T) {
	t.Parallel()

	if boolString(true) != "true" {
		t.Errorf("boolString(true) = %q, want true", boolString(true))
	}

	if boolString(false) != "false" {
		t.Errorf("boolString(false) = %q, want false", boolString(false))
	}
}

func TestEnvOr(t *testing.T) {
	t.Parallel()

	result := envOr("NAMECHEAP_CLI_NONEXISTENT_TEST_VAR_12345", "fallback")
	if result != "fallback" {
		t.Errorf("envOr(missing) = %q, want fallback", result)
	}
}
