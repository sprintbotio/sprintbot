package gchat

import "github.com/sprintbot.io/sprintbot/pkg/chat"

type HelpUseCase struct {
}

const (
	adminHelpResp = `available commands are:
` + "```@sprinbot add <user1>,<user2> to team```" + `
` + "```@sprintbot remove <user1>,<user2> from team```" + `
` + "```@sprintbot view team```" + `
` + "```@sprintbot schedule standup <HH:MM> <TimeZone>```" + `
`
	generalHelpResp = `available commands are:
` + "```@sprintbot view standup log```" + `
` + "```(DM Only) log standup status <your status message> ```" + `
`
)

func (hu *HelpUseCase) AdminHelp(cmd chat.Command, event *Event) (string, error) {
	return adminHelpResp, nil
}

func (hu *HelpUseCase) GeneralHelp(cmd chat.Command, event *Event) (string, error) {
	return generalHelpResp, nil
}

func (hu *HelpUseCase) Register() {
	actionHandlers[chat.CommandGeneralHelp] = hu.GeneralHelp
	actionHandlers[chat.CommandAdminHelp] = hu.AdminHelp
}

func NewHelpUseCase() *HelpUseCase {
	return &HelpUseCase{}
}
