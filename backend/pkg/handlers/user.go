package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"social/models"
	"social/pkg/db/sqlite"
	"social/queries/middleware"
	"social/queries/utils"

	"github.com/gorilla/mux"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := models.GetUserByID(sqlite.DB, userID)
	if err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetUserID(r.Context())

	if user.IsPrivate && currentUserID != user.ID {
		// TODO: Check if current user is a follower
		// For now, return limited info
		limitedUser := map[string]interface{}{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"avatar":     user.Avatar,
			"is_private": user.IsPrivate,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(limitedUser)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateProfilePrivacy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileUserID := vars["id"]

	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID != profileUserID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	var req struct {
		IsPrivate bool `json:"is_private"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := models.SetProfilePrivacy(sqlite.DB, currentUserID, req.IsPrivate); err != nil {
		http.Error(w, `{"error":"Failed to update privacy"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"is_private": req.IsPrivate})
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileUserID := vars["id"]

	currentUserID := middleware.GetUserID(r.Context())
	if currentUserID != profileUserID {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	user, err := models.GetUserByID(sqlite.DB, currentUserID)
	if err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, `{"error":"Invalid multipart form"}`, http.StatusBadRequest)
			return
		}

		if firstName := strings.TrimSpace(r.FormValue("first_name")); firstName != "" {
			user.FirstName = firstName
		}
		if lastName := strings.TrimSpace(r.FormValue("last_name")); lastName != "" {
			user.LastName = lastName
		}
		if nickname := strings.TrimSpace(r.FormValue("nickname")); nickname != "" {
			user.Nickname = &nickname
		}
		if aboutMe := strings.TrimSpace(r.FormValue("about_me")); aboutMe != "" {
			user.AboutMe = &aboutMe
		}

		file, header, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()
			avatarPath, err := utils.SaveImage(file, header)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			user.Avatar = &avatarPath
		} else if err != http.ErrMissingFile {
			http.Error(w, `{"error":"Invalid avatar upload"}`, http.StatusBadRequest)
			return
		}
	} else {
		var req struct {
			FirstName string  `json:"first_name"`
			LastName  string  `json:"last_name"`
			Avatar    *string `json:"avatar"`
			Nickname  *string `json:"nickname"`
			AboutMe   *string `json:"about_me"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.FirstName) != "" {
			user.FirstName = req.FirstName
		}
		if strings.TrimSpace(req.LastName) != "" {
			user.LastName = req.LastName
		}
		if req.Avatar != nil {
			user.Avatar = req.Avatar
		}
		if req.Nickname != nil {
			user.Nickname = req.Nickname
		}
		if req.AboutMe != nil {
			user.AboutMe = req.AboutMe
		}
	}

	if err := models.UpdateUser(sqlite.DB, user); err != nil {
		http.Error(w, `{"error":"Failed to update profile"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
