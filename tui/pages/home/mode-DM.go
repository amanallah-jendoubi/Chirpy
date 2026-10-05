package home

import (
	"encoding/json"
	"net/url"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"
)

type dmUserLoaded struct {
	id   uuid.UUID
	name string
}

func (h *home) startConversation() tea.Cmd {
	name := strings.TrimSpace(string(h.uName))
	if name == "" {
		h.dmErr = "Enter a username"
		return nil
	}

	for i, conversation := range h.convs {
		if !conversation.IsGroup && strings.EqualFold(conversation.Name, name) {
			h.cursor, h.mode, h.input = i, modeChat, nil
			return getConvMsgs(conversation.ID, conversation.Name, h.userID, h.accessToken)
		}
	}

	return func() tea.Msg {
		res, err := client.Get("/users/"+url.PathEscape(name), h.accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsDM, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsDM, res, "")
		}

		var user struct {
			ID uuid.UUID `json:"user_id"`
		}
		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			return apperr.HandleLocalError(apperr.AsDM, err)
		}
		return dmUserLoaded{id: user.ID, name: name}
	}
}

func (h *home) updateDM(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case dmUserLoaded:
		for i, conversation := range h.convs {
			if !conversation.IsGroup && conversation.ID == m.id {
				h.cursor, h.mode, h.input = i, modeChat, nil
				return h, getConvMsgs(conversation.ID, conversation.Name, h.userID, h.accessToken)
			}
		}
		h.convs = append(h.convs, receiver{ID: m.id, Name: m.name})
		h.cursor = len(h.convs) - 1
		h.mode, h.input, h.dmErr = modeChat, nil, ""
		return h, getConvMsgs(m.id, m.name, h.userID, h.accessToken)
	case apperr.DMErrMsg:
		h.dmErr = m.UserErr
	case tea.KeyPressMsg:
		switch m.String() {
		case "esc":
			h.mode, h.dmErr = modeList, ""
		case "enter":
			return h, h.startConversation()
		case "backspace":
			if len(h.uName) > 0 {
				h.uName = h.uName[:len(h.uName)-1]
			}
		default:
			if m.Text != "" && len(h.uName) < 20 && !strings.ContainsAny(m.Text, " \t\r\n") {
				h.uName = append(h.uName, []rune(m.Text)...)
				h.dmErr = ""
			}
		}
	}
	return h, nil
}

func (h *home) dmLines(frame int) []string {
	lines := make([]string, bodyH)
	lines[0] = styles.BannerStyle.Render("NEW CONVERSATION")
	lines[2] = styles.RainMid.Render("Username")

	cursor := " "
	if (frame/6)%2 == 0 {
		cursor = "█"
	}
	lines[3] = styles.RainBright.Render("▌ @" + string(h.uName) + cursor)

	if h.dmErr != "" {
		lines[5] = styles.ErrStyle.Render(styles.Clip(h.dmErr, rightW))
	}
	return lines
}
