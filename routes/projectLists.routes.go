package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carpentry-hub/woodys-backend/db"
	"github.com/carpentry-hub/woodys-backend/models"
	"github.com/gorilla/mux"
)

// obtener todas las listas de un usuario - Requiere user_id
func GetUsersProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userIDString := params["id"]

	// chequeo existencia del usuario
	userID, err := strconv.Atoi(userIDString) // cambio de str a int para evitar errores
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	// realizacion de la query y manejo de errores
	var lists []models.ProjectList
	if err := db.DB.Where("user_id = ?", userID).Find(&lists).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching Project Lists"))
		return
	}

	json.NewEncoder(w).Encode(&lists)
}

// obtener una lista - Requier id
func GetProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	db.DB.First(&list, params["id"])

	if list.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Project List not found"))
	} else {
		json.NewEncoder(w).Encode(&list)
	}
}

// postear una lista
func PostProjectLists(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	json.NewDecoder(r.Body).Decode(&list)

	createdList := db.DB.Create(&list)
	err := createdList.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // satatus code 400
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&list)
	}
}

// postear un project list item (aniadir un proyecto a una lista)
func AddProjectToList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	json.NewDecoder(r.Body).Decode(&item)

	createdItem := db.DB.Create(&item)
	err := createdItem.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&item)
}

// actualizar una lista - Requiere id
func PutProjectLists(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// chequeo que el proyecto ya exista
	var existing models.ProjectList
	if err := db.DB.First(&existing, params["id"]).Error; err != nil {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project List Not Found"))
		return
	}

	// leo el updated
	var updated models.ProjectList
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest) // status code 400
		w.Write([]byte("Error on json file"))
		return
	}

	// actualizar campos
	existing.Name = updated.Name
	existing.IsPublic = updated.IsPublic

	// guardar en DB
	if err := db.DB.Save(&existing).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save the Project List"))
		return
	}

	json.NewEncoder(w).Encode(&existing)
}

// borrar una lista - Requiere id
func DeleteProjectList(w http.ResponseWriter, r *http.Request) {
	var list models.ProjectList
	params := mux.Vars(r)
	db.DB.First(&list, params["id"])

	if list.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project List not found"))
	} else {
		db.DB.Unscoped().Delete(&list)
	}
}

// borrar un proyecto de una lista - Requiere id
func DeleteProjectFromList(w http.ResponseWriter, r *http.Request) {
	var item models.ProjectListItem
	params := mux.Vars(r)
	db.DB.First(&item, params["id"])

	if item.ID == 0 {
		w.WriteHeader(http.StatusNotFound) // status code 404
		w.Write([]byte("Project List Item not found"))
	} else {
		db.DB.Unscoped().Delete(&item)
	}
}

