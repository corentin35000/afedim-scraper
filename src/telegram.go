package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TelegramBotToken = "7423967574:AAHUuvNAsvLsTQ6bHMxHuOWxws_LXVeQUHw" // Token du bot Telegram (@AnnonceImmobilierScraperBot)
	TelegramChannel  = "@annonceimmobilierafidim"                       // Nom d'utilisateur public du canal public (https://t.me/annonceimmobiliers)
)

/**
 * sendTelegramMessageToPublicChannel envoie un message à un canal Telegram public.
 * @param {string} message - Le message à envoyer.
 * @return {void}
 */
func sendTelegramMessageToPublicChannel(message string) {
	// Initialiser le bot Telegram
	bot, err := tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		log.Fatalf("Erreur lors de la création du bot Telegram : %v", err)
	}

	// Créer un nouveau message pour le canal
	msg := tgbotapi.NewMessageToChannel(TelegramChannel, message)

	// Envoyer le message
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Erreur lors de l'envoi du message Telegram : %v", err)
	} else {
		log.Println("Message envoyé au canal Telegram.")
	}
}
