package model

import (
	"time"
	"fmt"
)

type Posts struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_date"`
	UpdatedAt time.Time `json:"updated_date"`
	Status    string    `json:"status"`
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type UpdatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type UpdatePostStatusRequest struct {
	Status string `json:"status"`
}

type CreatePostResponse struct {
	Message string `json:"message"`
}

type UpdatePostResponse struct {
	Message string `json:"message"`
}

type UpdatePostStatusResponse struct {
	Message string `json:"message"`
}

func (p *CreatePostRequest) Validate() error {
	if p.Title == "" {
		return fmt.Errorf("title is required")
	}

	if p.Content == "" {
		return fmt.Errorf("content is required")
	}

	if p.Category == "" {
		return fmt.Errorf("category is required")
	}

	if p.Status != "Publish" && p.Status != "Draft" && p.Status != "Trash" {
		return fmt.Errorf("wrong status")
	}

	return nil
}

func (p *UpdatePostRequest) ValidateUpdate() error {
	if p.Title == "" {
		return fmt.Errorf("title is required")
	}

	if p.Content == "" {
		return fmt.Errorf("content is required")
	}

	if p.Category == "" {
		return fmt.Errorf("category is required")
	}

	return nil
}

func (p *UpdatePostStatusRequest) ValidateStatus() error {
	if p.Status != "Publish" && p.Status != "Draft" && p.Status != "Trash" {
		return fmt.Errorf("wrong status")
	}
	return nil
}