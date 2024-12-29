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
	errChan   chan error
}

// Liste des User-Agents pour éviter le blocage
/*var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3328.95 Safari/537.36 QIHU 360SE",
	"Mozilla/5.0 (Linux; Android 14; SM-M136B Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.106 Mobile Safari/537.36 WebView MetaMaskMobile",
	"Mozilla/5.0 (Linux; U; Android 14; en-gb; RMX3612 Build/UKQ1.230924.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.88 Mobile Safari/537.36 HeyTapBrowser/45.11.4.1",
	"Mozilla/5.0 (Linux; Android 13; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.5563.116 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; RMX3690 Build/SP1A.210812.016) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6723.107 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11.12; rv:131.0) Gecko/20010101 Firefox/131.0",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.6556.192 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.2; Win64; x64; rv:132.0esr) Gecko/20100101 Firefox/132.0esr/0YoBqLP7z7eKob-09",
	"Mozilla/5.0 (Linux; Android 12; NCO-LX1; HMSCore 6.14.0.322) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.196 HuaweiBrowser/15.0.4.312 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 14; SM-S918U1 Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.106 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 17_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/128.0.6613.92 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 YaBrowser/22.11.2.803 Yowser/2.5 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 YaBrowser/22.11.5.711 Yowser/2.5 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11.12; rv:131.0) Gecko/20010101 Firefox/131.0",
	"Mozilla/5.0 (Linux; Android 14; SM-M136B Build/UP1A.231005.007; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/130.0.6723.106 Mobile Safari/537.36 WebView MetaMaskMobile",
}

// Générateur de nombres aléatoires
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

/**
 * randomUserAgent retourne un User-Agent aléatoire pour éviter le blocage.
 * @return {string} - User-Agent aléatoire.
*/
/*func randomUserAgent() string {
	return userAgents[random.Intn(len(userAgents))]
}*/

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
