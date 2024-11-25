package main

import (
	"fmt"
	"time"
)

/**
 * RunScraper lance le scraping des annonces immobilières à intervalles réguliers.
 * Cette fonction est appelée depuis le point d'entrée de l'application.
 * @param {int} intervalMinutes - Intervalle de temps en minutes entre chaque cycle de scraping
 * return {void}
 */
func RunScraper(intervalMinutes int) {
	// Map globale pour suivre les références des biens déjà traités par les différentes agences
	processedReferencesAfedim := make(map[string]bool)
	processedReferencesGiboire := make(map[string]bool)
	processedReferencesFoncia := make(map[string]bool)
	processedReferencesAgenceDuColombier := make(map[string]bool)

	for {
		// Lancer le scraping pour l'agence Afedim
		processAgencyScraping(processedReferencesAfedim, "https://www.afedim.fr/fr/location/annonces/Appartement-Maison-Parking-Garage/Rennes-France/1-5-pieces/surface-0-100-m2/budget-0-90000-euros/rayon-10-km/disponible-/options-/exclusPlafondRess-/Resultats", "AFEDIM", "Afedim")

		// Lancer le scraping pour l'agence Giboire
		processAgencyScraping(processedReferencesGiboire, "https://www.giboire.com/recherche-location/appartement/?searchBy=default&address%5B%5D=RENNES&address%5B%5D=CHANTEPIE&address%5B%5D=CESSON+SEVIGNE&priceMax=700&nbBedrooms%5B%5D=1&transactionType%5B%5D=Location&searchBy=default", "GIBOIRE", "Giboire")

		// Lancer le scraping pour l'agence Foncia
		processAgencyScraping(processedReferencesFoncia, "https://fr.foncia.com/location/rennes-35--chantepie-35135--cesson-sevigne-35510/appartement?nbPiece=2--&prix=--700&advanced=", "FONCIA", "Foncia")

		// Lancer le scraping pour l'agence Agence du Colombier
		processAgencyScraping(processedReferencesAgenceDuColombier, "https://agenceducolombier.com/annonces/?filter_search_action%5B%5D=louer&filter_search_type%5B%5D=&nb-pieces=&min-chambres=&min-surface=&max-surface=&price_low=0&price_max=6000000&submit=LANCER+MA+RECHERCHE", "AGENCE DU COLOMBIER", "Agence du Colombier")

		// Attendre avant le prochain cycle
		time.Sleep(time.Duration(intervalMinutes) * time.Minute)
	}
}

/**
 * processAgencyScraping lance le scraping pour une agence immobilière spécifique.
 * @param {map[string]bool} processedReferences - Map contenant les références des biens déjà traités.
 * @param {string} url - L'URL de la page de l'agence à scraper.
 * @param {string} titleMessageTelegram - Le titre du message Telegram.
 * @param {Agency} nameAgency - Le nom de l'agence.
 * @return {void}
 */
func processAgencyScraping(processedReferences map[string]bool, url string, titleMessageTelegram string, nameAgency Agency) {
	// Créer une nouvelle instance de CollyService
	collyService := NewCollyService()

	// Récupérer les annonces complètes depuis l'agence
	newAnnouncements := collyService.ScrapeAnnouncement(nameAgency, url)

	// Comparer les références des biens pour détecter les nouvelles annonces
	for _, announcement := range newAnnouncements {
		if _, exists := processedReferences[announcement.propertyReference]; !exists {
			// Nouvelle annonce détectée
			fmt.Println("Nouvelle annonce détectée référence :", announcement.propertyReference)
			processedReferences[announcement.propertyReference] = true

			// Envoie un message sur le canal Telegram
			sendTelegramMessageToPublicChannel(fmt.Sprintf(
				"%s\nNouvelle annonce immobilière !\nRéférence : %s\nURL : %s",
				titleMessageTelegram,
				announcement.propertyReference,
				announcement.url,
			))
		}
	}
}
