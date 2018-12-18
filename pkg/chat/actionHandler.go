package chat


type ActionHandler struct {
	handlers map[string]Handler
}

func (ah *ActionHandler)RegisterHandler(h Handler){
	ah.handlers[h.Platform()] = h
}

func NewActionHandler()*ActionHandler{
	return &ActionHandler{handlers: map[string]Handler{}}
}


func (ah *ActionHandler)Handle(m Message)string{
	switch m.Platform() {
	case "hangout":
		return ah.handlers["hangout"].Handle(m)

	}
	return ""
}





type Handler interface {
	Handle(m Message)string
	Platform()string
}

type Message interface {
	Platform()string
}
