package home

import (
	"strings"

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

type convsLoaded []receiver

type viewMode int

const (
	modeList  viewMode = iota // browsing conversations
	modeChat                  // typing a message
	modeGroup                 // new-group form
	modeDM                    // new conversation by username
)

type home struct {
	convs       []receiver
	cursor      int
	mode        viewMode
	input       []rune // message draft
	accessToken string
	listErr     string
	chatErr     string
	groupErr    string
	dmErr       string
	userID      uuid.UUID

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
		return convsLoaded(convs)
	}
}

func (h *home) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch m := msg.(type) {
	case convsLoaded:
		h.convs = m
		if len(h.convs) > 0 {
			return h, getConvMsgs(h.convs[h.cursor].ID, h.convs[h.cursor].Name, h.userID, h.accessToken)
		}
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
