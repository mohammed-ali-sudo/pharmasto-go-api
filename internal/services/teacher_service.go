// internal/services/teacher.go
package services

import (
	"database/sql"
	"errors"
	"fmt"
	"goapi/models"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
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

func GetTeacherByID(db *sql.DB, id int) (models.Teacher, error) {
	var teacher models.Teacher
	const query = `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Class,
		&teacher.Subject,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Teacher{}, errors.New("teacher not found")
		}
		return models.Teacher{}, err
	}

	return teacher, nil
}

func GetTeachers(db *sql.DB, ids []int) ([]models.Teacher, error) {

	if len(ids) == 0 {
		return nil, errors.New("teachers ids must be not nil")
	}

	query := `SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ANY($1)`
	rows, err := db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher

	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Email, &t.Class, &t.Subject); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, nil
}
