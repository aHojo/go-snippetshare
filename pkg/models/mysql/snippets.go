package mysql

import (
	"database/sql"

	"github.com/ahojo/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into our database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	// Create an insert statement

	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Execute the statement using the Exec method
	// The Exec method returns a Result object that can be used to get the last inserted ID
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use tthe LastInsertId() method to get the ID of the last inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}



// This function will return a snippet based on the ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Query the database with QueryRow(). It will only return one row
	row := m.DB.QueryRow(stmt, id)

	// Initialize a snippet struct pointer
	snippet := &models.Snippet{}

	// Scan the row into the snippet struct.
	// Arguements to the scan must be pointers
	// The number of arguements must be exactly the same as the number of columns in the result set
	// if now rows, sql.ErrNoRows, is returned
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecords
	} else if err != nil {
		return nil, err
	}

	return snippet, nil
}

// Return the 10 most recent snippets
func (m *SnippetModel) GetRecent() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Query() will return a slice of pointers to rows
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Make sure we clean up the connection before returning.
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		// Create a single snippet struct pointer
		snippet := &models.Snippet{}

		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)

	}

	return snippets, nil
}
