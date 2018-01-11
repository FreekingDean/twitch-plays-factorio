package main

import (
	//"github.com/go-vgo/robotgo"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	TOGGLE_TIME             time.Duration = time.Duration(25)
	MOUSE_TOGGLE_MULTIPLIER time.Duration = time.Duration(8)
)

func main() {
	commander := CreateTwitchCommander(
		os.Getenv("TWITCH_OAUTH"),
		os.Getenv("TWITCH_NICKNAME"),
		os.Getenv("TWITCH_CHANNEL_NAME"),
	)

	commander.AddCommand("inv", openInventory)
	commander.AddCommand("rot", rotateObject)

	commander.AddCommand("up", moveUp)
	commander.AddCommand("down", moveDown)
	commander.AddCommand("left", moveLeft)
	commander.AddCommand("right", moveRight)

	commander.AddCommand("p", primaryClick, "x", "y")
	commander.AddCommand("s", secondaryClick, "x", "y")
	commander.ConnectAndListen()
}

func openInventory(_ map[string]string) {
	//robotgo.KeyTap("e")
	log.Println("E")
}

func rotateObject(_ map[string]string) {
	//robotgo.KeyTap("r")
	log.Println("R")
}

func moveUp(_ map[string]string) {
	longKeystroke("w")
}
func moveDown(_ map[string]string) {
	longKeystroke("s")
}
func moveLeft(_ map[string]string) {
	longKeystroke("a")
}
func moveRight(_ map[string]string) {
	longKeystroke("d")
}

func longKeystroke(key string) {
	//robotgo.KeyToggle(key, "down")
	time.Sleep(TOGGLE_TIME * time.Millisecond)
	//robotgo.KeyToggle(key, "up")
	log.Println(key)
}

func primaryClick(args map[string]string) {
	toggleClick(args["x"], args["y"], "left")
}

func secondaryClick(args map[string]string) {
	toggleClick(args["x"], args["y"], "right")
}

func toggleClick(strX, strY, button string) {
	x, err := strconv.Atoi(strX)
	if err != nil {
		log.Printf("Error: %e\n", err)
		return
	}
	y, err := strconv.Atoi(strY)
	if err != nil {
		log.Printf("Error: %e\n", err)
		return
	}
	log.Println(x, y)
	log.Println(button)
	//robotgo.Move(x, y)
	//robotgo.MouseToggle("down", button)
	time.Sleep(TOGGLE_TIME * MOUSE_TOGGLE_MULTIPLIER * time.Millisecond)
	//robotgo.MouseToggle("up", button)
}
