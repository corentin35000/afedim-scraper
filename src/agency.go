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
	Afedim                 Agency = "Afedim"
	Giboire                Agency = "Giboire"
	Foncia                 Agency = "Foncia"
	AgenceDuColombier      Agency = "Agence du Colombier"
	LaFrancaiseImmobiliere Agency = "La Française Immobilière"
	Guenno                 Agency = "Guenno"
	LaMotte                Agency = "La Motte"
	Kermarrec              Agency = "Kermarrec"
	Nestenn                Agency = "Nestenn"
	SquareHabitat          Agency = "Square Habitat"
	CAImmobilier           Agency = "CA Immobilier"
	PigeaultImmobilier     Agency = "Pigeault Immobilier"
	LaForetImmobilier      Agency = "La Foret Immobilier"
	Cogir                  Agency = "Cogir"
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
		// Compter le nombre d'annonces et extraire les URLs
		e.ForEach("div.listing_wrapper", func(i int, annonce *colly.HTMLElement) {
			detailURL := annonce.ChildAttr("a", "href")
			if detailURL != "" {
				*detailPageURLs = append(*detailPageURLs, detailURL)
			}
		})
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

/**
 * setupMainPageLaFrancaiseImmobiliere configure le collecteur pour la page principale de l'agence La Française Immobilière.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageLaFrancaiseImmobiliere(collector *colly.Collector, detailPageURLs *[]string) {
	collector.OnHTML("div#liste_annonces div.row > article", func(e *colly.HTMLElement) {
		// Tente de récupérer l'attribut href du premier <a> dans chaque article
		detailLink := e.ChildAttr("a[rel='bookmark']", "href")
		if detailLink != "" {
			*detailPageURLs = append(*detailPageURLs, detailLink)
		} else {
			log.Println("Lien de détail introuvable dans cet article.")
		}
	})
}

/**
 * processDetailPagesLaFrancaiseImmobiliere extrait les références des annonces de la page de détail de l'agence La Française Immobilière.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesLaFrancaiseImmobiliere(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("p.ref.d-inline", func(detail *colly.HTMLElement) {
		// Récupérer le texte brut dans la balise
		fullValue := strings.TrimSpace(detail.Text) // Nettoyage de la chaîne

		// Essayer de parser la référence
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
			} else {
				log.Printf("Référence vide après extraction depuis : %s", fullValue)
			}
		} else {
			log.Printf("Erreur lors de l'extraction de la référence depuis : %s, erreur : %v", fullValue, err)
		}
	})
}

/**
 * setupMainPageGuenno configure le collecteur pour la page principale de l'agence Guenno.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageGuenno(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div contenant les annonces
	collector.OnHTML("div.section-content", func(e *colly.HTMLElement) {
		// Parcourir chaque balise <article> dans la section
		e.ForEach("article", func(_ int, article *colly.HTMLElement) {
			// Récupérer la valeur du href de la balise <a>
			href := article.ChildAttr("a", "href")

			// Vérifier si le lien est valide
			if href != "" {
				*detailPageURLs = append(*detailPageURLs, href)
			} else {
				log.Println("Aucun lien trouvé dans cet article.")
			}
		})
	})
}

/**
 * processDetailPagesGuenno extrait les références des annonces de la page de détail de l'agence Guenno.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesGuenno(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("div#realty_area.realty_details", func(detail *colly.HTMLElement) {
		// Rechercher l'élément contenant la référence dans l'attribut itemprop et span.grey-ref
		fullValue := detail.ChildText("span.grey-ref")
		fullValue = strings.TrimSpace(fullValue)

		// Vérification et extraction de la référence
		if strings.HasPrefix(fullValue, "Ref :") {
			// Extraire uniquement la partie après "Ref :"
			reference := strings.TrimSpace(strings.TrimPrefix(fullValue, "Ref :"))
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
			log.Printf("Impossible de trouver la référence dans : %s", fullValue)
		}
	})
}

/**
 * setupMainPageLaMotte configure le collecteur pour la page principale de l'agence La Motte.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageLaMotte(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div contenant les annonces
	collector.OnHTML("div.col-12.pr-md-0.col__list", func(e *colly.HTMLElement) {
		// Cibler la div avec id "result"
		e.ForEach("div#result", func(_ int, resultDiv *colly.HTMLElement) {
			// Parcourir chaque annonce
			resultDiv.ForEach("div.bien__wrapper--annonce", func(i int, annonce *colly.HTMLElement) {
				// Récupérer le lien dans <a>
				href := annonce.ChildAttr("a", "href")

				// Vérifier si le lien est valide
				if href != "" {
					*detailPageURLs = append(*detailPageURLs, href)
				} else {
					log.Printf("Annonce %d : Aucun lien trouvé", i+1)
				}
			})
		})
	})
}

/**
 * processDetailPagesLaMotte extrait les références des annonces de la page de détail de l'agence La Motte.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesLaMotte(collector *colly.Collector, announcements *[]Announcement) {
	collector.OnHTML("div.heading__delivery", func(detail *colly.HTMLElement) {
		// Récupérer le texte de la balise <p class="tva">
		fullValue := detail.ChildText("p.tva")
		fullValue = strings.TrimSpace(fullValue)

		// Vérifier si la valeur commence par "Lot"
		if strings.HasPrefix(fullValue, "Lot") {
			// Extraire la partie après "Lot"
			lot := strings.TrimSpace(strings.TrimPrefix(fullValue, "Lot"))
			if lot != "" {
				// URL de la page actuelle
				url := detail.Request.URL.String()

				// Ajouter l'annonce à la liste
				*announcements = append(*announcements, Announcement{
					propertyReference: lot,
					url:               url,
				})
			} else {
				log.Printf("Lot vide après extraction depuis : %s", fullValue)
			}
		} else {
			log.Printf("Impossible de trouver le lot dans : %s", fullValue)
		}
	})
}

/**
 * setupMainPageKermarrec configure le collecteur pour la page principale de l'agence Kermarrec.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageKermarrec(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div principale contenant les annonces
	collector.OnHTML("div#primary.content-area.listofposts.grid", func(e *colly.HTMLElement) {
		// Parcourir chaque balise <article> dans la div principale
		e.ForEach("article", func(_ int, article *colly.HTMLElement) {
			// Accéder à la div.panel > div.entry-content
			href := article.ChildAttr("div.panel div.entry-content a", "href")

			// Vérifier si le lien est valide
			if href != "" {
				*detailPageURLs = append(*detailPageURLs, href)
			} else {
				log.Println("Aucun lien trouvé dans cet article.")
			}
		})
	})
}

/**
 * processDetailPagesKermarrec extrait les références des annonces de la page de détail de l'agence Kermarrec.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesKermarrec(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler l'en-tête contenant la référence
	collector.OnHTML("header.container.entry-header", func(detail *colly.HTMLElement) {
		// Récupérer le texte contenant le ref dans la balise span.ref
		fullValue := detail.ChildText("span.ref")
		fullValue = strings.TrimSpace(fullValue)

		// Vérification et extraction de la référence
		if strings.HasPrefix(fullValue, "(ref :") {
			// Extraire uniquement la partie après "(ref :" et enlever la parenthèse fermante
			reference := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(fullValue, "(ref :"), ")"))
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
			log.Printf("Impossible de trouver la référence dans : %s", fullValue)
		}
	})
}

/**
 * setupMainPageNestenn configure le collecteur pour la page principale de l'agence Nestenn.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageNestenn(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div contenant les annonces
	collector.OnHTML("div#gridPropertyOnlyWidening", func(e *colly.HTMLElement) {
		// Parcourir chaque div contenant une annonce
		e.ForEach("div.relative.grid_map_container", func(_ int, property *colly.HTMLElement) {
			// Récupérer la valeur du href de la balise <a>
			href := property.ChildAttr("a", "href")

			// Vérifier si le lien est valide
			if href != "" {
				*detailPageURLs = append(*detailPageURLs, href)
			} else {
				log.Println("Aucun lien trouvé dans cette annonce.")
			}
		})
	})
}

/**
 * processDetailPagesNestenn extrait les références des annonces de la page de détail de l'agence Nestenn.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesNestenn(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler la div contenant la référence
	collector.OnHTML("div.property_ref", func(detail *colly.HTMLElement) {
		// Récupérer le texte brut dans la div
		fullValue := strings.TrimSpace(detail.Text)

		// Rechercher et extraire la référence après "Réf :"
		if strings.Contains(fullValue, "Réf :") {
			// Diviser la chaîne sur "Réf :" et récupérer la partie après
			parts := strings.Split(fullValue, "Réf :")
			if len(parts) > 1 {
				reference := strings.TrimSpace(parts[1])

				// Vérifier si la référence est non vide
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
				log.Printf("Impossible de diviser la chaîne pour trouver la référence : %s", fullValue)
			}
		} else {
			log.Printf("Pas de 'Réf :' trouvé dans : %s", fullValue)
		}
	})
}

/**
 * setupMainPageSquareHabitat configure le collecteur pour la page principale de Square Habitat.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageSquareHabitat(collector *colly.Collector, details *[]string) {
	// Cibler la div principale contenant les annonces
	collector.OnHTML("div.biens-container.afc-display-xs-flex.afc-width-xs-100", func(e *colly.HTMLElement) {
		// Parcourir chaque carte (card-container)
		e.ForEach("div.card-container", func(index int, property *colly.HTMLElement) {
			description := property.ChildText("app-card-bien > msl-card > div:nth-of-type(2) > div:nth-of-type(4) > app-texte-on-off > div.container > div.text-container > p")

			// Vérifier si une description est trouvée
			if description == "" {
				log.Printf("Aucune description trouvée pour l'annonce %d", index+1)
			} else {
				*details = append(*details, description)
			}
		})
	})
}

/**
 * setupMainPageCAImmobilier configure le collecteur pour la page principale de CA Immobilier.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]string} detailPageURLs - La liste des URLs des pages de détail.
 * @return {void}
 */
