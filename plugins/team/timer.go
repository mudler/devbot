package team

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mudler/anagent"
	"github.com/mudler/devbot/bot"
)

func SetupTimer(s string, a *anagent.Anagent) {
	split := strings.Split(s, " ")
	timer_name := split[0]
	team := split[1]
	channel := split[2]
	now := time.Now()
	when, err := time.Parse("3:04PM", split[3])
	if err != nil {
		fmt.Println("Error setting up the timer ", s)
		return
	}

	when = time.Date(now.Year(), now.Month(), now.Day(), when.Hour(), when.Minute(), 0, 0, time.UTC)
	if when.Before(now) {
		when = when.AddDate(0, 0, 1)
	}

	fmt.Println("Timer: " + timer_name)
	fmt.Println("Team: " + team)
	fmt.Println("Channel: " + channel)
	fmt.Println("Time: " + split[3])
	fmt.Println("Each (minutes): " + split[4])

	recurring, _ := strconv.Atoi(split[4]) // minutes

	a.Timer(anagent.TimerID(timer_name), when, time.Duration(recurring)*time.Minute, true, func() {

		// Avoid to buzz on weekends
		day := time.Now().Weekday().String()
		if day == "Saturday" || day == "Sunday" {
			return
		}

		if atom, err := bot.DBListKeys("team" + team); err == nil {
			bot.Conn.Privmsg(channel, team+" "+timer_name+": "+strings.Join(atom, " "))
		}
	})
}
