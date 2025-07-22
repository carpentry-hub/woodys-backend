package routes

import (
	"encoding/json"
	"net/http"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// obtener todos los comentarios de un proyecto - Requiere project_id
func GetProjectComments(w http.ResponseWriter, r *http.Request) {

}

// postear un comentario a un proyecto - Requiere project_id y parent_comment_id = 0
func PostProjectComment(w http.ResponseWriter, r *http.Request){
	var comment models.Comment
	json.NewDecoder(r.Body).Decode(&comment)

	createdComment := db.DB.Create(&comment)
	err := createdComment.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&comment)
	}
	
}

// borrar un comentario de un proyecto - Requiere id
func DeleteComment(w http.ResponseWriter, r *http.Request){
	var comment models.Comment
	params := mux.Vars(r)
	db.DB.First(&comment, params["id"])

	if comment.ID == 0 {
		w.WriteHeader(http.StatusNotFound) //status code 404
		w.Write([]byte("Comment not found"))
	} else {
		db.DB.Unscoped().Delete(&comment)
	}
}

// postear una respuesta a un comentario - Requiere project_id y parent_comment_id
func PostCommentReply(w http.ResponseWriter, r *http.Request){
	var commentReply models.Comment
	json.NewDecoder(r.Body).Decode(&commentReply)

	createdComment := db.DB.Create(&commentReply)
	err := createdComment.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&commentReply)
	}
}

// obtener las respuestas a un comentario - Requiere id
func GetCommentReplies(w http.ResponseWriter, r *http.Request){
	
}