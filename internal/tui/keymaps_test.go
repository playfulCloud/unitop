package tui

import (
	"strings"
	"testing"
)

func TestFooterTextIncludesLogsBinding(t *testing.T) {
	footer := footerText()

	if !strings.Contains(footer, "l logs") {
		t.Fatalf("expected footer to include logs binding, got %q", footer)
	}
}

func TestFooterTextMatchesExpectedLayout(t *testing.T) {
	footer := footerText()
	expected := "↑↓ move | / filter | l logs | r restart | s start | x stop | e/d enable | q quit"

	if footer != expected {
		t.Fatalf("expected footer %q, got %q", expected, footer)
	}
}

func TestFooterTextStaysCompact(t *testing.T) {
	footer := footerText()

	if len(footer) > 100 {
		t.Fatalf("expected compact footer, got %d characters: %q", len(footer), footer)
	}
}
