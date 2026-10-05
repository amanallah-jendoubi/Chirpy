package home

import (
	"strings"
	"time"

	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/auth"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"

	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"github.com/google/uuid"
)

type message struct {
	from string
	at   string
	body string
}
type receiver struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	IsGroup  bool      `json:"is_group"`
	Messages []message `json:"messages"`
}

type convsLoaded struct {
	convs  []receiver
	socket *client.WebSocket
}

type convsRefreshed []receiver

type realtimeMessageMsg client.RealtimeMessage

type websocketClosedMsg struct{}

type websocketConnectedMsg struct {
	socket *client.WebSocket
}

type viewMode int

const (
	modeList  viewMode = iota // browsing conversations
	modeChat                  // typing a message
	modeGroup                 // new-group form
	modeDM                    // new conversation by username
)

type home struct {
	convs         []receiver
	cursor        int
	messageScroll int
	mode          viewMode
	input         []rune // message draft
	accessToken   string
	socket        *client.WebSocket
	listErr       string
	chatErr       string
	groupErr      string
	dmErr         string
	userID        uuid.UUID

	// new-group form
	gName  []rune
	gFocus int // 0 = name, 1..len(contacts) = members, len(contacts)+1 = button

	// new conversation (by username)
	uID   string
	uName []rune
}

func NewHome() *home {
	return &home{}
}

func (h *home) Init() tea.Cmd {
	return func() tea.Msg {
		accessToken, err := auth.GetAccessToken()
		h.accessToken = accessToken
		if err != nil {
			return apperr.HandleLocalError(apperr.AsAuth, err)
		}
		//get userID
		res, err := client.Get("/users/me", accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		type user struct {
			ID uuid.UUID `json:"user_id"`
		}
		var u user
		if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		h.userID = u.ID
		res, err = client.Get("/conversations", accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsList, res, "")
		}
		var convs []receiver
		if err := json.NewDecoder(res.Body).Decode(&convs); err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		socket, _ := client.ConnectWebSocket(accessToken)
		return convsLoaded{convs: convs, socket: socket}
	}
}

func readWebSocket(socket *client.WebSocket) tea.Cmd {
	return func() tea.Msg {
		message, err := socket.ReadMessage()
		if err != nil {
			return websocketClosedMsg{}
		}
		return realtimeMessageMsg(message)
	}
}

func refreshConversations(accessToken string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.Get("/conversations", accessToken, nil)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsList, res, "")
		}
		var convs []receiver
		if err := json.NewDecoder(res.Body).Decode(&convs); err != nil {
			return apperr.HandleLocalError(apperr.AsList, err)
		}
		return convsRefreshed(convs)
	}
}

func reconnectWebSocket(accessToken string) tea.Cmd {
	return func() tea.Msg {
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		<-timer.C
		socket, _ := client.ConnectWebSocket(accessToken)
		return websocketConnectedMsg{socket: socket}
	}
}

func (h *home) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case convsLoaded:
		h.convs = m.convs
		h.socket = m.socket
		cmds := []tea.Cmd{reconnectWebSocket(h.accessToken)}
		if h.socket != nil {
			cmds[0] = readWebSocket(h.socket)
		}
		if len(h.convs) > 0 {
			cmds = append(cmds, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken))
		}
		return h, tea.Batch(cmds...)
	case realtimeMessageMsg:
		incoming := client.RealtimeMessage(m)
		conversationID := incoming.ReceiverID
		conversationIndex := -1
		for i := range h.convs {
			if h.convs[i].ID == conversationID {
				conversationIndex = i
				break
			}
		}
		if conversationIndex < 0 {
			if incoming.ReceiverID == h.userID {
				conversationID = incoming.SenderID
				for i := range h.convs {
					if h.convs[i].ID == conversationID {
						conversationIndex = i
						break
					}
				}
			}
		}
		if conversationIndex < 0 {
			return h, tea.Batch(readWebSocket(h.socket), refreshConversations(h.accessToken))
		}
		if conversationIndex >= 0 {
			from := h.convs[conversationIndex].Name
			if incoming.SenderID == h.userID {
				from = "me"
			}
			h.convs[conversationIndex].Messages = append(h.convs[conversationIndex].Messages, message{
				from: from,
				at:   formatMessageTime(incoming.CreatedAt),
				body: incoming.Body,
			})
		}
		return h, readWebSocket(h.socket)
	case convsRefreshed:
		h.convs = m
		if h.cursor >= len(h.convs) {
			h.cursor = 0
		}
		if len(h.convs) > 0 {
			return h, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken)
		}
	case websocketClosedMsg:
		h.socket = nil
		return h, reconnectWebSocket(h.accessToken)
	case websocketConnectedMsg:
		if m.socket == nil {
			return h, reconnectWebSocket(h.accessToken)
		}
		h.socket = m.socket
		return h, readWebSocket(h.socket)
	case convMsgs:
		for i := range h.convs {
			if h.convs[i].ID == m.ConvID {
				h.convs[i].Messages = m.Msgs
				break
			}
		}
	case apperr.ChatErrMsg:
		h.chatErr = m.UserErr
	case apperr.AuthErrMsg:
		h.listErr = m.UserErr
	case apperr.ListErrMsg:
		h.listErr = m.UserErr
	case tea.KeyPressMsg:
		switch h.mode {
		case modeChat:
			return h.updateChat(m)
		// case modeGroup:
		// 	return h.updateGroup(m)
		//case modeDM:
		//	return h.updateDM(m)
		default:
			return h.updateList(m)
		}
	case tea.PasteMsg:
		if h.mode == modeChat {
			return h.updateChat(m)
		}
	}
	return h, nil
}

// ---------- View ----------

func (h *home) Box(frame int) string {
	left := h.leftLines()
	var right []string
	switch h.mode {
	//case modeGroup:
	//	right = h.groupLines(frame)
	//case modeDM:
	//	right = h.dmLines(frame)
	default:
		right = h.chatLines(frame)
	}

	div := styles.RainDim.Render("│")
	rows := []string{pad(styles.BannerStyle.Render("Matrix Chat"), innerW)}
	for i := 0; i < bodyH; i++ {
		rows = append(rows, pad(left[i], leftW)+div+pad(right[i], rightW))
	}
	rows = append(rows, pad(styles.HelpStyle.Render(h.help()), innerW))

	return boxSt.Render(strings.Join(rows, "\n"))
}
