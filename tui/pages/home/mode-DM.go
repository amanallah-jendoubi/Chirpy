package home

/*import (
	tea "charm.land/bubbletea/v2"
	"encoding/json"
	"errors"
	"github.com/amanallah-jendoubi/Textio/tui/auth"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"
	"io"
	"strings"
)

func (h *home) startConversation() tea.Cmd {
	name := strings.TrimSpace(string(h.uName))
	if name == "" {
		h.uErr = "Enter a username"
		return nil
	}
	// already have a conversation with this user? just open it
	for i, c := range h.convs {
		if !c.IsGroup && strings.EqualFold(c.Name, name) {
			h.cursor, h.mode, h.input = i, modeChat, nil
			return nil
		}
	}
	return func() tea.Msg {
		accessToken, errMsg := auth.GetAccessToken()
		if errMsg.Err != nil {
			return errMsg
		}
		res, err := client.Get("/users/"+name, accessToken, nil)
		if err != nil {
			return client.ErrMsg{Err: errors.New("cannot reach server")}
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return client.HandleRequestError(res)
		}
		type receiver struct {
			ID uuid.UUID
		}
		var r receiver
		decoder := json.NewDecoder(io.LimitReader(res.Body, 1024))
		if err := decoder.Decode(&r); err != nil {
			return client.ErrMsg{Err: err}
		}
		h.convs = append(h.convs, conversation{
			ID:      r.ID,
			Name:    name,
			IsGroup: false,
		})
		h.cursor = len(h.convs) - 1
		h.mode, h.input = modeChat, nil
		return nil
	}
}

func (h *home) updateDM(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case client.ErrMsg:
		h.uErr = msg.Error()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			h.mode = modeList
		case "enter":
			return h, h.startConversation()
		case "backspace":
			if len(h.uName) > 0 {
				h.uName = h.uName[:len(h.uName)-1]
			}
		default:
			if msg.Text != "" && len(h.uName) < 20 && !strings.ContainsAny(msg.Text, " \t") {
				h.uName = append(h.uName, []rune(msg.Text)...)
				h.uErr = ""
			}
		}
	}
	return h, nil
}

//view for the new conversation form

func (h *home) dmLines(frame int) []string {
	lines := make([]string, bodyH)
	lines[0] = styles.BannerStyle.Render("NEW CONVERSATION")
	lines[2] = styles.RainMid.Render("Username")

	cur := " "
	if (frame/6)%2 == 0 {
		cur = "█"
	}
	lines[3] = styles.RainBright.Render("▌ @" + string(h.uName) + cur)

	if h.uErr != "" {
		lines[5] = styles.ErrStyle.Render(styles.Clip(h.uErr, rightW))
	}
	return lines
}
*/
