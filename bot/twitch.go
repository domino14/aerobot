package bot

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	irc "github.com/thoj/go-ircevent"
)

const serverssl = "irc.chat.twitch.tv"

// Aerobot encapsulates the connection and some current settings
type Aerobot struct {
	isRunning bool
	irccon    *irc.Connection
	ircchan   string
	db        *sql.DB
	wg        sync.WaitGroup
}

func (a *Aerobot) Init(nick string, pass string, ircchan string,
	databaseFile string) error {
	var err error

	if a.isRunning {
		return errors.New("Aerobot has already been initialized. Please quit and start again.")
	}

	a.db, err = sql.Open("sqlite3", databaseFile)
	log.Println("opened", databaseFile)
	if err != nil {
		return err
	}

	irccon := irc.IRC(nick, nick) //Create new ircobj
	if irccon == nil {
		return errors.New("failed to initialize")
	}
	//Set options
	irccon.UseTLS = true //default is false
	//ircobj.TLSOptions //set ssl options
	irccon.Password = pass
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	//Commands
	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join(ircchan)
		a.ircchan = ircchan
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	// irccon.AddCallback("PING", func(e *irc.Event) {
	// 	irccon.SendRaw("PONG :tmi.twitch.tv")
	// })/
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		//fmt.Printf("Got a msg, %s", e.Arguments)
		cmd := e.Arguments[1]

		if len(cmd) > 1 && cmd[0] == '!' {
			strs := strings.Split(cmd[1:], " ")
			switch strs[0] {
			case "define":
				if len(strs) <= 1 {
					break
				}
				word := strs[1]
				fmt.Printf("[debug] should define %s", word)

				rows, err := a.db.Query("SELECT definition FROM words WHERE word = ?",
					strings.ToUpper(word))
				if err != nil {
					fmt.Printf("[ERROR] %s", err)
					break
				}
				var definition string
				for rows.Next() {
					err = rows.Scan(&definition)
				}
				rows.Close()
				if definition != "" {
					irccon.SendRaw("PRIVMSG " + ircchan + " :" +
						strings.ToUpper(word) + ": " + definition)
				}
			}
		}
	})
	// ircobj.SendRaw("<string>") //sends string to server. Adds \r\n
	// ircobj.SendRawf("<formatstring>", ...) //sends formatted string to server.n

	err = irccon.Connect(serverssl + ":443")
	if err != nil {
		fmt.Printf("Err %s", err)
		return err
	}
	a.irccon = irccon
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.irccon.Loop() // XXX: this might reconnect; check this.
	}()
	return nil
}

func (a *Aerobot) Close() {
	a.irccon.ClearCallback("001")
	a.irccon.ClearCallback("366")
	a.irccon.ClearCallback("PRIVMSG")
	a.irccon.Quit()
	a.db.Close()
	a.wg.Wait()
	a.isRunning = false
}

func (a *Aerobot) ChangeLexicon(databaseFile string) error {
	a.db.Close()
	var err error
	a.db, err = sql.Open("sqlite3", databaseFile)
	if err != nil {
		return err
	}
	paths := strings.Split(databaseFile, "/")
	lexName := strings.Split(paths[len(paths)-1], ".")[0]
	a.irccon.SendRaw("PRIVMSG " + a.ircchan + " :" +
		"Changed lexicon to " + lexName)

	return nil
}
