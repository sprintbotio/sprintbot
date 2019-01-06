package main

import (
	"context"
	"flag"

	"net/http"

	"github.com/sprintbot.io/sprintbot/pkg/standup"

	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/data/bolt"
	sprintBotGchat "github.com/sprintbot.io/sprintbot/pkg/gchat"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"github.com/sprintbot.io/sprintbot/pkg/user"
	"github.com/sprintbot.io/sprintbot/pkg/web"
	"golang.org/x/oauth2/google"
	gchat "google.golang.org/api/chat/v1"
)

var (
	logLevel string
	dbLoc    string
	platform string
)

func main() {
	flag.StringVar(&logLevel, "log-level", "debug", "use this to set log level: error, info, debug")
	flag.StringVar(&dbLoc, "db-loc", "./bot-db", "set the location of the db file")
	flag.StringVar(&platform, "platform", "hangouts", "choose the chat platform to target")
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

	db, err := bolt.Connect(dbLoc)
	if err != nil {
		panic(err)
	}
	defer bolt.Disconnect()

	if err := bolt.Setup(); err != nil {
		panic(err)
	}

	// generic stuff
	userRepo := bolt.NewUserRepository(db)
	teamRepo := bolt.NewTeamRespository(db)
	standUpRepo := bolt.NewStandUpRepo(db)
	scheduleRepo := bolt.NewStandUpRepository(db)
	chatActionHandler := chat.NewActionHandler()
	teamService := team.NewService(userRepo, teamRepo)
	userService := user.NewService(userRepo)
	router := web.BuildRouter()
	logger := logrus.StandardLogger()
	standupService := standup.NewStandUpService(teamRepo, scheduleRepo, standUpRepo)
	var runner standup.Runner

	if platform == "hangouts" {
		// hangout client

		gClient, err := google.DefaultClient(context.TODO(), "https://www.googleapis.com/auth/chat.bot")
		if err != nil {
			panic(err)
		}
		Gservice, err := gchat.New(gClient)
		if err != nil {
			panic(err)
		}
		spacesService := gchat.NewSpacesService(Gservice)
		hangoutService := sprintBotGchat.NewService(spacesService)
		runner = sprintBotGchat.NewStandUpRunner(hangoutService, teamService, standUpRepo)

		hangoutChatHandler := sprintBotGchat.NewActionHandler(teamService, standupService, userService)

		chatActionHandler.RegisterHandler(hangoutChatHandler)

		handler := web.NewHangoutHandler(chatActionHandler)
		web.MountHangoutHandler(router, handler)
		// register commands
		sprintBotGchat.NewRegisterationUseCase(userService, teamService).Register()
		sprintBotGchat.NewTeamUseCase(teamService).Register()
		sprintBotGchat.NewUserUseCases(userService).Register()
		sprintBotGchat.NewHelpUseCase().Register()
		sprintBotGchat.NewStandUpUseCase(standupService, userService).Register()

	}

	httpHandler := web.BuildHTTPHandler(router)

	go standupService.Schedule(context.TODO(), runner)

	//sys
	{
		web.MountSystemHandler(router)
	}

	logger.Println("starting api on 8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		logger.Fatal(err)
	}

}
