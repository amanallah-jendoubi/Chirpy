package home

import (
	"testing"
	"time"
)

func TestFormatMessageTime(t *testing.T) {
	got := formatMessageTime(time.Date(2026, time.August, 17, 14, 40, 0, 0, time.UTC))
	if want := "17 août 2026, 14:40"; got != want {
		t.Fatalf("formatMessageTime() = %q, want %q", got, want)
	}
}
