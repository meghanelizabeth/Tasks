package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/meghanelizabeth/Tasks/db"
	"github.com/meghanelizabeth/Tasks/sessions"

	"context"
	"github.com/DataDog/dd-trace-go/tracer/contrib/gorilla/muxtrace"
)

var homeTemplate *template.Template
var deletedTemplate *template.Template
var completedTemplate *template.Template
var editTemplate *template.Template
var searchTemplate *template.Template
var templates *template.Template
var loginTemplate *template.Template

var message string //message will store the message to be shown as notification
var err error

//ShowAllTasksFunc is used to handle the "/" URL which is the default ons
//TODO add http404 error
func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	span, _ := muxtrace.GetRequestSpan(r)
	ctx := span.Context(context.Background())
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		context2, err := db.GetTasks(ctx, username, "pending", "")
		log.Println(context2)
		categories := db.GetCategories(ctx, username)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		} else {
			if message != "" {
				context2.Message = message
			}
			context2.CSRFToken = "abcd"
			context2.Categories = categories
			message = ""
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
			http.SetCookie(w, &cookie)
			homeTemplate.Execute(w, context2)
		}
	}
}

//ShowTrashTaskFunc is used to handle the "/trash" URL which is used to show the deleted tasks
func ShowTrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	span, _ := muxtrace.GetRequestSpan(r)
	ctx := span.Context(context.Background())
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		categories := db.GetCategories(ctx, username)
		context, err := db.GetTasks(ctx, username, "deleted", "")
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/trash", http.StatusInternalServerError)
		}
		if message != "" {
			context.Message = message
			message = ""
		}
		err = deletedTemplate.Execute(w, context)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//ShowCompleteTasksFunc is used to populate the "/completed/" URL
func ShowCompleteTasksFunc(w http.ResponseWriter, r *http.Request) {
	span, _ := muxtrace.GetRequestSpan(r)
	ctx := span.Context(context.Background())
	if r.Method == "GET" {
		username := sessions.GetCurrentUserName(r)
		categories := db.GetCategories(ctx, username)
		context, err := db.GetTasks(ctx, username, "completed", "")
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/completed", http.StatusInternalServerError)
		}
		completedTemplate.Execute(w, context)
	}
}

//ShowCategoryFunc will populate the /category/<id> URL which shows all the tasks related
// to that particular category
func ShowCategoryFunc(w http.ResponseWriter, r *http.Request) {
	span, _ := muxtrace.GetRequestSpan(r)
	ctx := span.Context(context.Background())
	if r.Method == "GET" && sessions.IsLoggedIn(r) {
		category := r.URL.Path[len("/category/"):]
		username := sessions.GetCurrentUserName(r)
		context, err := db.GetTasks(ctx, username, "", category)
		categories := db.GetCategories(ctx, username)

		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		if message != "" {
			context.Message = message
		}
		context.CSRFToken = "abcd"
		context.Categories = categories
		message = ""
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
		http.SetCookie(w, &cookie)
		homeTemplate.Execute(w, context)
	}
}
