package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rivo/uniseg"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	var token_file = "token.txt"
	if !(fileExists(token_file)) {
		fmt.Println("token file:", token_file, "was not found, so the file was created, please put your bot token there")
		emptyFile, err := os.Create(token_file)
		if err != nil {
			log.Fatal(err)
		}
		emptyFile.Close()
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(token_file)
	Token := string(content)

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		if string(err.Error()) == "websocket: close 4004: Authentication failed." {
			fmt.Println("Make sure your Token is valid!")
			fmt.Println("Token:", Token)
		}
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// if message sender is the bot, ignore message
	if m.Author.ID == s.State.User.ID {
		return
	}

	// my userid: 321028131982934017
	var owner_id = "321028131982934017"

	switch m.Content {
	case "!ping":
		s.ChannelMessageSend(m.ChannelID, "pong")

	case "!userid":
		// prints the id of the user doing !userid
		s.ChannelMessageSend(m.ChannelID, m.Author.ID)

	case "!ownertest":
		// tests if the user is the owner by comparing the userid to owner_id
		if m.Author.ID == owner_id {
			s.ChannelMessageSend(m.ChannelID, "you own the bot")
		} else {
			s.ChannelMessageSend(m.ChannelID, "you don't own the bot")
		}

	case "!jebaited":
		// you just got jebaited
		s.ChannelMessageSend(m.ChannelID, "http://www.gardling.com/coolvideo6")

	case "bruh":
		// responds to every message "bruh" with "moment"
		s.ChannelMessageSend(m.ChannelID, "moment")

	case "!website":
		// links to my website
		s.ChannelMessageSend(m.ChannelID, "<http://www.gardling.com>")

	case "!neofetch":
		// runs neofetch
		if m.Author.ID == owner_id {
			cmdstring := "/usr/bin/neofetch --stdout --color_blocks off | sed 's/\x1B[[0-9;?]*[a-zA-Z]//g' | sed '/^[[:space:]]*$/d'"

			fmt.Println("running neofetch")
			cmd := exec.Command("sudo", "su", "discord", "bash", "-c", cmdstring)

			out, err := cmd.CombinedOutput()
			s.ChannelMessageSend(m.ChannelID, "```\n"+string(out)+"\n```")
			if err != nil {
				error_str := string(err.Error())
				fmt.Println(error_str)
				s.ChannelMessageSend(m.ChannelID, "```\n"+string(error_str)+"\n```")
			}
			fmt.Println("ran neofetch")

		} else {
			s.ChannelMessageSend(m.ChannelID, "you have to be the owner to do that!")
		}
	case "!uptime":
		// does uptime -p
		cmdstring := "uptime -p | sed '/^[[:space:]]*$/d'"

		fmt.Println("running uptime")
		cmd := exec.Command("bash", "-c", cmdstring)
		out, err := cmd.CombinedOutput()

		s.ChannelMessageSend(m.ChannelID, "```\n"+string(out)+"\n```")
		if err != nil {
			error_str := string(err.Error())
			fmt.Println(error_str)
			s.ChannelMessageSend(m.ChannelID, "```\n"+string(error_str)+"\n```")
		}

	case "!stop":
		// kills the bot
		if m.Author.ID == owner_id {

			fmt.Println("shutting down bot")
			s.ChannelMessageSend(m.ChannelID, "shutting down bot!")
			os.Exit(0)
		} else {
			s.ChannelMessageSend(m.ChannelID, "you have to be the owner to shut down the bot!")
		}
	case "thx bot":
		// takes complements
		s.ChannelMessageSend(m.ChannelID, "np bro")

	case "!temps":
		// runs the temps command (my own shell script)
		cmdstring := "temps"

		fmt.Println("running temps")
		cmd := exec.Command("bash", "-c", cmdstring)
		out, err := cmd.CombinedOutput()

		s.ChannelMessageSend(m.ChannelID, "```\n"+string(out)+"\n```")
		if err != nil {
			error_str := string(err.Error())
			fmt.Println(error_str)
			s.ChannelMessageSend(m.ChannelID, "```\n"+string(error_str)+"\n```")
		}

	case "pog":
		// responds to "pog" with "poggers1"
		s.ChannelMessageSend(m.ChannelID, "poggers")

	case "!crab":
		// crab rave
		s.ChannelMessageSend(m.ChannelID, ":crab: :crab: :crab: :crab: :crab: :crab: :crab:")
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=LDU_Txk06tM")

	default:
		// bash stuff, bc why not?
		if strings.HasPrefix(m.Content, "!bash") {
			if m.Author.ID == owner_id {
				cmdstring := strings.Replace(m.Content, "!bash ", "", -1)

				cmd := exec.Command("sudo", "su", "discord", "bash", "-c", cmdstring)

				fmt.Println("running bash command:", cmdstring)
				out, err := cmd.CombinedOutput()

				if uniseg.GraphemeClusterCount(string(out)) > 2000 {
					s.ChannelMessageSend(m.ChannelID, "Sorry, the output of that command is over the 2000 character limit set by discord")
					return
				}

				s.ChannelMessageSend(m.ChannelID, "```\n"+string(out)+"\n```")
				if err != nil {
					error_str := string(err.Error())
					fmt.Println(error_str)
					s.ChannelMessageSend(m.ChannelID, "```\n"+error_str+"\n```")
				}

				fmt.Println("ran command")

			} else {
				s.ChannelMessageSend(m.ChannelID, "you have to be the owner to do that!")
			}
		}
	}

	// send message: s.ChannelMessageSend(m.ChannelID, message)
	// delete message: s.ChannelMessageDelete(m.ChannelID, m.ID)

}
