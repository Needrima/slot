package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"slot/models"
	"strconv"
	"strings"
	"time"

	"github.com/shomali11/slacker"
)

var (
	tpl *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*"))
}

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

	bot.Command("covid info <country>", &slacker.CommandDefinition{
		Description: "covid info checker",
		Example:     "covid info Nigeria",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			country := strings.ToLower(request.Param("country"))
			uri := fmt.Sprintf(`https://coronavirus-19-api.herokuapp.com/countries/%s`, country)
			resp, err := http.Get(uri)
			if err != nil {
				fmt.Println("Resp error:", err)
				response.Reply("I'm sleeping. Try again later")
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Readall error:", err)
				response.Reply("I'm sleeping. Try again later")
				return
			}

			if string(body) == "Country not found" {
				response.Reply(string(body) + ".Ex: covid info Nigeria")
				return
			}

			var covidInfo models.CovidInfo
			if err := json.Unmarshal(body, &covidInfo); err != nil {
				fmt.Println("JSON decoder error:", err)
				response.Reply("Something went wrong. Try again later")
				return
			}

			buf := &bytes.Buffer{}

			if err := tpl.ExecuteTemplate(buf, "covidTemp.txt", covidInfo); err != nil {
				response.Reply("Something went wrong. Try again later")
				return
			}

			reply := buf.String()
			response.Reply(reply)
		},
	})

	bot.Command("weather <city>", &slacker.CommandDefinition{
		Description: "weather forecast",
		Example:     "weather Lagos",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			city := strings.ToLower(request.Param("city"))
			uri := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=da8d44e4505f5153cf700b5eeeb1885d&units=metric", city)
			resp, err := http.Get(uri)
			if err != nil {
				fmt.Println("Resp error:", err)
				response.Reply("I'm sleeping. Try again later")
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Readall error:", err)
				response.Reply("I'm sleeping. Try again later")
				return
			}

			var errMsg = map[string]interface{}{}
			if err := json.Unmarshal(body, &errMsg); err != nil {
				fmt.Println("JSON decoder error:", err)
				response.Reply("Something went wrong. Try again later")
				return
			}

			if errMsg["message"] == "city not found" {
				response.Reply("city not found. Ex: weather Lagos")
				return
			}

			var weatherInfo models.Weather
			if err := json.Unmarshal(body, &weatherInfo); err != nil {
				fmt.Println("JSON decoder error:", err)
				response.Reply("Something went wrong. Try again later")
				return
			}

			buf := &bytes.Buffer{}

			if err := tpl.ExecuteTemplate(buf, "weatherTemp.txt", weatherInfo); err != nil {
				response.Reply("Something went wrong. Try again later")
				return
			}

			reply := buf.String()
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
