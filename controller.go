package main

import (
	"bufio"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/joho/godotenv"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var token string
var nickname string
var channel string

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	token = os.Getenv("TWITCH_OAUTH")
	nickname = os.Getenv("TWITCH_NICKNAME")
	channel = os.Getenv("TWITCH_CHANNEL")
	lines := make(chan string)
	conn := connect(lines)
	defer conn.Close()
	robotgo.SetKeyDelay(250)

	for line := range lines {
		if strings.Contains(line, "PING") {
			fmt.Fprintf(conn, "PONG :tmi.twitch.tv\r\n")
		} else {
			fmt.Println(line)
			data := strings.Split(line, " PRIVMSG ")
			if len(data) > 1 {
				message := strings.Split(data[1], fmt.Sprintf("#%s :", channel))
				if len(message) > 1 {
					parseCommand(message[1])
				}
			}
		}
	}
}

func connect(lines chan string) net.Conn {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				fmt.Println(err)
			} else {
				lines <- string(line)
			}
		}
	}()
	fmt.Fprintf(conn, "PASS oauth:%s\n", token)
	fmt.Fprintf(conn, "NICK %s\n", nickname)
	fmt.Fprintf(conn, "JOIN #%s\n", channel)
	return conn
}

var commands = []string{}
var toggleTime = time.Duration(25)
var mouseMultiplier = time.Duration(8)

func parseCommand(command string) {
	fmt.Printf("\"%s\"\n", command)
	re := regexp.MustCompile("(up$|down$|left$|right$|e$|r$|p\\((\\d+),(\\d+)\\)$|s\\((\\d+),(\\d+)\\)$)")
	commandStringsTop := re.FindAllStringSubmatch(command, -1)
	if len(commandStringsTop) <= 0 {
		return
	}
	commandStrings := commandStringsTop[0]
	if len(commandStrings) <= 0 {
		return
	}
	fmt.Println("\"" + commandStrings[1] + "\"")
	if commandStrings[1] == "e" {
		robotgo.KeyTap("e")
	} else if commandStrings[1] == "r" {
		robotgo.KeyTap("r")
	} else if commandStrings[1] == "up" {
		robotgo.KeyToggle("w", "down")
		time.Sleep(toggleTime * time.Millisecond)
		robotgo.KeyToggle("w", "up")
	} else if commandStrings[1] == "down" {
		robotgo.KeyToggle("s", "down")
		time.Sleep(toggleTime * time.Millisecond)
		robotgo.KeyToggle("s", "up")
	} else if commandStrings[1] == "right" {
		robotgo.KeyToggle("d", "down")
		time.Sleep(toggleTime * time.Millisecond)
		robotgo.KeyToggle("d", "up")
	} else if commandStrings[1] == "left" {
		robotgo.KeyToggle("a", "down")
		time.Sleep(toggleTime * time.Millisecond)
		robotgo.KeyToggle("a", "up")
	} else if strings.HasPrefix(commandStrings[1], "p") {
		x, _ := strconv.Atoi(commandStrings[2])
		y, _ := strconv.Atoi(commandStrings[3])
		if x > 1770 && x < 1835 && y > 50 && y < 110 {
			return
		}
		robotgo.Move(x, y)
		robotgo.MouseToggle("down", "left")
		time.Sleep(toggleTime * mouseMultiplier * time.Millisecond)
		robotgo.MouseToggle("up", "left")
	} else if strings.HasPrefix(commandStrings[1], "s") {
		x, _ := strconv.Atoi(commandStrings[4])
		y, _ := strconv.Atoi(commandStrings[5])
		robotgo.Move(x, y)
		robotgo.MouseToggle("down", "right")
		time.Sleep(toggleTime * mouseMultiplier * time.Millisecond)
		robotgo.MouseToggle("up", "right")
	}
}
