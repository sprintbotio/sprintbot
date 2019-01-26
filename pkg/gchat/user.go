package gchat

import (
	"strings"

	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/user"
)

type UserUseCases struct {
	userService *user.Service
}

func NewUserUseCases(userServices *user.Service) *UserUseCases {
	return &UserUseCases{userService: userServices}
}

func (uuc *UserUseCases) SetUsersTimeZone(cmd chat.Command, event *Event) (string, error) {
	updated := []string{}
	if _, ok := cmd.MappedArgs["self"]; ok {
		if err := uuc.userService.UpdateTZ(cmd.RequesterID, cmd.MappedArgs["zone"]); err != nil {
			return "", err
		}
		updated = append(updated, cmd.Requester)
		return `updated timezone to ` + cmd.MappedArgs["zone"] + ` for ` + strings.Join(updated, " , "), nil
	}
	for _, a := range event.Message.Annotations {
		if a != nil && (a.UserMention.User.DisplayName != "sprintbot" && a.UserMention.User.DisplayName != "") {
			userToAdd := a.UserMention.User.DisplayName
			userID := a.UserMention.User.Name
			if err := uuc.userService.UpdateTZ(userID, cmd.MappedArgs["zone"]); err != nil {
				return "", err
			}
			updated = append(updated, userToAdd)
		}
	}
	return `updated timezone to ` + cmd.MappedArgs["zone"] + ` for ` + strings.Join(updated, " , "), nil
}

func (uuc *UserUseCases) Register() {
	actionHandlers[chat.CommandSetUserTZ] = uuc.SetUsersTimeZone
}
