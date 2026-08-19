package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sucip70/SVTest-demo-be/internal/model"
)

type PostsHandler struct {
	DB *sql.DB
}

func (h *PostsHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	query := strings.Builder{}
	query.WriteString("SELECT id, title, content, category, created_date, updated_date, status FROM posts")
	conditions := make([]string, 0, 5)
	args := make([]any, 0, 5)

	if title := r.URL.Query().Get("title"); title != "" {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+title+"%")
	}
	if content := r.URL.Query().Get("content"); content != "" {
		conditions = append(conditions, "content LIKE ?")
		args = append(args, "%"+content+"%")
	}
	if status := r.URL.Query().Get("status"); status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if category := r.URL.Query().Get("category"); category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}
	if createdDate := r.URL.Query().Get("created_date"); createdDate != "" {
		conditions = append(conditions, "DATE(created_date) = ?")
		args = append(args, createdDate)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	sortColumns := map[string]string{
		"title":        "title",
		"content":      "content",
		"status":       "status",
		"category":     "category",
		"created_date": "created_date",
	}
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_date"
	}
	sortColumn, ok := sortColumns[sortBy]
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid sort field")
		return
	}

	sortOrder := strings.ToUpper(r.URL.Query().Get("sort_order"))
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	if sortOrder != "ASC" && sortOrder != "DESC" {
		writeError(w, http.StatusBadRequest, "Invalid sort order")
		return
	}

	query.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder))
	rows, err := h.DB.Query(query.String(), args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}
	defer rows.Close()

	posts := []model.Posts{}
	for rows.Next() {
		var p model.Posts
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.UpdatedAt, &p.Status); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to retrieve posts")
			return
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve posts")
		return
	}

	writeJSON(w, http.StatusOK, posts)
}

func (h *PostsHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	row := h.DB.QueryRow("SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?", id)
	var p model.Posts
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.UpdatedAt, &p.Status); err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "Post not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to retrieve post")
		}
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var p model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := p.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO posts (title, content, category, created_date, updated_date, status) 
		VALUES (?, ?, ?, NOW(), NOW(), ?)`,
		p.Title, p.Content, p.Category, p.Status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	writeJSON(w, http.StatusCreated, model.CreatePostResponse{Message: "Post created successfully"})
}

func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var p model.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = p.ValidateUpdate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.DB.Exec("UPDATE posts SET title = ?, content = ?, category = ? WHERE id = ?",
		p.Title, p.Content, p.Category, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	writeJSON(w, http.StatusOK, model.UpdatePostResponse{Message: "Post updated successfully"})
}

func (h *PostsHandler) UpdatePostStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var status model.UpdatePostStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	err = status.ValidateStatus()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.DB.Exec("UPDATE posts SET status = ? WHERE id = ?", status.Status, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update post status")
		return
	}

	writeJSON(w, http.StatusOK, model.UpdatePostStatusResponse{Message: "Post status updated successfully"})
}

func writeJSON(w http.ResponseWriter, status int, response interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, model.APIResponse{
        Message: message,
    })
}