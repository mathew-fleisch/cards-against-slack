# Cards Against Slack (bot)

Based on a list of regular expressions, a slack bot will be triggered to read a random line from two files. One card will fill in the blank of the other. Pipe/mount any three text files to set the "questions", "answers", and "triggers" and a "legacy app" slack token to substitute your own deck of cards.

### Configuration/Environment-Variables
```bash
# Required Environment Variable
SLACK_TOKEN
# Optional Environment Variable (defaults listed below)
QUESTIONS_FILE_URL="https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/questions.txt"
ANSWERS_FILE_URL="https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/answers.txt"
TRIGGERS_FILE_URL="https://raw.githubusercontent.com/mathew-fleisch/cards-against-slack/main/files/triggers.txt"
DISPLAY_ICON_URL="https://static.thenounproject.com/png/30134-200.png"
DISPLAY_USERNAME="Cards Against Slack"
```

### Build & Run

This bot can be run stand-alone as a go binary, or inside a container. Listed below are a few deployment options.

#### Go Binary

```bash
# Build and Run via Go Binary
go build bot.go
# Set environment variables here
./bot [path-to-questions.txt] [path-to-answers.txt] [path-to-triggers.txt]
```

#### Docker

```bash
# --------------------------------------------- #
#  Build and Run via Docker Container
export SLACK_TOKEN=xoxb-000000000000-000000000000-000000000000000000000000
export QUESTIONS_FILE_URL="https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/questions.txt"
export ANSWERS_FILE_URL="https://raw.githubusercontent.com/nodanaonlyzuul/against-humanity/master/answers.txt"
export TRIGGERS_FILE_URL="https://raw.githubusercontent.com/mathew-fleisch/cards-against-slack/main/files/triggers.txt"
export DISPLAY_ICON_URL="https://static.thenounproject.com/png/30134-200.png"
export DISPLAY_USERNAME="Cards Against Slack"
export CARDS_NAME=cards-build
docker build -t $CARDS_NAME .
docker run -it \
    -e SLACK_TOKEN=$SLACK_TOKEN \
    -e QUESTIONS_FILE_URL=$QUESTIONS_FILE_URL \
    -e ANSWERS_FILE_URL=$ANSWERS_FILE_URL \
    -e TRIGGERS_FILE_URL=$TRIGGERS_FILE_URL \
    -e DISPLAY_ICON_URL=$DISPLAY_ICON_URL \
    -e DISPLAY_USERNAME=$DISPLAY_USERNAME \
    -v ${PWD}:/cards \
    -w /cards \
    --name $CARDS_NAME \
    golang:alpine
```

#### Kubernetes

You must first set a kubernetes secret in the namespace you wish to deploy this to. The secret should be called 'slack' and have a key 'token' so that this deployment can map it to the environment variable 'SLACK_TOKEN'

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: cards
  name: cards
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cards
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: cards
    spec:
      containers:
      - image: mathewfleisch/cards-against-slack:v1.0.2
        name: cards-against-slack
        env:
        - name: DISPLAY_USERNAME
          value: Cards Against Humanity
        - name: DISPLAY_ICON_URL
          value: https://static.thenounproject.com/png/30134-200.png
        - name: QUESTIONS_FILE_URL
          value: https://raw.githubusercontent.com/cardsagainstcontainers/deck/master/questions.txt
        - name: ANSWERS_FILE_URL
          value: https://raw.githubusercontent.com/cardsagainstcontainers/deck/master/answers.txt
        - name: TRIGGERS_FILE_URL
          value: https://raw.githubusercontent.com/mathew-fleisch/cards-against-slack/main/files/triggers.txt
        - name: SLACK_TOKEN
          valueFrom:
            secretKeyRef:
              key: token
              name: slack
```


### Debug/Contribute
```bash
export SLACK_TOKEN=xoxb-000000000000-000000000000-000000000000000000000000
export CARDS_NAME=cards-debug
docker rm -f $(docker ps --all | grep $CARDS_NAME | awk '{print $1}') || true
docker run -it \
    -e SLACK_TOKEN=$SLACK_TOKEN \
    -v ${PWD}:/cards \
    -w /cards \
    --name $CARDS_NAME \
    golang:alpine

# Inside the container
apk update
apk add bash git vim
```

## LICENSE

ISC.
