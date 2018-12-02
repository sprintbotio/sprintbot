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


func (ah *ActionHandler)Handle(m Message)(string,error){
	switch m.Platform() {
	case "hangout":
		return ah.handlers["hangout"].Handle(m)

	}
	return "", nil
}





type Handler interface {
	Handle(m Message)(string,error)
	Platform()string
}

type Message interface {
	Platform()string
}
