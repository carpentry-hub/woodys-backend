package main

import (
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/middlewares"
	"github.com/carpentry-hub/woodys-backend/routes"
	"github.com/gorilla/mux"
)

func main() {

	db.DBConnection()

	
	r := mux.NewRouter()
	r.HandleFunc("/", routes.HomeHandler)


	// user routes handlers
	r.HandleFunc("/users/{id}", routes.GetUser).Methods("GET") 
	r.HandleFunc("/users/{id}/projects", routes.GetUserProjects).Methods("GET")
	r.HandleFunc("/users", routes.PostUser).Methods("POST")
	r.HandleFunc("/users/{id}", routes.PutUser).Methods("PUT")
	r.HandleFunc("/users/{id}", routes.DeleteUser).Methods("DELETE")
	

	// project routes handlers
	r.HandleFunc("/projects/search", routes.SearchProjects).Methods("GET")
	r.HandleFunc("/projects/{id:[0-9]+}", routes.GetProject).Methods("GET")
	r.HandleFunc("/projects", routes.PostProject).Methods("POST")
	r.HandleFunc("/projects/{id}", routes.PutProject).Methods("PUT")
	r.HandleFunc("/projects/{id}", routes.DeleteProject).Methods("DELETE")


	// comment routes handlers
	r.HandleFunc("/projects/{id}/comments", routes.GetProjectComments).Methods("GET")
	r.HandleFunc("/projects/{id}/comments", routes.PostProjectComment).Methods("POST") 
	r.HandleFunc("/comments/{id}", routes.DeleteComment).Methods("DELETE")
	r.HandleFunc("/comments/{id}/reply", routes.PostCommentReply).Methods("POST")
	r.HandleFunc("/comments/{id}/replies", routes.GetCommentReplies).Methods("GET")


	// rating routes handlers
	r.HandleFunc("/projects/{id}/ratings", routes.PostRating).Methods("POST")
	r.HandleFunc("/projects/{id}/ratings", routes.PutRating).Methods("PUT")
	r.HandleFunc("/projects/{id}/ratings", routes.GetRating).Methods("GET") 


	// project list routes handlers
	r.HandleFunc("/users/{id}/project-lists", routes.GetUsersProjectLists).Methods("GET")
	r.HandleFunc("/project-lists/{id}", routes.GetProjectLists).Methods("GET")
	r.HandleFunc("/project-lists", routes.PostProjectLists).Methods("POST")
	r.HandleFunc("/project-lists/{id}/projects", routes.AddProjectToList).Methods("POST") //This is a post method to project_list_item table
	r.HandleFunc("/project-lists/{id}", routes.PutProjectLists).Methods("PUT")
	r.HandleFunc("/project-list/{id}", routes.DeleteProjectList).Methods("DELETE")
	r.HandleFunc("/project-list/{list_id}/projects/{project_id}", routes.DeleteProjectFromList).Methods("DELETE") //Deletes the chosen project from the list
	

	http.ListenAndServe(":8080", middlewares.EnableCors(r))
}