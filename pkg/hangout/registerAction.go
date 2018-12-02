package hangout

import "github.com/sirupsen/logrus"

func (ah *ActionHandler )handleRegister(m *Event) (string,error) {
	if err := ah.teamService.RegisterAdmin(m.Space.Name,m.User.Name); err != nil{
		logrus.Error("failed to register admin ", err)
		return "I was unable to register you", err
	}
	return "thank you for registering SprintBot. You have been added as the admin for this space", nil
}
