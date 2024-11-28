package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	// Ignorer les erreurs de certificat TLS
	collyService.collector.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

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
	case LaFrancaiseImmobiliere:
		setupMainPageLaFrancaiseImmobiliere(collyService.collector, &detailPageURLs)
	case Guenno:
		setupMainPageGuenno(collyService.collector, &detailPageURLs)
	case LaMotte:
		setupMainPageLaMotte(collyService.collector, &detailPageURLs)
	case Kermarrec:
		setupMainPageKermarrec(collyService.collector, &detailPageURLs)
	case Nestenn:
		setupMainPageNestenn(collyService.collector, &detailPageURLs)
	case SquareHabitat:
		setupMainPageSquareHabitat(collyService.collector, &detailPageURLs)
	case CAImmobilier:
		setupMainPageCAImmobilier(collyService.collector, &detailPageURLs)
	case PigeaultImmobilier:
		setupMainPagePigeaultImmobilier(collyService.collector, &detailPageURLs)
	case LaForetImmobilier:
		setupMainPageLaForetImmobilier(collyService.collector, &detailPageURLs)
	case Cogir:
		setupMainPageCogir(collyService.collector, &detailPageURLs)
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

	if agency == SquareHabitat {
		// La référence devient la description pour Square Habitat
		var announcements []Announcement
		for _, ref := range detailPageURLs {
			announcements = append(announcements, Announcement{propertyReference: ref, url: ""})
		}
		return announcements
	} else if agency == CAImmobilier {
		// La référence devient la description pour CA Immobilier
		var announcements []Announcement
		for _, ref := range detailPageURLs {
			// Extraire la partie après le dernier "/" pour le propertyReference
			lastSlashIndex := strings.LastIndex(ref, "/")
			var propertyReference string
			if lastSlashIndex != -1 && lastSlashIndex+1 < len(ref) {
				propertyReference = ref[lastSlashIndex+1:]
			} else {
				propertyReference = "unknown" // Valeur par défaut si la référence est mal formée
			}

			// Ajouter à la liste des annonces
			announcements = append(announcements, Announcement{
				propertyReference: propertyReference,
				url:               ref,
			})
		}
		return announcements
	}

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

	// Ignorer les erreurs de certificat TLS
	detailCollector.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

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
	case LaFrancaiseImmobiliere:
		processDetailPagesLaFrancaiseImmobiliere(detailCollector, &announcements)
	case Guenno:
		processDetailPagesGuenno(detailCollector, &announcements)
	case LaMotte:
		processDetailPagesLaMotte(detailCollector, &announcements)
	case Kermarrec:
		processDetailPagesKermarrec(detailCollector, &announcements)
	case Nestenn:
		processDetailPagesNestenn(detailCollector, &announcements)
	case PigeaultImmobilier:
		processDetailPagesPigeaultImmobilier(detailCollector, &announcements)
	case LaForetImmobilier:
		processDetailPagesLaForetImmobilier(detailCollector, &announcements)
	case Cogir:
		processDetailPagesCogir(detailCollector, &announcements)
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
