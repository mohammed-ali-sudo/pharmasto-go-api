// internal/services/teacher.go
package services

import (
	"database/sql"
	"fmt"
	"goapi/models"

	"github.com/go-playground/validator/v10"
)

// AddTeacher builds the email, runs validation, inserts into DB, and returns the new record.
func AddTeacher(db *sql.DB, t models.Teacher) (models.Teacher, error) {
    // 1. Build derived field
    t.Email = fmt.Sprintf("%s.%s@example.com", t.FirstName, t.LastName)

    // 2. Validate struct tags
    if err := validator.New().Struct(t); err != nil {
        // return the raw validator error (you can wrap or map to APIError here)
        return models.Teacher{}, err
    }

    // 3. Insert & RETURNING id
    const sqlStmt = `
        INSERT INTO teachers (first_name, last_name, email, class, subject)
        VALUES ($1,$2,$3,$4,$5)
        RETURNING id`
    if err := db.QueryRow(sqlStmt,
        t.FirstName, t.LastName, t.Email, t.Class, t.Subject,
    ).Scan(&t.ID); err != nil {
        return models.Teacher{}, err
    }

    return t, nil
}