func setupMainPageCAImmobilier(collector *colly.Collector, detailPageURLs *[]string) {
	collector.OnHTML("div.results-container.mosaic", func(e *colly.HTMLElement) {
		// Parcourir chaque annonce (chaque div enfant)
		e.ForEach("div.columns.large-3", func(i int, annonceDiv *colly.HTMLElement) {
			// Vérifier qu'il ne s'agit pas de la dernière div (liens supplémentaires)
			if !annonceDiv.DOM.HasClass("sub_card-entities--blocliens") {
				// Récupérer la balise <article>
				annonceDiv.ForEach("article.sub_card-entities", func(_ int, article *colly.HTMLElement) {
					// Récupérer le lien dans la balise <a> avec le texte "Découvrir"
					href := article.ChildAttr("div.bottom-container div.bottom-bar a", "href")
					if href != "" {
						*detailPageURLs = append(*detailPageURLs, "https://www.ca-immobilier.fr/"+href)
					}
				})
			} else {
				log.Println("Div supplémentaire ignorée.")
			}
		})
	})
}

/**
 * processDetailPagesCAImmobilier extrait les références des annonces de la page de détail de CA Immobilier.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func setupMainPagePigeaultImmobilier(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div contenant toutes les annonces
	collector.OnHTML("div#liste_annonces", func(e *colly.HTMLElement) {
		log.Println("Div principale 'liste_annonces' trouvée")

		// Parcourir chaque <article> dans la div "row"
		e.ForEach("div.row article", func(_ int, article *colly.HTMLElement) {
			// Accéder à la balise <a> avec le href pour l'URL de l'annonce
			href := article.ChildAttr("a[rel='bookmark']", "href")
			if href != "" {
				log.Printf("URL trouvée : %s", href)
				*detailPageURLs = append(*detailPageURLs, href)
			} else {
				log.Println("Aucun href trouvé pour cet article")
			}
		})
	})
}

/**
 * processDetailPagesPigeaultImmobilier extrait les références des annonces de la page de détail de Pigeault Immobilier.
 * @param {colly.Collector} collector - Le collecteur à configurer.
 * @param {[]Announcement} announcements - La liste des annonces à remplir.
 * @return {void}
 */
