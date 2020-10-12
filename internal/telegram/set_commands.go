package telegram

import (
	"encoding/json"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type botCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// setMyCommands changes the list of the bot's commands. This code
// is copy-pasted from
// https://github.com/go-telegram-bot-api/telegram-bot-api/commit/4a2c8c4547a868841c1ec088302b23b59443de2b
// Unfortunately there are no new releases: https://github.com/go-telegram-bot-api/telegram-bot-api/issues/378
func setMyCommands(bot *tgbotapi.BotAPI, commands []botCommand) error {
	v := url.Values{}
	data, err := json.Marshal(commands)
	if err != nil {
		return err
	}
	v.Add("commands", string(data))
	_, err = bot.MakeRequest("setMyCommands", v)
	if err != nil {
		return err
	}
	return nil
}

func (t *Telegram) setMyCommands() error {
	commands := []botCommand{
		// Homies
		{
			Command:     "newhomie",
			Description: "create a new homie (telegram bot user)",
		},
		{
			Command:     "listhomies",
			Description: "list homies (telegram bot users)",
		},
		{
			Command:     "deletehomie",
			Description: "delete an existing homie (telegram bot user)",
		},
		// Devices
		{
			Command:     "newdevice",
			Description: "create a new device",
		},
		{
			Command:     "listdevices",
			Description: "list devices",
		},
		{
			Command:     "deletedevice",
			Description: "delete an existing device",
		},
		// Subnets
		{
			Command:     "newsubnet",
			Description: "create a new subnet",
		},
		{
			Command:     "listsubnets",
			Description: "list subnets",
		},
		{
			Command:     "deletesubnet",
			Description: "delete an existing subnet",
		},
		// Server
		{
			Command:     "sethostname",
			Description: "set server endpoint hostname",
		},
		{
			Command:     "setport",
			Description: "set server listen port",
		},
		{
			Command:     "stats",
			Description: "show server info and device statistics",
		},
		{
			Command:     "restart",
			Description: "restart the server",
		},
		// Common
		{
			Command:     "cancel",
			Description: "cancel the current operation",
		},
		{
			Command:     "help",
			Description: "Show help",
		},
	}
	return setMyCommands(t.bot, commands)
}
