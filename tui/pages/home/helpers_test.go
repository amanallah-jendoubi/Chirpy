package home

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"
)

func TestFormatMessageTime(t *testing.T) {
	got := formatMessageTime(time.Date(2026, time.August, 17, 14, 40, 0, 0, time.UTC))
	if want := "17 août 2026, 14:40"; got != want {
		t.Fatalf("formatMessageTime() = %q, want %q", got, want)
	}
}

func TestUpdateChatAddMemberShortcutDoesNotConsumeLetterA(t *testing.T) {
	groupID := uuid.New()
	h := &home{
		mode:  modeChat,
		convs: []receiver{{ID: groupID, IsGroup: true}},
	}

	h.updateChat(tea.KeyPressMsg(tea.Key{Text: "a", Code: 'a'}))
	if got := string(h.input); got != "a" {
		t.Fatalf("typed 'a' = %q, want it in the message draft", got)
	}
	if h.mode != modeChat {
		t.Fatalf("mode after typing 'a' = %v, want modeChat", h.mode)
	}

	h.updateChat(tea.KeyPressMsg(tea.Key{Code: 'a', Mod: tea.ModCtrl}))
	if h.mode != modeGroupMember {
		t.Fatalf("mode after Ctrl+A = %v, want modeGroupMember", h.mode)
	}
	if h.memberGroupID != groupID {
		t.Fatalf("memberGroupID = %v, want %v", h.memberGroupID, groupID)
	}
}