func processDetailPagesPigeaultImmobilier(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler la div contenant les informations principales
	collector.OnHTML("div#top_infos", func(detail *colly.HTMLElement) {
		// Récupérer la référence dans la balise <p> avec la classe "ref"
		fullValue := detail.ChildText("p.ref")
		fullValue = strings.TrimSpace(fullValue)

		// Vérification et extraction de la référence
		if strings.HasPrefix(fullValue, "Réf :") {
			// Extraire uniquement la partie après "Réf :"
			reference := strings.TrimSpace(strings.TrimPrefix(fullValue, "Réf :"))
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
			log.Printf("Impossible de trouver la référence dans : %s", fullValue)
		}
	})
}

func setupMainPageLaForetImmobilier(collector *colly.Collector, detailPageURLs *[]string) {
	// Cibler la div principale contenant les annonces
	collector.OnHTML("div.properties__list", func(e *colly.HTMLElement) {
		// Parcourir chaque div avec la classe "row"
		e.ForEach("div.row", func(_ int, row *colly.HTMLElement) {
			// Parcourir chaque div contenant les annonces
			row.ForEach("div.col-md-6.col-lg-6.col-xl-4", func(_ int, annonce *colly.HTMLElement) {
				// Récupérer la valeur du href dans la balise <a>
				href := annonce.ChildAttr("a.apartment-card__link", "href")
				if href != "" {
					// Ajouter l'URL complète à la liste
					fullURL := e.Request.AbsoluteURL(href)
					*detailPageURLs = append(*detailPageURLs, fullURL)
				}
			})
		})
	})
}

