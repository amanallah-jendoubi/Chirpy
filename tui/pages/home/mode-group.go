package home

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"
)

type groupCreated struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type groupMemberAdded struct {
	Username string
}

func (h *home) createGroup() tea.Cmd {
	name := strings.TrimSpace(string(h.gName))
	if name == "" {
		h.groupErr = "Enter a group name"
		return nil
	}
	if len([]rune(name)) > 80 {
		h.groupErr = "Group names can be at most 80 characters"
		return nil
	}

	accessToken := h.accessToken
	return func() tea.Msg {
		payload, err := json.Marshal(struct {
			Name string `json:"name"`
		}{Name: name})
		if err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		res, err := client.Post("/groups", accessToken, bytes.NewReader(payload))
		if err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsGroup, res, "")
		}
		var result groupCreated
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		return result
	}
}

func (h *home) updateGroup(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case groupCreated:
		for i, conversation := range h.convs {
			if conversation.ID == m.ID {
				h.cursor = i
				h.mode, h.gName, h.groupErr = modeChat, nil, ""
				return h, getConvMsgs(m.ID, m.Name, h.userID, h.accessToken)
			}
		}
		h.convs = append(h.convs, receiver{ID: m.ID, Name: m.Name, IsGroup: true})
		h.cursor = len(h.convs) - 1
		h.mode, h.gName, h.groupErr = modeChat, nil, ""
		return h, getConvMsgs(m.ID, m.Name, h.userID, h.accessToken)
	case tea.KeyPressMsg:
		switch m.String() {
		case "esc":
			h.mode, h.groupErr = modeList, ""
		case "enter":
			return h, h.createGroup()
		case "backspace":
			if len(h.gName) > 0 {
				h.gName = h.gName[:len(h.gName)-1]
			}
		default:
			if m.Text != "" && len(h.gName) < 80 {
				h.gName = append(h.gName, []rune(m.Text)...)
				h.groupErr = ""
			}
		}
	case tea.PasteMsg:
		pasted := strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ").Replace(m.Content)
		remaining := 80 - len(h.gName)
		if remaining > 0 {
			runes := []rune(pasted)
			if len(runes) > remaining {
				runes = runes[:remaining]
			}
			h.gName = append(h.gName, runes...)
		}
		h.groupErr = ""
	}
	return h, nil
}

func (h *home) startAddGroupMember() tea.Cmd {
	if h.cursor >= len(h.convs) || !h.convs[h.cursor].IsGroup {
		return nil
	}
	h.memberGroupID = h.convs[h.cursor].ID
	h.memberName = nil
	h.groupErr = ""
	h.mode = modeGroupMember
	return nil
}

func (h *home) addGroupMember() tea.Cmd {
	username := strings.TrimSpace(string(h.memberName))
	if username == "" {
		h.groupErr = "Enter a username"
		return nil
	}
	groupID := h.memberGroupID
	accessToken := h.accessToken

	return func() tea.Msg {
		userResponse, err := client.Get("/users/"+url.PathEscape(username), accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		defer userResponse.Body.Close()
		if userResponse.StatusCode < 200 || userResponse.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsGroup, userResponse, "")
		}
		var user struct {
			ID uuid.UUID `json:"user_id"`
		}
		if err := json.NewDecoder(userResponse.Body).Decode(&user); err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}

		payload, err := json.Marshal(struct {
			UserID uuid.UUID `json:"user_id"`
		}{UserID: user.ID})
		if err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		memberResponse, err := client.Post(
			fmt.Sprintf("/groups/%s/members", groupID),
			accessToken,
			bytes.NewReader(payload),
		)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsGroup, err)
		}
		defer memberResponse.Body.Close()
		if memberResponse.StatusCode < 200 || memberResponse.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsGroup, memberResponse, "")
		}
		return groupMemberAdded{Username: username}
	}
}

func (h *home) updateGroupMember(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case groupMemberAdded:
		h.groupErr = "Added @" + m.Username + " to the group"
		h.memberName = nil
	case apperr.GroupErrMsg:
		h.groupErr = m.UserErr
	case tea.KeyPressMsg:
		switch m.String() {
		case "esc":
			h.mode, h.groupErr = modeChat, ""
		case "enter":
			return h, h.addGroupMember()
		case "backspace":
			if len(h.memberName) > 0 {
				h.memberName = h.memberName[:len(h.memberName)-1]
			}
		default:
			if m.Text != "" && len(h.memberName) < 20 && !strings.ContainsAny(m.Text, " \t\r\n") {
				h.memberName = append(h.memberName, []rune(m.Text)...)
				h.groupErr = ""
			}
		}
	case tea.PasteMsg:
		pasted := strings.Map(func(r rune) rune {
			if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, m.Content)
		remaining := 20 - len(h.memberName)
		if remaining > 0 {
			runes := []rune(pasted)
			if len(runes) > remaining {
				runes = runes[:remaining]
			}
			h.memberName = append(h.memberName, runes...)
		}
		h.groupErr = ""
	}
	return h, nil
}

func (h *home) groupLines(frame int) []string {
	lines := make([]string, bodyH)
	lines[0] = styles.BannerStyle.Render("NEW GROUP")
	lines[2] = styles.RainMid.Render("Group name")

	cursor := " "
	if (frame/6)%2 == 0 {
		cursor = "█"
	}
	lines[3] = styles.RainBright.Render("▌ " + string(h.gName) + cursor)
	if h.groupErr != "" {
		lines[5] = styles.ErrStyle.Render(styles.Clip(h.groupErr, rightW))
	}
	return lines
}

func (h *home) groupMemberLines(frame int) []string {
	lines := make([]string, bodyH)
	lines[0] = styles.BannerStyle.Render("ADD GROUP MEMBER")
	if h.cursor < len(h.convs) {
		lines[2] = styles.RainMid.Render(styles.Clip("Group: "+h.convs[h.cursor].Name, rightW))
	}
	lines[4] = styles.RainMid.Render("Username")

	cursor := " "
	if (frame/6)%2 == 0 {
		cursor = "█"
	}
	lines[5] = styles.RainBright.Render("▌ @" + string(h.memberName) + cursor)
	if h.groupErr != "" {
		lines[7] = styles.ErrStyle.Render(styles.Clip(h.groupErr, rightW))
	}
	return lines
}
