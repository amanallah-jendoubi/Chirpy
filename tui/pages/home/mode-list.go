package home

import (
	"encoding/json"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"
)

type convMsgs struct {
	ConvID uuid.UUID
	Msgs   []message
}

type messageApi struct {
	CreatedAt  time.Time `json:"created_at"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Body       string    `json:"body"`
}

func getConvMsgs(receiverID uuid.UUID, receiverName string, userID uuid.UUID, accessToken string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.Get(fmt.Sprintf("/conversations/%s/messages", receiverID), accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsChat, err)
		}
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsChat, res, "")
		}
		defer res.Body.Close()

		var messages []messageApi
		var msgs []message
		if err := json.NewDecoder(res.Body).Decode(&messages); err != nil {
			return apperr.HandleLocalError(apperr.AsChat, err)
		}
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			from := receiverName
			if msg.SenderID == userID {
				from = "me"
			}
			msgs = append(msgs, message{
				from: from,
				at:   formatMessageTime(msg.CreatedAt),
				body: msg.Body,
			})
		}
		return convMsgs{ConvID: receiverID, Msgs: msgs}
	}
}

func (h *home) updateList(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "q":
			return h, tea.Quit
		case "up", "k":
			if len(h.convs) == 0 {
				break
			}
			h.cursor = (h.cursor + len(h.convs) - 1) % len(h.convs)
			h.messageScroll = 0
			return h, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken)
		case "down", "j", "tab":
			if len(h.convs) == 0 {
				break
			}
			h.cursor = (h.cursor + 1) % len(h.convs)
			h.messageScroll = 0
			return h, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken)
		case "enter":
			h.mode = modeChat
		case "c":
			h.mode, h.uName, h.dmErr = modeDM, nil, ""
		case "n":
			h.mode, h.gName, h.groupErr = modeGroup, nil, ""
		}

	case apperr.ChatErrMsg:
		h.chatErr = m.UserErr
	}
	return h, nil
}

//sidebar conversation list view

func (h *home) leftLines() []string {
	lines := make([]string, bodyH)
	lines[0] = styles.BannerStyle.Render("CONVERSATIONS")

	rows := bodyH - 2
	start := 0
	if h.cursor >= rows {
		start = h.cursor - rows + 1
	}
	for i := 0; i < rows && start+i < len(h.convs); i++ {
		idx := start + i
		c := h.convs[idx]

		marker := "@ "
		if c.IsGroup {
			marker = "# "
		}
		label := styles.Clip(marker+c.Name, 17)
		text := padPlain(label, leftW-2)

		switch {
		case idx == h.cursor && h.mode == modeList:
			lines[2+i] = selectSt.Render("▌ " + text)
		case idx == h.cursor:
			lines[2+i] = styles.RainBright.Render("▌ " + text)
		default:
			lines[2+i] = styles.RainMid.Render("  " + text)
		}
	}
	return lines
}
