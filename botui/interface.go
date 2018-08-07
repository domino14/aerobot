// The interface to the twitch bot

package botui

import (
	"html/template"
	"net/http"

	"github.com/domino14/aerobot/bot"
)

// IndexHandler just handles the call to the main interface page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("botui/template/index.html")
	t.Execute(w, nil)
}

type AerobotService struct {
	bot *bot.Aerobot
}

type StartArgs struct {
	Username  string `json:"username"`
	Channel   string `json:"channel"`
	LexiconDB string `json:"lexiconDb"`
	Password  string `json:"password"`
}

type StopArgs struct{}

type LoadLexiconArgs struct {
	NewLexiconDB string `json:"newLexiconDb"`
}

type GenericReply struct {
	Message string `json:"message"`
}

func (a *AerobotService) Start(r *http.Request, args *StartArgs,
	reply *GenericReply) error {
	if a.bot == nil {
		a.bot = new(bot.Aerobot)
	}
	// starts the robot. XXX Make sure channel has a # preceding it.
	err := a.bot.Init(args.Username, args.Password, args.Channel, args.LexiconDB)
	if err == nil {
		reply.Message = "Success"
	}
	return err

}

func (a *AerobotService) Stop(r *http.Request, args *StopArgs,
	reply *GenericReply) error {

	a.bot.Close()
	reply.Message = "Success"
	return nil
}

func (a *AerobotService) LoadLexicon(r *http.Request, args *LoadLexiconArgs,
	reply *GenericReply) error {

	err := a.bot.ChangeLexicon(args.NewLexiconDB)
	if err == nil {
		reply.Message = "Success"
	}
	return err
}
