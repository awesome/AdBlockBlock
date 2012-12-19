package index

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))

func init() {  
	http.HandleFunc("/", index)
}

func index(response http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		indexPost(response, request)
	} else {
		renderTemplate(response, "index", nil)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, model interface{}) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexPost(response http.ResponseWriter, request *http.Request) {
	url64 := request.FormValue("Id")
	if decBytes, decErr := base64.StdEncoding.DecodeString(url64); decErr == nil {
		url := string(decBytes)
		if strings.HasPrefix(url, "//") {
			url = "http:" + url
		}

		ctx := appengine.NewContext(request)
		client := urlfetch.Client(ctx)
		if getRes, getErr := client.Get(url); getErr == nil {
			defer getRes.Body.Close()
			if data, bodyErr := ioutil.ReadAll(getRes.Body); bodyErr == nil {
				// If content type is set to text/javascript the browser will execute it immediately.
				response.Header().Set("content-type", "text/plain")
				response.Write(data)
			} else {
				//log
			}
		}
	} else {
		//log
	}
}