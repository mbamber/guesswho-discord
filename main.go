package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	players     map[discordgo.User]string
	playerOrder []discordgo.User
	Token       string
	started     bool
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	channel, _ := s.Channel(m.ChannelID)
	if channel.Type != discordgo.ChannelTypeDM {
		return
	}

	// Make commands case insensitve
	msg := strings.ToLower(m.Content)

	cmd := strings.TrimSpace(msg)

	switch cmd {
	case "new":
		fmt.Println("New game started")
		players = map[discordgo.User]string{}
	case "join":
		for player := range players {
			if player == *m.Author {
				fmt.Printf("%s tried to join the game, but they are already in it\n", m.Author.Username)
				break
			}
		}

		players[*m.Author] = ""
		playerOrder = append(playerOrder, *m.Author)
		fmt.Printf("Added %s to the game\n", m.Author.Username)
	case "start":
		// Randomize the order
		playerOrder = []discordgo.User{}
		for player := range players {
			playerOrder = append(playerOrder, player)
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(playerOrder), func(i, j int) { playerOrder[i], playerOrder[j] = playerOrder[j], playerOrder[i] })

		fmt.Printf("Randomized player order to: %+v\n", playerOrder)

		// Let everyone know who they are picking a word for
		for _, player := range playerOrder {
			_, chooseFor := playerIsChoosingFor(player)

			channel, _ := s.UserChannelCreate(player.ID)
			s.ChannelMessageSend(channel.ID, fmt.Sprintf("Choose a character for: %s\n", chooseFor.Username))
			fmt.Printf("%s is choosing a character for %s\n", player.Username, chooseFor.Username)
		}

	default:
		// Find the player who sent the message
		// Identify who they were chosing for
		// record the string against that person

		_, chooseFor := playerIsChoosingFor(*m.Author)
		players[chooseFor] = m.Content
		fmt.Printf("%s recorded character for %s\n", m.Author.Username, chooseFor.Username)

		// Check if all players have now got a character
		for _, character := range players {
			if character == "" {
				return
			}
		}
		fmt.Println("All players now have a character!")

		// All players have got a character, so send them out
		for player := range players {
			msg := ""
			for guesser, character := range players {
				if player == guesser {
					continue
				}

				msg += fmt.Sprintf("%s: %s\n", guesser.Username, character)
			}

			channel, _ := s.UserChannelCreate(player.ID)
			s.ChannelMessageSend(channel.ID, msg)
		}
	}
}

func playerIsChoosingFor(p discordgo.User) (index int, user discordgo.User) {
	for index, player := range playerOrder {
		if player == p {
			if index == len(playerOrder)-1 {
				return 0, playerOrder[0]
			}
			return index + 1, playerOrder[index+1]
		}
	}

	fmt.Printf("error: could not find who %s is chosing for", p.Username)
	return -1, discordgo.User{}
}

// https://discordapp.com/api/oauth2/authorize?client_id=703671133014327347&scope=bot&permissions=67584
