package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/data/memory"
	"github.com/sprintbot.io/sprintbot/pkg/hangout"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"github.com/sprintbot.io/sprintbot/pkg/web"
	"net/http"
)

var logLevel string

func main() {
	flag.StringVar(&logLevel, "log-level", "debug", "use this to set log level: error, info, debug")
	flag.Parse()
	switch logLevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Info("log-level set to info")
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Error("log-level set to error")
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("log-level set to debug")
	default:
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Error("log-level set to error")
	}
	logrus.SetLevel(logrus.InfoLevel)
	router := web.BuildRouter()
	logger := logrus.StandardLogger()
	httpHandler := web.BuildHTTPHandler(router)

	userRepo := memory.NewUserRespository()
	chatActionHandler := chat.NewActionHandler()
	teamService := team.NewService(userRepo)

	//sys
	{
		web.MountSystemHandler(router)
	}

	// hangouts
	{
		hangoutChatHandler := hangout.NewActionHandler(teamService)
		chatActionHandler.RegisterHandler(hangoutChatHandler)
		handler := web.NewHangoutHandler(chatActionHandler)
		web.MountHangoutHandler(router,handler)
	}

	logger.Println("starting api on 8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		logger.Fatal(err)
	}

}