package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
  Title string 
  Body []byte

}

func (p *Page) save() error {  
  filename := p.Title + ".txt"
  return os.WriteFile(filename, p.Body, 0600)

}

func GetTitle(w http.ResponseWriter, r *http.Request) (string, error){
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("Invalid Page Title") 
  }
  return m[2], nil
}


func LoadPage(title string) (*Page, error){
  filename := title + ".txt"
  body, err := os.ReadFile(filename)
  if err != nil{
    return nil, err
  }

  return &Page{Title: title, Body: body}, nil

}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string){
  p, err := LoadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  RenderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string){
  p, err := LoadPage(title)
  if err != nil{
    p = &Page{Title: title}
  }  
  RenderTemplate(w, "edit", p)

}


func SaveHandler(w http.ResponseWriter, r *http.Request, title string){
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html",p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

}

func MakeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request){
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}


func main(){
  http.HandleFunc("/view/", MakeHandler(ViewHandler))
  http.HandleFunc("/edit/", MakeHandler(EditHandler))
  http.HandleFunc("/save/", MakeHandler(SaveHandler))
  log.Fatal(http.ListenAndServe(":8080", nil))

}
