package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

/**
 * Agency est un type énuméré pour les agences immobilières.
 */
type Agency string

/**
 * Constantes pour les agences immobilières.
 */
const (
	Afedim            Agency = "Afedim"
	Giboire           Agency = "Giboire"
	Foncia            Agency = "Foncia"
	AgenceDuColombier Agency = "Agence du Colombier"
)

/**
 * setupMainPageAfedim configure le collecteur pour la page principale de l'agence Afedim.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageAfedim(collector *colly.Collector, detailPageURLs *[]string) {
	collector.OnHTML("#C\\:blocRecherche\\.blocRechercheDesk\\.P\\.C\\:U", func(e *colly.HTMLElement) {
		e.ForEach("li.item", func(_ int, li *colly.HTMLElement) {
			li.ForEach("div div div:last-child span a", func(_ int, el *colly.HTMLElement) {
				detailPageURL := el.Attr("href")
				if detailPageURL != "" {
					*detailPageURLs = append(*detailPageURLs, "https://www.afedim.fr"+detailPageURL)
				}
			})
		})
	})
}

/**
 * processDetailPagesAfedim extrait les références des annonces de la page de détail de l'agence Afedim.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesAfedim(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("span[class*='note']", func(detail *colly.HTMLElement) {
		fullValue := detail.Text

		var reference string
		fmt.Sscanf(fullValue, "Référence du bien : %s", &reference)

		if reference != "" {
			url := detail.Request.URL.String()
			*announcements = append(*announcements, Announcement{
				propertyReference: reference,
				url:               url,
			})
		}
	})
}

/**
 * setupMainPageGiboire configure le collecteur pour la page principale de l'agence Giboire.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageGiboire(collector *colly.Collector, detailPageURLs *[]string) {
	collector.OnHTML(".result-grid_wrap", func(e *colly.HTMLElement) {
		// Parcourir chaque div représentant une annonce
		e.ForEach("div", func(_ int, div *colly.HTMLElement) {
			// Chercher l'article à l'intérieur de chaque div
			article := div.DOM.Find("article")
			if article.Length() > 0 {
				// Récupérer la deuxième div dans l'article
				secondDiv := article.Find("div:nth-child(2)")
				if secondDiv.Length() > 0 {
					// Trouver la balise <h2> contenant le lien <a>
					h2 := secondDiv.Find("h2 a")
					href, exists := h2.Attr("href")
					if exists && href != "" {
						// Ajouter le lien complet à la liste des URLs
						*detailPageURLs = append(*detailPageURLs, href)
					}
				}
			}
		})
	})
}

/**
 * processDetailPagesGiboire extrait les références des annonces de la page de détail de l'agence Giboire.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesGiboire(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("p.presentation-bien_exclu_desc_ref", func(detail *colly.HTMLElement) {
		// Récupérer le texte brut dans la balise
		fullValue := detail.Text

		// Nettoyer le texte pour extraire uniquement la référence
		var reference string
		if _, err := fmt.Sscanf(fullValue, "Réf : %s", &reference); err == nil {
			if reference != "" {
				// URL de la page actuelle
				url := detail.Request.URL.String()

				// Ajouter l'annonce à la liste
				*announcements = append(*announcements, Announcement{
					propertyReference: reference,
					url:               url,
				})
			}
		} else {
			log.Printf("Impossible d'extraire la référence depuis : %s", fullValue)
		}
	})
}

/**
 * setupMainPageFoncia configure le collecteur pour la page principale de l'agence Foncia.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageFoncia(collector *colly.Collector, detailPageURLs *[]string) {
	// Utiliser un ensemble pour éviter les doublons
	seenURLs := make(map[string]struct{})

	// Cibler la div contenant toutes les annonces
	collector.OnHTML("div.p-col-12.mosaic-list.large.ng-star-inserted", func(e *colly.HTMLElement) {
		// Itérer sur chaque div enfant représentant une annonce
		e.ForEach("div", func(_ int, annonce *colly.HTMLElement) {
			// Cibler la deuxième div dans chaque annonce
			annonce.ForEach("div:nth-child(2)", func(_ int, secondDiv *colly.HTMLElement) {
				// Trouver la balise <a> et extraire l'attribut href
				href := secondDiv.ChildAttr("a", "href")
				if href != "" {
					// Construire l'URL complète si nécessaire
					fullURL := "https://fr.foncia.com" + href

					// Vérifier si l'URL est déjà dans l'ensemble
					if _, exists := seenURLs[fullURL]; !exists {
						// Ajouter à l'ensemble et à la liste
						seenURLs[fullURL] = struct{}{}
						*detailPageURLs = append(*detailPageURLs, fullURL)
					}
				}
			})
		})
	})
}

/**
 * processDetailPagesFoncia extrait les références des annonces de la page de détail de l'agence Foncia.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesFoncia(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("p.section-reference", func(detail *colly.HTMLElement) {
		// Récupérer le texte brut dans la balise
		fullValue := strings.TrimSpace(detail.Text) // Nettoyage de la chaîne

		// Essayer de parser la référence
		var reference string
		if _, err := fmt.Sscanf(fullValue, "Réf. %s", &reference); err == nil {
			if reference != "" {
				// URL de la page actuelle
				url := detail.Request.URL.String()

				// Ajouter l'annonce à la liste
				*announcements = append(*announcements, Announcement{
					propertyReference: reference,
					url:               url,
				})
			} else {
				log.Printf("Référence vide après extraction depuis : %s", fullValue)
			}
		} else {
			log.Printf("Erreur lors de l'extraction de la référence depuis : %s, erreur : %v", fullValue, err)
		}
	})
}

/**
 * setupMainPageAgenceDuColombier configure le collecteur pour la page principale de l'agence Agence du Colombier.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageAgenceDuColombier(collector *colly.Collector, detailPageURLs *[]string) {
	collector.OnHTML("div#listing_ajax_container", func(e *colly.HTMLElement) {
		log.Println("Entrée dans le conteneur principal des annonces")

		// Compter le nombre d'annonces et extraire les URLs
		e.ForEach("div.listing_wrapper", func(i int, annonce *colly.HTMLElement) {
			log.Printf("Annonce trouvée : %d", i+1)
			detailURL := annonce.ChildAttr("a", "href")
			if detailURL != "" {
				*detailPageURLs = append(*detailPageURLs, detailURL)
				log.Printf("URL de la page de détail : %s", detailURL)
			}
		})

		log.Printf("Nombre total d'annonces : %d", len(*detailPageURLs))
	})

}

/**
 * processDetailPagesAgenceDuColombier extrait les références des annonces de la page de détail de l'agence Agence du Colombier.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesAgenceDuColombier(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler la div contenant les informations principales, notamment la référence
	collector.OnHTML("div.wpestate_estate_property_design_intext_details", func(detail *colly.HTMLElement) {
		// Trouver la balise <p> contenant "REF:"
		detail.ForEach("p", func(_ int, el *colly.HTMLElement) {
			// Vérifier si la balise contient "REF:"
			if strings.Contains(el.Text, "REF:") {
				// Extraire le texte brut et isoler la référence
				fullText := strings.TrimSpace(el.Text)
				var reference string

				// Extraire la partie après "REF:"
				if _, err := fmt.Sscanf(fullText, "REF: %s", &reference); err == nil {
					if reference != "" {
						// URL de la page actuelle
						url := detail.Request.URL.String()

						// Ajouter l'annonce à la liste des résultats
						*announcements = append(*announcements, Announcement{
							propertyReference: reference,
							url:               url,
						})
					}
				} else {
					log.Printf("Impossible d'extraire la référence depuis : %s", fullText)
				}
			}
		})
	})
}
