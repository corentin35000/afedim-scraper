package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

/**
 * CollyService est une structure qui encapsule le collecteur Colly pour le scraping de données.
 * @property {colly.Collector} collector - Instance du collecteur Colly pour le scraping.
 * @property {chan error} errChan - Canal pour signaler les erreurs pendant le scraping.
 */
type CollyService struct {
	collector *colly.Collector
	errChan   chan error // Canal pour signaler les erreurs
}

/**
 * Announcement est une structure pour stocker les informations sur les annonces de bien immobilier.
 * @property {string} propertyReference - Référence du bien immobilier.
 */
type Announcement struct {
	propertyReference string
}

/**
 * NewCollyService crée une nouvelle instance de CollyService avec une configuration de collecteur prédéfinie.
 * @return {CollyService} - Retourne une instance configurée de CollyService.
 */
func NewCollyService() *CollyService {
	// Configuration de base du collecteur
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"), // User-Agent pour éviter le blocage
		colly.IgnoreRobotsTxt(), // Ignorer les règles du fichier robots.txt
		colly.MaxDepth(1),       // Limiter la profondeur de recherche à 1 pour éviter les liens externes
		colly.Async(true),       // Activer le mode asynchrone pour le scraping
		colly.CacheDir("./tmp"), // Définir le répertoire de cache pour éviter de re-scraper les pages
		colly.DetectCharset(),   // Détecter automatiquement l'encodage de la page
	)

	// Définition des limites de requêtes pour éviter les blocages
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",             // Applique cette règle à tous les domaines visités par le collecteur
		Parallelism: 2,               // Limite à 2 requêtes simultanées pour ne pas surcharger le serveur
		Delay:       2 * time.Second, // Attente de 2 secondes entre chaque requête pour éviter un blocage par le serveur
	})

	// Retourne une nouvelle instance de CollyService avec un canal d'erreur
	return &CollyService{
		collector: c,
		errChan:   make(chan error), // Initialiser le canal d'erreurs
	}
}

/**
 * ScrapeAnnouncement lance le scraping des annonces immobilières à partir de la page spécifiée.
 * @param {string} url - L'URL de la page à scraper.
 * @return {[]string} - Slice contenant les références des biens.
 */
func (collyService *CollyService) ScrapeAnnouncement(url string) []string {
	// Slice pour stocker les URLs des pages de détails
	var detailPageURLs []string

	// Afficher un message de démarrage
	fmt.Println("Démarrage du scraping des annonces immobilières depuis :", url)

	// Callback pour scraper les annonces sur la page principale
	collyService.collector.OnHTML("#C\\:blocRecherche\\.blocRechercheDesk\\.P\\.C\\:U", func(e *colly.HTMLElement) {
		e.ForEach("li.item", func(_ int, li *colly.HTMLElement) {
			li.ForEach("div div div:last-child span a", func(_ int, el *colly.HTMLElement) {
				detailPageURL := el.Attr("href")

				if detailPageURL != "" {
					detailPageURL = "https://www.afedim.fr" + detailPageURL
					detailPageURLs = append(detailPageURLs, detailPageURL)
				}
			})
		})
	})

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

	// Récupérer les références des pages de détails
	return collyService.processDetailPages(detailPageURLs)
}

/**
 * processDetailPages traite les pages de détails des annonces immobilières.
 * @param {[]string} detailPageURLs - Slice contenant les URLs des pages de détails.
 * @return {[]string} - Slice contenant les références des biens.
 */
func (collyService *CollyService) processDetailPages(detailPageURLs []string) []string {
	// Slice pour stocker les références scrappées
	var scrapedReferences []string

	// Créer un nouveau collector pour les pages de détails
	detailCollector := colly.NewCollector()

	// Callback pour scraper les détails
	detailCollector.OnHTML("span[class*='note']", func(detail *colly.HTMLElement) {
		fullValue := detail.Text

		// Extraire la référence si le format correspond
		var reference string
		fmt.Sscanf(fullValue, "Référence du bien : %s", &reference)

		// Ajouter la référence à la slice si elle est valide
		if reference != "" {
			fmt.Println("Référence du bien :", reference)
			scrapedReferences = append(scrapedReferences, reference)
		}
	})

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

	// Retourner toutes les références trouvées
	return scrapedReferences
}

/**
 * ErrorChannel retourne le canal d'erreurs pour le scraping.
 * @return {<-chan error} - Canal en lecture seule pour les erreurs de scraping.
 */
func (collyService *CollyService) ErrorChannel() <-chan error {
	return collyService.errChan
}
