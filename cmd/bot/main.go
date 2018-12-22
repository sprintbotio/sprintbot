package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/sprintbot.io/sprintbot/pkg/chat"
	"github.com/sprintbot.io/sprintbot/pkg/data/bolt"
	"github.com/sprintbot.io/sprintbot/pkg/hangout"
	"github.com/sprintbot.io/sprintbot/pkg/team"
	"github.com/sprintbot.io/sprintbot/pkg/user"
	"github.com/sprintbot.io/sprintbot/pkg/web"
	"golang.org/x/oauth2/google"
	gchat "google.golang.org/api/chat/v1"
	"net/http"
)

var (
	logLevel string
	dbLoc string
)

func main() {
	flag.StringVar(&logLevel, "log-level", "debug", "use this to set log level: error, info, debug")
	flag.StringVar(&dbLoc, "db-loc", "./bot-db", "set the location of the db file")
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
	if err != nil{
		panic(err)
	}
	defer bolt.Disconnect()

	if err := bolt.Setup(); err != nil{
		panic(err)
	}


	// hangout background jobs
	{
		gClient, err := google.DefaultClient(context.TODO(),"https://www.googleapis.com/auth/chat.bot")
		if err != nil{
			panic(err)
		}
		Gservice,err := gchat.New(gClient)
		if err != nil{
			panic(err)
		}
		spacesService := gchat.NewSpacesService(Gservice)
		hangoutService := hangout.NewService(spacesService)
		hangoutService.MonitorInBackground(context.TODO())
	}


	router := web.BuildRouter()
	logger := logrus.StandardLogger()
	httpHandler := web.BuildHTTPHandler(router)

	userRepo := bolt.NewUserRepository(db)
	teamRepo := bolt.NewTeamRespository(db)
	chatActionHandler := chat.NewActionHandler()
	teamService := team.NewService(userRepo, teamRepo)
	userService := user.NewService(userRepo)

	//sys
	{
		web.MountSystemHandler(router)
	}

	// hangouts
	{
		hangoutChatHandler := hangout.NewActionHandler(teamService,userService)
		chatActionHandler.RegisterHandler(hangoutChatHandler)
		handler := web.NewHangoutHandler(chatActionHandler)
		web.MountHangoutHandler(router,handler)
	}

	logger.Println("starting api on 8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		logger.Fatal(err)
	}

}