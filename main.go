package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/mb-14/gomarkov"
)

var (
	chain  *gomarkov.Chain
	botCfg config
)

func init() {
	log.Println("Loading and parsing config...")
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("[error] failed to open config.json: %v\n", err)
	}
	cfgParser := json.NewDecoder(cfgFile)
	err = cfgParser.Decode(&botCfg)
	if err != nil {
		log.Fatalf("[error] failed to open config.json: %v\n", err)
	}
}

func main() {
	fmt.Println("Markov-Bot")

	_, err := os.Stat("model.json")
	if err == nil {
		log.Println("Loading model from model.json...")
		data, err := ioutil.ReadFile("model.json")
		if err != nil {
			log.Fatalf("[error] failed to load model: %v\n", err)
		}
		err = json.Unmarshal(data, &chain)
		if err != nil {
			log.Fatalf("[error] failed to parse model.json: %v\n", err)
		}

	} else {
		log.Println("[info] creating new chain...")
		chain = gomarkov.NewChain(1)
	}

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", botCfg.Token))
	if err != nil {
		log.Fatalf("[error] failed to create Discord session: %v\n", err)
	}

	// Handlers go here
	dg.AddHandler(messageCreate)
	dg.AddHandler(discordReady)

	err = dg.Open()
	if err != nil {
		log.Fatalf("[error] failed to open Discord session: %v\n", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	saveModel()

}

func saveModel() {
	obj, _ := json.Marshal(chain)

	err := ioutil.WriteFile("model.json", obj, 0644)
	if err != nil {
		log.Printf("[error] failed to write model to file: %v\n", err)
	}
}
