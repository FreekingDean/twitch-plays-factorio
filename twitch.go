package main

type config struct {
	oauthToken  string
	nickname    string
	channelName string
}

type TwitchCommander struct {
	config config

	commands []*command
}

type command struct {
	name      string
	arguments []string
}

func createTwitchCommander(oauthToken, nickname, channelName string) *TwitchCommander {
	newConfig := &config{
		oauthToken:  oauthToken,
		nickname:    nickname,
		channelName: channelName,
	}

	return &TwitchCommander{
		config: newConfig,
	}
}

func (t *TwitchCommander) AddCommand(name string, args ...string) {
	newCommand := &command{
		name:      name,
		arguments: args,
	}
	t.commands = append(t.commands, newCommand)
}
