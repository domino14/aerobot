package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent"
)

const serverssl = "irc.chat.twitch.tv"

func main() {
	var nick = flag.String("nickname", "", "Twitch Nickname")
	var pass = flag.String("password", "", "Twitch OAuth token")
	var channel = flag.String("channel", "", "Twitch channel (twitch.tv/{channel})")
	var databaseFile = flag.String("dbfile", "/Users/cesar/coding/webolith/db/CSW15.db",
		"Dictionary sqlite db file")
	flag.Parse()
	ircchan := "#" + *channel

	var err error

	db, err := sql.Open("sqlite3", *databaseFile)
	if err != nil {
		fmt.Printf("ERROR %s", err)
		return
	}

	irccon := irc.IRC(*nick, *nick) //Create new ircobj
	//Set options
	irccon.UseTLS = true //default is false
	//ircobj.TLSOptions //set ssl options
	irccon.Password = *pass
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	//Commands
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(ircchan) })
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

				rows, err := db.Query("SELECT definition FROM words WHERE word = ?",
					strings.ToUpper(word))
				if err != nil {
					fmt.Printf("[ERROR]", err)
					break
				}
				var definition string
				for rows.Next() {
					err = rows.Scan(&definition)
				}
				if definition != "" {
					irccon.SendRaw("PRIVMSG " + ircchan + " :" + definition)
				}
			}
		}
	})
	// ircobj.SendRaw("<string>") //sends string to server. Adds \r\n
	// ircobj.SendRawf("<formatstring>", ...) //sends formatted string to server.n

	err = irccon.Connect(serverssl + ":443")
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
