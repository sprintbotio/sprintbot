package chat

import "fmt"

type UnkownCommand struct {
	Command string
}

func (uc *UnkownCommand)Error() string {
	return fmt.Sprintf("uknown command %s use ```@sprintbot help``` to see available options", uc.Command)
}

func NewUknownCommand(c string) *UnkownCommand {
	return &UnkownCommand{Command:c}
}

func IsUnkownCommandErr(err error)  bool {
	_, ok := err.(*UnkownCommand)
	return ok
}

type MissingArgs struct {

}

func (ma *MissingArgs)Error()string  {
	return "the command was missing arguments"
}

func IsMissingArgsErr(err error)  bool {
	_, ok := err.(*MissingArgs)
	return ok
}