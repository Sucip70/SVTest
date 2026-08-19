package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/Sucip70/SVTest-demo-be/internal/model"
	"fmt"
)

type PostsHandler struct {
	DB *sql.DB
}

func (h *PostsHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, content, category, created_date, updated_date, status FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []model.Posts{}
	for rows.Next() {
		var p model.Posts
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.UpdatedAt, &p.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostsHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	row := h.DB.QueryRow("SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?", id)
	var p model.Posts
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.UpdatedAt, &p.Status); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var p model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := p.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO posts (title, content, category, created_date, updated_date, status) 
		VALUES (?, ?, ?, NOW(), NOW(), ?)`,
		p.Title, p.Content, p.Category, p.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.CreatePostResponse{Message: "Post created successfully"})
}

func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var p model.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = p.ValidateUpdate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("UPDATE posts SET title = ?, content = ?, category = ? WHERE id = ?",
		p.Title, p.Content, p.Category, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UpdatePostResponse{Message: "Post updated successfully"})
}

func (h *PostsHandler) UpdatePostStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var status model.UpdatePostStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = status.ValidateStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("UPDATE posts SET status = ? WHERE id = ?", status.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UpdatePostStatusResponse{Message: "Post status updated successfully"})
}