func processDetailPagesLaForetImmobilier(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler la section contenant les informations de l'annonce
	collector.OnHTML("section.property__block.property-content", func(detail *colly.HTMLElement) {
		log.Println("Section 'property__block' trouvée")

		// Récupérer la référence web
		webRef := detail.ChildText("h5.text-base.text-ref:contains('Référence web')")
		webRef = strings.TrimSpace(strings.TrimPrefix(webRef, "Référence web :"))

		// Récupérer la référence agence
		agencyRef := detail.ChildText("h5.text-base.text-ref:contains('Référence Agence')")
		agencyRef = strings.TrimSpace(strings.TrimPrefix(agencyRef, "Référence Agence :"))

		// Vérifier si des références valides sont trouvées
		if webRef != "" || agencyRef != "" {
			// URL actuelle de la page
			url := detail.Request.URL.String()

			// Ajouter l'annonce avec les références à la liste
			*announcements = append(*announcements, Announcement{
				propertyReference: fmt.Sprintf("Web: %s, Agence: %s", webRef, agencyRef),
				url:               url,
			})
			log.Printf("Références trouvées - Web: %s, Agence: %s, URL: %s", webRef, agencyRef, url)
		} else {
			log.Println("Aucune référence trouvée dans cette annonce")
		}
	})
}

func setupMainPageCogir(collector *colly.Collector, announcementURLs *[]string) {
	// Cibler la div principale contenant les articles
	collector.OnHTML("div.listing_article.clearfix", func(element *colly.HTMLElement) {
		// Parcourir chaque balise <article> dans la div
		element.ForEach("article", func(_ int, article *colly.HTMLElement) {
			// Extraire la valeur de l'attribut href du lien dans la balise <article>
			href := article.ChildAttr("a.item-link", "href")
			if href != "" {
				// Ajouter l'URL complète à la liste
				*announcementURLs = append(*announcementURLs, href)
			}
		})
	})
}

func processDetailPagesCogir(collector *colly.Collector, announcements *[]Announcement) {
	// Cibler la section contenant les informations de l'annonce
	collector.OnHTML("div.detail_header", func(detail *colly.HTMLElement) {
		// Récupérer la référence de l'annonce
		ref := detail.ChildText("div.crit span:contains('Réf.')")
		ref = strings.TrimSpace(strings.TrimPrefix(ref, "Réf."))

		// Vérifier si une référence valide est trouvée
		if ref != "" {
			// URL actuelle de la page
			url := detail.Request.URL.String()

			// Ajouter l'annonce avec la référence à la liste
			*announcements = append(*announcements, Announcement{
				propertyReference: ref,
				url:               url,
			})
		}
	})
}
