package main

import (
	//"fmt"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// DATA STRUCTURES
type Page struct {
  Title string
  Body []byte // byte slice. what is expected by the io lib
}

// GLOBAL VARIABLES
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


/*
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:]) //slicing drops the leading /
}
*/

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return "", errors.New("invalid Page Title")
  }
  return m[2], nil // the title is the second subexpression.
}


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl+".html",p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// FUNCTION LITERALS and CLOSURES
/*
The closure returned by makeHandler is a function that takes an http.ResponseWriter and http.Request (in other words, an http.HandlerFunc). The closure extracts the title from the request path, and validates it with the validPath regexp. If the title is invalid, an error will be written to the ResponseWriter using the http.NotFound function. If the title is valid, the enclosed handler function fn will be called with the ResponseWriter, Request, and title as arguments.
*/
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  // called a closure. fn is one of the xxxxHandlers
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}


func viewHandler(w http.ResponseWriter,r *http.Request, title string) {
  /* transcended with the closures
  title, err := getTitle(w, r)
  if err != nil {
    return
  }*/
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  /* new and improved version above of the below: error handling!
  title := r.URL.Path[len("/view/"):]
  p, _ := loadPage(title)
  */
  renderTemplate(w, "view", p)
  //fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  /* same regex error handling as in other handlers.
  title := r.URL.Path[len("/edit/"):]
  p, err := loadPage(title)
  */
  /* closures. preventing code repetition.
  title, err := getTitle(w, r)
  if err != nil {
    return
  }*/
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
  /* Hardcoded html:
  fmt.Fprintf(w, "<h1>Editing %s</h1>"+
    "<form action=\"/save/%s\" method=\"POST\">"+
    "<textarea name=\"body\">%s</textarea><br>"+
    "<input type=\"submit\" value=\"Save\">"+
    "</form>",
    p.Title, p.Title, p.Body)
  */
  /* Code repetition:
  t, _ := template.ParseFiles("edit.html")
  t.Execute(w,p)
  */
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  /* closured;
  title, err := getTitle(w, r)
  if err != nil {
    return
  }*/
  //title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return os.WriteFile(filename, p.Body, 0600) // returns nil if allwell
  // 0600 is octal integer -> read-write
}

func loadPage (title string) (*Page, error) {
  filename := title + ".txt"
  body, err := os.ReadFile(filename) // dunder throws away error
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func main() {
  //p1 := &Page{Title: "TestPage", Body:[]byte("This is a sample Page.")}
  //p1.save()
  //p2, _ :=loadPage("TestPage")
  //fmt.Println(string(p2.Body))
  http.HandleFunc("/view/", makeHandler(viewHandler)) // http package handles requests to web root
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  log.Fatal(http.ListenAndServe(":8080",nil)) // wrapped with log.Fatal. blocking
}
