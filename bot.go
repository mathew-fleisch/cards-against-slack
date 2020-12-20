package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/slack-go/slack"
)

// Slacking off with global vars
var specials []func(event *slack.MessageEvent) bool
var questions []string
var answers []string
var triggers []string
var api *slack.Client
var rtm *slack.RTM
var channelsByName map[string]string

func makeChannelMap() {
	log.Println("Intializing...")
	channelsByName = make(map[string]string)
	channels, err := api.GetChannels(true)
	if err != nil {
		return
	}

	for _, v := range channels {
		channelsByName[v.Name] = v.ID
	}
	log.Println("Cards Against Slack Initialized! Read to deal...")
}

func parseUserMessage(event *slack.MessageEvent) bool {
	if !isTriggered(event.Text) {
		return false
	}
	var rejoinder = randomLine(questions) + "\n> " + randomLine(answers)

	sendMessage(event, rejoinder)
	return true
}

func isTriggered(msg string) bool {

	input := strip.StripTags(msg)

	// log.Printf("Input Before: `%s`", input)
	// log.Printf("Triggers: `%s`", triggers)
	for _, s := range triggers {
		// fmt.Println(i, s)
		input = regexp.MustCompile("^(?i)"+s+"$").ReplaceAllLiteralString(input, "")
	}
	log.Printf("Input After: `%s`", input)

	if len(input) == 0 {
		return true
	}

	return false
}

func sendMessage(event *slack.MessageEvent, msg string) {
	channelID, _, err := api.PostMessage(event.Channel,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionUsername("Cards Against Slack"),
		slack.MsgOptionTS(event.ThreadTimestamp),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			UnfurlLinks: true,
			UnfurlMedia: true,
		}))

	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	log.Printf("Sent to %s: `%s`", channelID, msg)
}

func handleMessage(event *slack.MessageEvent) {
	if event.SubType == "bot_message" {
		return
	}

	for _, handler := range specials {
		if handler(event) {
			break
		}
	}
}

func randomLine(textArr []string) string {
	rand.Seed(time.Now().Unix())
	return textArr[rand.Intn(len(textArr))]
}

func main() {
	// godotenv.Load(".env")

	qptr := flag.String("questions-path", "files/questions.txt", "file path to read questions from")
	aptr := flag.String("answers-path", "files/answers.txt", "file path to read answers from")
	tptr := flag.String("triggers-path", "files/triggers.txt", "file path to read triggers from")
	flag.Parse()
	questionsFile, err := ioutil.ReadFile(*qptr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	// fmt.Println("Questions File:", *qptr)
	questions = strings.Split(regexp.MustCompile("\n\n*").ReplaceAllLiteralString(string(questionsFile), "\n"), "\n")
	// fmt.Println("Questions:", questions)

	answersFile, err := ioutil.ReadFile(*aptr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	// fmt.Println("Answers File:", *aptr)
	answers = strings.Split(regexp.MustCompile("\n\n*").ReplaceAllLiteralString(string(answersFile), "\n"), "\n")
	// fmt.Println("answers:", answers)

	triggersFile, err := ioutil.ReadFile(*tptr)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	// fmt.Println("Triggers File:", *tptr)
	triggers = strings.Split(regexp.MustCompile("\n\n*").ReplaceAllLiteralString(string(triggersFile), "\n"), "\n")
	// fmt.Println("triggers:", triggers)
	// triggersSplit := strings.Split(string(triggers), "\n")

	// fmt.Println("Contents of questions file:", string(questions))
	// rand.Seed(time.Now().Unix())
	// fmt.Println("Random Question: ", questions[rand.Intn(len(questions))])
	// fmt.Println("Contents of answers file:", string(answers))
	// rand.Seed(time.Now().Unix())
	// fmt.Println("Random Answer: ", answers[rand.Intn(len(answers))])
	// fmt.Println("Triggers:", string(triggers))

	slacktoken, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("You must provide an access token in SLACK_TOKEN")
	}

	// Our special handlers. If they handled a message, they return true.
	specials = []func(event *slack.MessageEvent) bool{
		parseUserMessage,
	}

	api = slack.New(slacktoken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			makeChannelMap()

		case *slack.MessageEvent:
			// fmt.Printf("Message: %v\n", ev)
			handleMessage(ev)

		case *slack.PresenceChangeEvent:
			// fmt.Printf("Presence Change: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Fatal("Invalid credentials")

		case *slack.ConnectionErrorEvent:
			fmt.Printf("Event: %v\n", msg)
			log.Fatal("Can't connect")

		default:
			// Ignore other events..
			// fmt.Printf("Event: %v\n", msg)
		}
	}
}
