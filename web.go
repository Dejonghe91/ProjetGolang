package main

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"fmt"
	"strings"
)


//structure des pages {wiki, tablature, ...}
//structure chargement / sauvegarde
type PageTablature struct{
	Titre string
	Tab []byte
}

//structure d'affichage
type TablaturePage struct {
	Titre string
	Tab template.JS
}

//structure liste tablature
type PageTab struct {
	Tabs []string
}

//charge les pages en cache
var index = template.Must(template.ParseFiles("views/index.html"))
var list = template.Must(template.ParseFiles("views/list.html"))
var view = template.Must(template.ParseFiles("views/view.html"))
var create = template.Must(template.ParseFiles("views/create.html"))
var wiki = template.Must(template.ParseFiles("views/wiki.html"))


func main() {
	http.HandleFunc("/", homeHandler)		// Point d'entree de l'application
	http.HandleFunc("/wiki", wikiHandler)		// Page wiki de l'application how to build a tab
	http.HandleFunc("/tabs", tabsHandler)		// Liste des tableatures deja enregistrer sur le server
	http.HandleFunc("/create", createHandler)	// Page de creation / viso d'une tableatures
	http.HandleFunc("/view", viewHandler)		// Vue d'une tablature
	http.HandleFunc("/save", saveHandler)		// Sauvegarde d'une tablature
	
	//ressources externe de l'application, css / js / images / fonts
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	// tablatures de l'application
	http.Handle("/tablatures/", http.StripPrefix("/tablatures/", http.FileServer(http.Dir("tablatures"))))	
	
	fmt.Printf("Server listen on the port 9999 ...")
	http.ListenAndServe(":9999", nil)
}

//afficheras la page présentation
func homeHandler (w http.ResponseWriter, r *http.Request) {
	err := index.ExecuteTemplate(w, "index.html", nil)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//afficheras la page du wiki
func wikiHandler (w http.ResponseWriter, r *http.Request) { 

	err := wiki.ExecuteTemplate(w, "wiki.html", nil)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//afficheras la liste des tablatures/view"
func tabsHandler (w http.ResponseWriter, r *http.Request) {
	tabs, _ := ioutil.ReadDir("tablatures")
	
	//gestion erreur retour loadPage (redirection sur http.NotFound)
	
	p := PageTab{Tabs: nil}
	
	for _, v := range tabs {
		p.Tabs = append(p.Tabs, strings.Replace(v.Name(), ".txt", "", -1))
    }

	err := list.ExecuteTemplate(w, "list.html", p)	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//afficheras la page de cration d'une tablature
func createHandler (w http.ResponseWriter, r *http.Request) {
	err := create.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler (w http.ResponseWriter, r *http.Request) {

	//on sauvegarde dans un fichier
	p := &PageTablature{Titre:r.FormValue("titre"), Tab:[]byte(r.FormValue("sandbox"))}
	
	err := p.save()
	
	if err != nil {
		panic(err)
		// redirection vers page creation car erreur
	}
	
	http.Redirect(w, r, "/view?tab="+p.Titre, 201)
}

//afficheras la page de vue std d'une tablature
func viewHandler (w http.ResponseWriter, r *http.Request) {
	
	//chargement du fichier
	p, _:= loadPage(r.FormValue("tab"))
	
	//execution du template
	err := view.ExecuteTemplate(w, "view.html", &TablaturePage{Titre:p.Titre, Tab:template.JS(p.Tab)})
	if err != nil {
	      http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


//Tools
func loadPage(title string) (*PageTablature, error) {
    filename := "tablatures/"+ title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &PageTablature{Titre: title, Tab: body}, nil
}

func (p *PageTablature) save() error {
    filename := "tablatures/" + p.Titre + ".txt"
    return ioutil.WriteFile(filename, p.Tab, 0600)
}

