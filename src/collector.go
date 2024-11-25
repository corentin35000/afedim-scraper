package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

/**
 * Announcement est une structure pour stocker les informations sur les annonces de bien immobilier.
 * @property {string} propertyReference - Référence du bien immobilier.
 * @property {string} url - URL de la page de détails de l'annonce.
 */
type Announcement struct {
	propertyReference string
	url               string
}

/**
 * ScrapeAnnouncement lance le scraping des annonces immobilières à partir de la page spécifiée.
 * @param {string} url - L'URL de la page à scraper.
 * @return {[]Announcement} - Slice contenant les annonces.
 */
func (collyService *CollyService) ScrapeAnnouncement(agency Agency, url string) []Announcement {
	// Slice pour stocker les URLs des pages de détails
	var detailPageURLs []string

	// Afficher un message de démarrage
	fmt.Println("Démarrage du scraping des annonces immobilières de l'agence :", agency)

	// Utiliser un switch pour configurer les callbacks spécifiques à l'agence
	switch agency {
	case Afedim:
		setupMainPageAfedim(collyService.collector, &detailPageURLs)
	case Giboire:
		setupMainPageGiboire(collyService.collector, &detailPageURLs)
	case Foncia:
		setupMainPageFoncia(collyService.collector, &detailPageURLs)
	case AgenceDuColombier:
		setupMainPageAgenceDuColombier(collyService.collector, &detailPageURLs)
	default:
		log.Fatalf("Agence inconnue : %s", agency)
	}

	// Gestion des erreurs pour la page principale
	collyService.collector.OnError(func(_ *colly.Response, err error) {
		log.Printf("Erreur pendant le scraping de la page principale : %v", err)
	})

	// Démarrer le scraping de la page principale
	if err := collyService.collector.Visit(url); err != nil {
		log.Printf("Erreur lors de la visite de l'URL principale : %v", err)
	}

	// Attendre la fin des requêtes asynchrones
	collyService.collector.Wait()

	// Récupérer les annonces complètes (références et URLs)
	return collyService.processDetailPages(detailPageURLs, agency)
}

/**
 * processDetailPages traite les pages de détails des annonces immobilières.
 * @param {[]string} detailPageURLs - Slice contenant les URLs des pages de détails.
 * @return {[]Announcement} - Slice contenant les annonces.
 */
func (collyService *CollyService) processDetailPages(detailPageURLs []string, agency Agency) []Announcement {
	// Slice pour stocker les annonces
	var announcements []Announcement

	// Créer un nouveau collector pour les pages de détails
	detailCollector := colly.NewCollector()

	// Utiliser un switch pour appeler la fonction spécifique à l'agence
	switch agency {
	case Afedim:
		processDetailPagesAfedim(detailCollector, &announcements)
	case Giboire:
		processDetailPagesGiboire(detailCollector, &announcements)
	case Foncia:
		processDetailPagesFoncia(detailCollector, &announcements)
	case AgenceDuColombier:
		processDetailPagesAgenceDuColombier(detailCollector, &announcements)
	default:
		log.Fatalf("Agence inconnue : %s", agency)
	}

	// Gestion des erreurs pour les détails
	detailCollector.OnError(func(_ *colly.Response, err error) {
		log.Printf("Erreur pendant le scraping de la page de détails : %v", err)
	})

	// Visiter chaque URL dans la slice
	for _, url := range detailPageURLs {
		fmt.Println("Visite de la page de détails :", url)
		if err := detailCollector.Visit(url); err != nil {
			log.Printf("Erreur lors de la visite de la page de détails : %v", err)
		}
	}

	// Attendre la fin des requêtes asynchrones
	detailCollector.Wait()

	// Retourner toutes les annonces trouvées
	return announcements
}
