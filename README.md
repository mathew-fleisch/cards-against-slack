# Cards Against Slack (bot)

Based on a list of regular expressions, a slack bot will be triggered to read a random line from two files. One card will fill in the blank of the other. Pipe/mount any three text files to set the "questions", "answers", and "triggers" and a "legacy app" slack token to substitute your own deck of cards.

### Build
```bash
export CARDS_NAME=cards-build
docker build -t $CARDS_NAME .
```

### Debug
```bash
export CARDS_NAME=cards-debug
docker rm -f $(docker ps --all | grep $CARDS_NAME | awk '{print $1}') || true
docker run -it -e SLACK_TOKEN=$SLACK_TOKEN -v ${PWD}:/cards -w /cards --name $CARDS_NAME golang:alpine

apk update
apk add bash git vim
```

## LICENSE

ISC.
