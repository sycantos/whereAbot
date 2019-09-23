package main

import (
	"log"
	"os"
	"fmt"
	"github.com/Krognol/go-wolfram"
	"github.com/christianrondeau/go-wit"
	"github.com/nlopes/slack"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const confidenceThreshold = 0.5
var (
	slackClient   *slack.Client
	witClient     *wit.Client
	wolframClient *wolfram.Client	
	err			  error
)

type DB struct {
    *sql.DB
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost)/test_db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	
	slack_token := "xoxb-765689382215-751070033618-zsg196Mr0o3cZX2hUcb4OwTU"
	slackClient  = slack.New(slack_token)
	witai_token := "APUKCLVTHDDGKMPT6LCS73FNDUBDIVRE"
	witClient = wit.NewClient(witai_token)
	wolframClient = &wolfram.Client{"XQ5AUJ-53AET27QQ5"}

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if len(ev.BotID) == 0 {
				go handleMessage(ev, db)
			}
		}
	}
}

func handleMessage(ev *slack.MessageEvent, db *sql.DB) {
	result, err := witClient.Message(ev.Msg.Text)
	if err != nil {
		log.Printf("unable to get wit.ai result: %v", err)
		return
	}

	var (
		topEntity    wit.MessageEntity
		topEntityKey string
	)

	for key, entityList := range result.Entities {
		for _, entity := range entityList {
			if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
				topEntity = entity
				topEntityKey = key
			}
		}
	}

	replyToUser(ev, topEntity, topEntityKey, db)
}

func replyToUser(ev *slack.MessageEvent, topEntity wit.MessageEntity, topEntityKey string, db *sql.DB) {
	switch topEntityKey {
	case "greetings":
		slackClient.PostMessage(ev.User, slack.MsgOptionText("Hello user! How can I help you?", false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{AsUser: true,}))
		return
	case "wolfram_search_query":
		res, err := wolframClient.GetShortAnswerQuery(topEntity.Value.(string), wolfram.Metric, 0)
		if err == nil {
			slackClient.PostMessage(ev.User, slack.MsgOptionText(res, false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{AsUser: true,}))
			log.Printf("Result: %v", res)
			return
		}
		log.Printf("unable to get data from wolfram: %v", err)
	case "how_to_learn":
		var tech_to_learn = topEntity.Value.(string)

		rows, err := db.Query("select link from resourcelist where tech_id = (select id from techtools where name = ?)", tech_to_learn)
		
		if err != nil {
			fmt.Printf("failed to enumerate tables: %v", err)
		}
		var table string
		for rows.Next() {
			if rows.Scan(&table) == nil {
				slackClient.PostMessage(ev.User, slack.MsgOptionText(table, false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{AsUser: true,}))
			}
		}
		return
	}
	slackClient.PostMessage(ev.User, slack.MsgOptionText("i dunno! sorry!", false), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{AsUser: true,}))
}