package home

import (
	tea "charm.land/bubbletea/v2"

	"bytes"
	"encoding/json"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"

	"fmt"
	"strings"
	"time"
)

type messageSentMsg message

func sendMessage(receiverID, userID uuid.UUID, body string, createdAt time.Time, accessToken string) tea.Cmd {
	return func() tea.Msg {
		msg := messageApi{
			CreatedAt:  createdAt,
			SenderID:   userID,
			ReceiverID: receiverID,
			Body:       body,
		}
		payload, err := json.Marshal(msg)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsChat, err)
		}
		res, err := client.Post(
			fmt.Sprintf("/conversations/%s/messages", receiverID),
			accessToken,
			bytes.NewReader(payload),
		)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsChat, err)
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsChat, res, "")
		}
		return nil
	}
}

func (h *home) updateChat(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			h.mode, h.input = modeList, nil

		case "enter":
			text := strings.TrimSpace(string(h.input))
			if text == "" {
				return h, nil
			}
			now := time.Now()
			c := &h.convs[h.cursor]
			c.Messages = append(c.Messages, message{
				from: "me",
				at:   formatMessageTime(now),
				body: text,
			})
			h.input = nil

			return h, sendMessage(c.ID, h.userID, text, now, h.accessToken)

		case "backspace":
			if len(h.input) > 0 {
				h.input = h.input[:len(h.input)-1]
			}
		default:
			if msg.Text != "" {
				h.input = append(h.input, []rune(msg.Text)...)
			}
		}
	case apperr.ChatErrMsg:
		h.chatErr = msg.UserErr
	}

	return h, nil
}

// view for the right side of the screen (chat window)

func (h *home) chatLines(frame int) []string {

	lines := make([]string, bodyH)
	if h.cursor >= len(h.convs) {
		lines[0] = styles.BannerStyle.Render("No conversations")
		lines[bodyH-1] = h.inputLine(frame)
		return lines
	}
	c := h.convs[h.cursor]

	head := "@ " + c.Name
	if c.IsGroup {
		head = "# " + c.Name // add the member count once your struct has it
	}
	lines[0] = styles.BannerStyle.Render(styles.Clip(head, rightW))

	var out []string
	for _, m := range c.Messages {
		out = append(out, renderMsg(m, rightW)...)
	}
	if avail := bodyH - 3; len(out) > avail {
		out = out[len(out)-avail:] // show the newest messages
	}
	copy(lines[1:], out)

	lines[bodyH-1] = h.inputLine(frame)
	return lines
}

func (h *home) inputLine(frame int) string {
	if h.mode != modeChat {
		return styles.RainDim.Render("> press enter to write a message")
	}
	r := h.input
	if maxText := rightW - 3; len(r) > maxText {
		r = r[len(r)-maxText:]
	}
	cur := " "
	if (frame/6)%2 == 0 {
		cur = "█"
	}
	return styles.RainBright.Render("> " + string(r) + cur)
}
