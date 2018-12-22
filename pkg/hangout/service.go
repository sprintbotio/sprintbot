package hangout

import (
	"context"
	gchat "google.golang.org/api/chat/v1"
)

type Service struct {
	spaces *gchat.SpacesService
}

func NewService(spacesClient *gchat.SpacesService)*Service  {
	return &Service{spaces:spacesClient}
}

func(s *Service)MonitorInBackground(ctx context.Context)  {

}
