package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func printEvents(eventsChannel <-chan *slacker.CommandEvent) {
	for event := range eventsChannel {
		fmt.Printf("timestamp: %v\n", event.Timestamp)
		fmt.Printf("command: %v\n", event.Command)
		fmt.Printf("parameters: %v\n", event.Parameters)
		fmt.Printf("event: %v\n", event.Event)
	}
}

func main() {
	botToken := os.Getenv("slot-bot-token")
	appToken := os.Getenv("slot-socket-token")

	bot := slacker.NewClient(botToken, appToken)

	go printEvents(bot.CommandEvents())

	bot.Command("i was born in <year>", &slacker.CommandDefinition{
		Description: "age calculator",
		Example:     "i was born in 2005",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				response.Reply("invalid year of birth.Ex: i was born in 2005")
				return
			}

			age := time.Now().Year() - yob
			if age < 0 {
				response.Reply("year exceeds current year")
				return
			}

			reply := fmt.Sprintf("You are %d years old", age)
			response.Reply(reply)
		},
	})

	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("pong")
		},
	})

	context, cancle := context.WithCancel(context.TODO())
	defer cancle()
	if err := bot.Listen(context); err != nil {
		log.Fatal("Error listening to bot:", err)
	}
}
