package main

import (
  "os"
  "log"
  "net/http"
  "html/template"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
  Title string 
  Body []byte

}

func (p *Page) save() error {  
  filename := p.Title + ".txt"
  return os.WriteFile(filename, p.Body, 0600)

}

func LoadPage(title string) (*Page, error){
  filename := title + ".txt"
  body, err := os.ReadFile(filename)
  if err != nil{
    return nil, err
  }

  return &Page{Title: title, Body: body}, nil

}

func ViewHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/view/"):]
  p, _ := LoadPage(title)
  RenderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/edit/"):]
  p, err := LoadPage(title)
  if err != nil{
    p = &Page{Title: title}
  }  
  RenderTemplate(w, "edit", p)


}


func SaveHandler(w http.ResponseWriter, r *http.Request){
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err :=  p.save()
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



func main(){
  http.HandleFunc("/view/", ViewHandler)
  http.HandleFunc("/edit/", EditHandler)
  http.HandleFunc("/save/", SaveHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))


}
