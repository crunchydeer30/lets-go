package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
	Insert(title string, content string) (int, error)
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() AND id = $1`

	s := &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}

func (m *SnippetModel) Insert(title string, content string) (int, error) {
	var id int
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, NOW(), NOW() + INTERVAL '7 days') RETURNING id`

	err := m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
