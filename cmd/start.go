package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Brotchu/bot/db"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(startCmd)
}

var selfID string //botid
// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the bot to listen to messages",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("start called")
		home, _ := homedir.Dir()
		dbPath := filepath.Join(home, "botid.db")

		must(db.Init(dbPath))

		botToken, err := db.GetBotID("botid")
		must(err)

		dg, err := discordgo.New("Bot " + botToken)
		must(err)

		u, err := dg.User("@me")
		must(err)

		selfID = u.ID
		// fmt.Println(SelfID)
		dg.AddHandler(messageHandlerFunc)
		must(dg.Open())
		printdash()

		go handleMessage(dg)
		<-make(chan struct{})
	},
}

func messageHandlerFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	others := color.New(color.FgWhite, color.BgBlue).Render
	sep := color.New(color.FgWhite, color.BgRed).Render
	selfc := color.New(color.FgBlack, color.BgYellow).Render
	chs, err := db.GetChannels()
	must(err)
	cname, err := channelIdtoName(m.ChannelID, chs)
	if err != nil {
		cname = m.ChannelID
	}
	// fmt.Printf("DEBUG %s\n", m.ChannelID)
	if m.Author.ID == selfID {
		fmt.Printf("%s%s%s%s %s\n", selfc(m.Author), sep("#"), selfc(cname), selfc(" >>"), m.Content)
	} else {
		fmt.Printf("%s%s%s%s %s\n", others(m.Author), sep("#"), others(cname), others(" >>"), m.Content)
	}

}

func handleMessage(s *discordgo.Session) {
	chs, err := db.GetChannels()
	must(err)

	defaultch, err := getdefault(chs, "")
	warn(err)
	fmt.Printf("/addc  <channel name> to add new channel \n")
	fmt.Printf("/listc to list channels\n")
	fmt.Printf("/default to set default channel\n")
	fmt.Printf("/send <channel name> <message> to send to specific channel\n")
	fmt.Printf("/c to exit \n")
	fmt.Printf("type in messages to send to default channel\n")
	printdash()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		must(err)

		text = text[:len(text)-1]
		words := strings.Split(text, " ")

		if words[0] == "/addc" {
			if len(words) > 2 {
				warn(db.AddChannel(words[1], words[2]))
				chs, err = db.GetChannels()
				listchannels(chs)
				must(err)
			} else {
				fmt.Println("Specify channel name and id")
			}
		} else if words[0] == "/default" {
			if len(words) > 1 {
				defaultch, err = getdefault(chs, words[1])
				warn(err)
				if err == nil {
					fmt.Printf("default channel set - %s\n", defaultch)
					printdash()
				}
			} else {
				fmt.Printf("default Channel : %s\n", defaultch)
				fmt.Println("Specify channel name to set default [/default <channel name>]")
				printdash()
			}
		} else if words[0] == "/send" {
			if len(words) > 2 {
				s.ChannelMessageSend(chs[words[1]], strings.Join(words[2:], " "))
			} else {
				fmt.Printf("specify channel name and message to send to specific channel - /send <channel name> <msg> \n")
				printdash()
			}
		} else if words[0] == "/listc" {
			listchannels(chs)
		} else if words[0] == "/c" {
			os.Exit(0)
		} else {
			if dc, ok := chs[defaultch]; ok {
				s.ChannelMessageSend(dc, strings.Join(words, " "))
			} else {
				fmt.Println("default channel not set. use /default <channel name> to set.")
				printdash()
			}
		}
	}
}

func getdefault(cm map[string]string, cname string) (string, error) {
	var keys []string
	for k := range cm {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return "", errors.New("no channels, Use /addc to add a channel")
	}
	if cname == "" {
		return keys[0], nil
	}
	if _, ok := cm[cname]; ok {
		return cname, nil
	} else {
		return "", errors.New("channel not found")
	}
}

func listchannels(cm map[string]string) {
	for k, v := range cm {
		fmt.Printf("%s  ==> %s\n", k, v)
	}
	printdash()
}

func channelIdtoName(id string, cm map[string]string) (string, error) {
	for k, v := range cm {
		if v == id {
			return k, nil
		}
	}
	return "", errors.New("id not found")
}

func printdash() {
	fmt.Printf("--------------------------------------------------------------\n")
}
func warn(err error) {
	if err != nil {
		fmt.Println("Something went wrong : ", err.Error())
	}
}
