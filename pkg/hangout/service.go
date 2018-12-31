package hangout

import (
	gchat "google.golang.org/api/chat/v1"
)

type MessageBuilder struct {
	message gchat.Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (mb *MessageBuilder) Text(msg string) *MessageBuilder {
	mb.message.Text = msg
	return mb
}

func (mb *MessageBuilder) Thread(thread string) *MessageBuilder {
	t := &gchat.Thread{Name: thread}
	mb.message.Thread = t
	return mb
}

func (mb *MessageBuilder) Mention(userID, userName string, start int64) *MessageBuilder {
	u := &gchat.User{
		Name:        userID,
		DisplayName: userName,
		Type:        "HUMAN",
	}
	ann := &gchat.Annotation{
		Type: "USER_MENTION",
		UserMention: &gchat.UserMentionMetadata{
			Type: "MENTION",
			User: u,
		},
		StartIndex: start,
		Length:     int64(len(userName) + 1),
	}
	if mb.message.Annotations == nil {
		mb.message.Annotations = []*gchat.Annotation{
			ann,
		}
	} else {
		mb.message.Annotations = append(mb.message.Annotations, ann)
	}
	return mb
}

func (mb MessageBuilder) Build() *gchat.Message {
	return &mb.message
}

func (mb *MessageBuilder) Reset() *MessageBuilder {
	mb.message = gchat.Message{}
	return mb
}

type Service struct {
	spaces *gchat.SpacesService
}

func NewService(spacesClient *gchat.SpacesService) *Service {
	return &Service{spaces: spacesClient}
}

func (s *Service) SendMessageToTeam(teamID string, msg *gchat.Message) (string, error) {
	m, err := s.spaces.Messages.Create(teamID, msg).Do()
	if err != nil {
		return "", err
	}
	return m.Thread.Name, nil
}
