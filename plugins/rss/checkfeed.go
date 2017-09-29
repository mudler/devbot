package rss

import (
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mudler/devbot/bot"
)

func CheckFeed() {
	if atom, err := bot.DBListKeys("subscriptions"); err == nil {
		for _, rss := range atom {
			fp := gofeed.NewParser()
			if feed, err := fp.ParseURL(rss); err == nil {
				for _, item := range feed.Items {
					if v, err := bot.DBGetSingleKey(rss+"articles", item.Title); err == nil {
						if len(v) == 0 || item.Title != v {
							bot.DBPutKeyValue(rss+"articles", item.Title, item.Title)
							SendMessageToSubscribed(rss, item.Title+" - "+item.Link)
						}
					}
				}
			}
		}
	}
}

func SendMessageToSubscribed(atom, message string) {
	if destinations, err := bot.DBListKeys("atom" + atom); err == nil {
		for _, d := range destinations {
			time.Sleep(time.Second)
			bot.Conn.Privmsg(d, message)
		}
	}
}
