package main

/**
 * This is the main file for the Task application
 * License: MIT
 **/
import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/contrib/gorilla/muxtrace"
	"github.com/gorilla/mux"
	"github.com/thewhitetulip/Tasks/config"
	"github.com/thewhitetulip/Tasks/views"
)

func main() {
	values, err := config.ReadConfig("config.json")
	var port *string

	if err != nil {
		port = flag.String("port", "", "IP address")
		flag.Parse()

		//User is expected to give :8080 like input, if they give 8080
		//we'll append the required ':'
		if !strings.HasPrefix(*port, ":") {
			*port = ":" + *port
			log.Println("port is " + *port)
		}

		values.ServerPort = *port
	}

	views.PopulateTemplates()

	router := mux.NewRouter()
	muxTracer := muxtrace.NewMuxTracer("my-tasks-app", tracer.DefaultTracer)

	//Login logout
	muxTracer.HandleFunc(router, "/login/", views.LoginFunc)
	muxTracer.HandleFunc(router, "/logout/", views.RequiresLogin(views.LogoutFunc))
	muxTracer.HandleFunc(router, "/signup/", views.SignUpFunc)

	muxTracer.HandleFunc(router, "/add-category/", views.RequiresLogin(views.AddCategoryFunc))
	muxTracer.HandleFunc(router, "/add-comment/", views.RequiresLogin(views.AddCommentFunc))
	muxTracer.HandleFunc(router, "/add/", views.RequiresLogin(views.AddTaskFunc))

	//these handlers are used to delete
	muxTracer.HandleFunc(router, "/del-comment/", views.RequiresLogin(views.DeleteCommentFunc))
	muxTracer.HandleFunc(router, "/del-category/", views.RequiresLogin(views.DeleteCategoryFunc))
	muxTracer.HandleFunc(router, "/delete/", views.RequiresLogin(views.DeleteTaskFunc))

	//these handlers update
	muxTracer.HandleFunc(router, "/upd-category/", views.RequiresLogin(views.UpdateCategoryFunc))
	muxTracer.HandleFunc(router, "/update/", views.RequiresLogin(views.UpdateTaskFunc))

	//these handlers are used for restoring tasks
	muxTracer.HandleFunc(router, "/incomplete/", views.RequiresLogin(views.RestoreFromCompleteFunc))
	muxTracer.HandleFunc(router, "/restore/", views.RequiresLogin(views.RestoreTaskFunc))

	//these handlers fetch set of tasks
	muxTracer.HandleFunc(router, "/", views.RequiresLogin(views.ShowAllTasksFunc))
	muxTracer.HandleFunc(router, "/category/", views.RequiresLogin(views.ShowCategoryFunc))
	muxTracer.HandleFunc(router, "/deleted/", views.RequiresLogin(views.ShowTrashTaskFunc))
	muxTracer.HandleFunc(router, "/completed/", views.RequiresLogin(views.ShowCompleteTasksFunc))

	//these handlers perform action like delete, mark as complete etc
	muxTracer.HandleFunc(router, "/complete/", views.RequiresLogin(views.CompleteTaskFunc))
	muxTracer.HandleFunc(router, "/files/", views.RequiresLogin(views.UploadedFileHandler))
	muxTracer.HandleFunc(router, "/trash/", views.RequiresLogin(views.TrashTaskFunc))
	muxTracer.HandleFunc(router, "/edit/", views.RequiresLogin(views.EditTaskFunc))
	muxTracer.HandleFunc(router, "/search/", views.RequiresLogin(views.SearchTaskFunc))

	router.Handle("/static/", http.FileServer(http.Dir("public")))

	muxTracer.HandleFunc(router, "/api/get-task/", views.GetTasksFuncAPI)
	muxTracer.HandleFunc(router, "/api/get-deleted-task/", views.GetDeletedTaskFuncAPI)
	muxTracer.HandleFunc(router, "/api/add-task/", views.AddTaskFuncAPI)
	muxTracer.HandleFunc(router, "/api/update-task/", views.UpdateTaskFuncAPI)
	muxTracer.HandleFunc(router, "/api/delete-task/", views.DeleteTaskFuncAPI)

	muxTracer.HandleFunc(router, "/api/get-token/", views.GetTokenHandler)
	muxTracer.HandleFunc(router, "/api/get-category/", views.GetCategoryFuncAPI)
	muxTracer.HandleFunc(router, "/api/add-category/", views.AddCategoryFuncAPI)
	muxTracer.HandleFunc(router, "/api/update-category/", views.UpdateCategoryFuncAPI)
	muxTracer.HandleFunc(router, "/api/delete-category/", views.DeleteCategoryFuncAPI)

	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, router))
}
