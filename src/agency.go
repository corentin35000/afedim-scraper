package main

import (
	"fmt"
	"log"

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
	Afedim  Agency = "Afedim"
	Giboire Agency = "Giboire"
)

/**
 * setupMainPageAfedim configure le collecteur pour la page principale de l'agence Afedim.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
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
