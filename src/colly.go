package main

import (
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
 * ErrorChannel retourne le canal d'erreurs pour le scraping.
 * @return {<-chan error} - Canal en lecture seule pour les erreurs de scraping.
 */
func (collyService *CollyService) ErrorChannel() <-chan error {
	return collyService.errChan
}
