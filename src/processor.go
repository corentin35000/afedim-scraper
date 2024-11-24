package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TelegramBotToken = "7423967574:AAHUuvNAsvLsTQ6bHMxHuOWxws_LXVeQUHw" // Token du bot Telegram
	TelegramChannel  = "@annonceimmobilierafidim"                       // Nom d'utilisateur public du canal
)

/**
 * RunScraper initialise les services NATS et Colly, et lance le scraping en continu.
 * Cette fonction est appelée depuis le point d'entrée de l'application.
 * @param {int} intervalMinutes - Intervalle de temps en minutes entre chaque cycle de scraping
 * return {void}
 */
func RunScraper(intervalMinutes int) {
	// Map globale pour suivre les références des biens déjà traités
	// Tableau d'annonces immobilières déjà traitées
	processedReferences := make(map[string]bool)

	for {
		// Créer une nouvelle instance de CollyService
		collyService := NewCollyService()

		// Récupérer les références des annonces
		newReferences := collyService.ScrapeAnnouncement("https://www.afedim.fr/fr/location/annonces/Appartement-Maison-Parking-Garage/Rennes-France/1-5-pieces/surface-0-100-m2/budget-0-90000-euros/rayon-10-km/disponible-/options-/exclusPlafondRess-/Resultats")

		// Vérifier les nouvelles annonces
		for _, ref := range newReferences {
			if _, exists := processedReferences[ref]; !exists {
				// Nouvelle annonce détectée
				fmt.Println("Nouvelle annonce détectée :", ref)
				processedReferences[ref] = true

				// Envoyer une notification Telegram (fonction fictive pour l'exemple)
				sendTelegramMessageToPublicChannel(fmt.Sprintf("Nouvelle annonce référence : %s", ref))
			}
		}

		// Attendre avant le prochain cycle
		time.Sleep(time.Duration(intervalMinutes) * time.Minute)
	}
}

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
		log.Println("Message envoyé au canal Telegram :", message)
	}
}
