package home

import (
	"strings"
	"time"

	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/matrix-chat/tui/auth"
	"github.com/amanallah-jendoubi/matrix-chat/tui/client"
	apperr "github.com/amanallah-jendoubi/matrix-chat/tui/errors"
	"github.com/amanallah-jendoubi/matrix-chat/tui/nav"

	"github.com/amanallah-jendoubi/matrix-chat/tui/styles"
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
	modeList        viewMode = iota // browsing conversations
	modeChat                        // typing a message
	modeGroup                       // new-group form
	modeDM                          // new conversation by username
	modeGroupMember                 // adding a member to a group
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
	gName         []rune
	memberName    []rune
	memberGroupID uuid.UUID

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
			cmds = append(cmds, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken, h.convs[h.cursor].IsGroup))
		}
		return h, tea.Batch(cmds...)
	case realtimeMessageMsg:
		return h, tea.Batch(readWebSocket(h.socket), refreshConversations(h.accessToken))
	case messageSentMsg:
		return h, refreshConversations(h.accessToken)
	case convsRefreshed:
		selectedID := uuid.Nil
		if h.cursor >= 0 && h.cursor < len(h.convs) {
			selectedID = h.convs[h.cursor].ID
		}
		h.convs = m
		if selectedID != uuid.Nil {
			for i := range h.convs {
				if h.convs[i].ID == selectedID {
					h.cursor = i
					break
				}
			}
		} else if h.cursor >= len(h.convs) {
			h.cursor = 0
		}
		if len(h.convs) > 0 {
			return h, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken, h.convs[h.cursor].IsGroup)
		}
	case dmUserLoaded:
		return h.updateDM(m)
	case groupCreated:
		return h.updateGroup(m)
	case groupMemberAdded:
		return h.updateGroupMember(m)
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
	case apperr.DMErrMsg:
		h.dmErr = m.UserErr
	case apperr.GroupErrMsg:
		h.groupErr = m.UserErr
	case tea.KeyPressMsg:
		switch h.mode {
		case modeChat:
			return h.updateChat(m)
		case modeDM:
			return h.updateDM(m)
		case modeGroup:
			return h.updateGroup(m)
		case modeGroupMember:
			return h.updateGroupMember(m)
		default:
			return h.updateList(m)
		}
	case tea.PasteMsg:
		switch h.mode {
		case modeChat:
			return h.updateChat(m)
		case modeGroup:
			return h.updateGroup(m)
		case modeGroupMember:
			return h.updateGroupMember(m)
		}
	}
	return h, nil
}

// ---------- View ----------

func (h *home) Box(frame int) string {
	left := h.leftLines()
	var right []string
	switch h.mode {
	case modeGroup:
		right = h.groupLines(frame)
	case modeGroupMember:
		right = h.groupMemberLines(frame)
	case modeDM:
		right = h.dmLines(frame)
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
