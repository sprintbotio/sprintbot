package gchat

import (
	"time"

	"github.com/sprintbot.io/sprintbot/pkg/domain"
)

type Event struct {
	Type      string    `json:"type"`
	Token     string    `json:"token"`
	EventTime time.Time `json:"eventTime"`
	Space     struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Type        string `json:"type"`
	} `json:"space"`
	Message Message `json:"message"`
	User    User    `json:"user"`
}

type Message struct {
	Name   string `json:"name"`
	Sender struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		AvatarURL   string `json:"avatarUrl"`
		Email       string `json:"email"`
	} `json:"sender"`
	CreateTime   time.Time `json:"createTime"`
	Text         string    `json:"text"`
	ArgumentText string    `json:"argumentText"`
	Thread       struct {
		Name string `json:"name"`
	} `json:"thread"`
	Annotations []*Annotation `json:"annotations"`
}

func (e Event) Platform() string {
	return "hangout"
}

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Email       string `json:"email"`
}

type Annotation struct {
	Length      int    `json:"length"`
	StartIndex  int    `json:"startIndex"`
	Type        string `json:"type"`
	UserMention struct {
		Type string `json:"type"`
		User struct {
			AvatarURL   string `json:"avatarUrl"`
			DisplayName string `json:"displayName"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		} `json:"user"`
	} `json:"userMention"`
}

type command struct {
	actionType string
	teamID     string
	requester  string
	name       string
	args       []string
	room       string
	space      string
	team       *domain.Team
}

func (c command) NoEmptyArgs() bool {
	if len(c.args) == 0 {
		return true
	}
	for _, v := range c.args {
		if v == "" {
			return false
		}
	}
	return true
}
