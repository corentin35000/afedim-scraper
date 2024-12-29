package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TelegramBotToken = "7423967574:AAHUuvNAsvLsTQ6bHMxHuOWxws_LXVeQUHw" // Token du bot Telegram (@AnnoncesImmobiliersScraperBot)
	TelegramChannel  = "@annonceimmobiliers"                            // Canal Telegram public (https://t.me/annonceimmobiliers)
	MaxRetries       = 5                                                // Nombre maximal de tentatives
)

// sendTelegramMessageToPublicChannel envoie un message à un canal Telegram public.
func sendTelegramMessageToPublicChannel(message string) {
	// Initialiser le bot Telegram
	bot, err := tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		log.Fatalf("Erreur lors de la création du bot Telegram : %v", err)
	}

	// Créer un nouveau message pour le canal
	msg := tgbotapi.NewMessageToChannel(TelegramChannel, message)

	retries := 0

	for {
		// Envoyer le message
		_, err = bot.Send(msg)
		if err != nil {
			// Vérifier si l'erreur est liée aux limites de débit
			if apiErr, ok := err.(*tgbotapi.Error); ok && apiErr.RetryAfter > 0 {
				log.Printf("Trop de requêtes pour l'API Telegram. Réessayer après %d secondes. Scraper mis en pause en attendant", apiErr.RetryAfter)
				time.Sleep(time.Duration(apiErr.RetryAfter) * time.Second)
			} else {
				log.Printf("Erreur lors de l'envoi du message Telegram : %v", err)
				return
			}
		} else {
			log.Println("Message envoyé au canal Telegram.")
			return
		}

		retries++
		if retries >= MaxRetries {
			log.Println("Nombre maximal de tentatives atteint. Abandon de l'envoi.")
			return
		}
	}
}
