package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type config struct {
	oauthToken  string
	nickname    string
	channelName string
}

type TwitchCommander struct {
	config *config
	conn   net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	messageMatcher *regexp.Regexp

	commands []*command
}

type command struct {
	name      string
	arguments []string
	handler   commandFunc

	matcher *regexp.Regexp
}

const (
	PRIVMSG_REGEX string = "^:[a-zA-Z0-9_]+![a-zA-Z0-9_]+@[a-zA-Z0-9_]+\\.tmi\\.twitch\\.tv\\ PRIVMSG\\ #%s\\ :(.*)$"
)

func CreateTwitchCommander(oauthToken, nickname, channelName string) *TwitchCommander {
	newConfig := &config{
		oauthToken:  oauthToken,
		nickname:    nickname,
		channelName: channelName,
	}

	return &TwitchCommander{
		config: newConfig,

		messageMatcher: regexp.MustCompile(fmt.Sprintf(PRIVMSG_REGEX, channelName)),
	}
}

type commandFunc func(map[string]string)

func (t *TwitchCommander) AddCommand(name string, handler commandFunc, args ...string) {
	commandRegex := "^[!@#$]?" + name
	if len(args) > 0 {
		commandRegex += "\\("
		commandArgRegexes := make([]string, len(args))
		for i, _ := range args {
			commandArgRegexes[i] = "(.*)"
		}
		commandRegex += strings.Join(commandArgRegexes, ",")
		commandRegex += "\\)"
	}
	commandRegex += "$"
	newCommand := &command{
		name:      name,
		arguments: args,
		handler:   handler,
		matcher:   regexp.MustCompile(commandRegex),
	}

	t.commands = append(t.commands, newCommand)
}

func (t *TwitchCommander) Connect() error {
	if t.conn != nil {
		return errors.New("Already connected!")
	}

	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return err
	}

	t.conn = conn
	t.reader = bufio.NewReader(conn)
	t.writer = bufio.NewWriter(conn)
	return nil
}

func (t *TwitchCommander) Close() {
	t.conn.Close()
}

func (t *TwitchCommander) ConnectAndListen() {
	t.Connect()
	defer t.Close()
	t.Authenticate()
	t.StartListener()
}

func (t *TwitchCommander) Authenticate() {
	t.Send("PASS oauth:" + t.config.oauthToken)
	t.Send("NICK " + t.config.nickname)
	t.Send("JOIN #" + t.config.channelName)
}

func (t *TwitchCommander) Send(message string) error {
	fmt.Fprintf(t.conn, message+"\r\n")
	return nil
}

const (
	PRIVMSG string = "PRIVMSG"
	PING    string = "PING"
)

func (t *TwitchCommander) StartListener() {
	for {
		line, _, err := t.reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			panic(err)
		} else {
			t.handleMessage(string(line))
		}
	}
}

func (t *TwitchCommander) handleMessage(message string) {
	if message == "PING :tmi.twitch.tv" {
		t.Send("PONG :tmi.twitch.tv")
		return
	}

	messageParts := strings.Split(message, " ")

	if messageParts[1] == "PRIVMSG" {
		t.parseCommand(messageParts[3][1:])
	}
}

func (t *TwitchCommander) parseCommand(commandString string) {
	for _, command := range t.commands {
		if commandArgs, matches := matchAndParseCommand(command, commandString); matches {
			command.handler(commandArgs)
			return
		}
	}
}

func matchAndParseCommand(c *command, cString string) (map[string]string, bool) {
	matchData := c.matcher.FindAllStringSubmatch(cString, -1)
	if len(matchData) == 0 {
		return nil, false
	}
	argReturn := make(map[string]string)
	for i, argName := range c.arguments {
		argReturn[argName] = matchData[0][i+1]
	}
	return argReturn, true
}
