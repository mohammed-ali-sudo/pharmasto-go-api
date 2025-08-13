// internal/services/teacher.go
package services

import (
	"database/sql"
	"errors"
	"fmt"
	"goapi/models"
	"goapi/shared"
	"strings"

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
			return models.Teacher{}, shared.ErrorHandling(err, "teacher not found")
		}
		return models.Teacher{}, shared.ErrorHandling(err, "error during data process")
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
		return nil, shared.ErrorHandling(err, "database error")
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

// put/teachers/id
func UpdateTeacherService(db *sql.DB, t models.Teacher, id int) (models.Teacher, error) {
	// will hold the updated row
	var teacher models.Teacher

	// single query to update + return the new values
	const qry = `
        UPDATE teachers
           SET first_name = $1,
               last_name  = $2,
               email      = $3,
               class      = $4,
               subject    = $5
         WHERE id = $6
     RETURNING id, first_name, last_name, email, class, subject
    `

	// run it and scan into teacher
	err := db.QueryRow(
		qry,
		t.FirstName,
		t.LastName,
		t.Email,
		t.Class,
		t.Subject,
		id,
	).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Class,
		&teacher.Subject,
	)
	if err != nil {
		return models.Teacher{}, fmt.Errorf("update failed: %w", err)
	}

	// success! return the updated struct
	return teacher, nil
}

func PatchTeacherService(db *sql.DB, t models.Teacher, id int) (models.Teacher, error) {
	const qry = `
        UPDATE teachers
           SET first_name = COALESCE(NULLIF($1, ''), first_name),
               last_name  = COALESCE(NULLIF($2, ''), last_name),
               email      = COALESCE(NULLIF($3, ''), email),
               class      = COALESCE(NULLIF($4, ''), class),
               subject    = COALESCE(NULLIF($5, ''), subject)
         WHERE id = $6
     RETURNING id, first_name, last_name, email, class, subject`

	var teacher models.Teacher
	err := db.QueryRow(
		qry,
		t.FirstName,
		t.LastName,
		t.Email,
		t.Class,
		t.Subject,
		id,
	).Scan(
		&teacher.ID,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Class,
		&teacher.Subject,
	)
	if err != nil {
		return models.Teacher{}, fmt.Errorf("patch failed: %w", err)
	}
	return teacher, nil
}

func FilterServices(db *sql.DB, filters map[string]string) ([]models.Teacher, error) {
	// 1. Base query
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers"

	// 2. Build WHERE clauses with $1, $2, …
	var where []string
	var args []interface{}
	idx := 1

	for _, key := range []string{"first_name", "last_name", "email", "class", "subject"} {
		if val, ok := filters[key]; ok && val != "" {
			where = append(where, fmt.Sprintf("%s = $%d", key, idx))
			args = append(args, val)
			idx++
		}
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// 3. Execute
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 4. Scan results
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

// SortServices returns all teachers ordered by the given sortParams.
// sortParams entries must be in the form "field" or "field:asc" or "field:desc".
// If any entry is malformed or uses an unsupported field/direction, an error is returned.
//
// Example valid sortParams:
//
//	[]string{"last_name:desc", "id:asc"}
//	[]string{"first_name"}            // defaults to ASC
func SortServices(db *sql.DB, sortParams []string) ([]models.Teacher, error) {
	// 1) Base SELECT
	query := `
        SELECT
            id,
            first_name,
            last_name,
            email,
            class,
            subject
        FROM teachers
    `

	// 2) Whitelist of allowed fields
	allowed := map[string]bool{
		"id":         true,
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	// 3) Build ORDER BY clauses
	var orders []string
	for _, p := range sortParams {
		// split into field and optional direction
		parts := strings.SplitN(p, ":", 2)
		field := strings.TrimSpace(parts[0])
		if field == "" {
			return nil, fmt.Errorf("invalid sort parameter: %q", p)
		}
		if !allowed[field] {
			return nil, fmt.Errorf("unsupported sort field: %q", field)
		}

		// determine direction
		dir := "ASC"
		if len(parts) == 2 {
			switch strings.ToUpper(strings.TrimSpace(parts[1])) {
			case "ASC", "":
				dir = "ASC"
			case "DESC":
				dir = "DESC"
			default:
				return nil, fmt.Errorf("invalid sort direction %q in parameter %q", parts[1], p)
			}
		}

		orders = append(orders, fmt.Sprintf("%s %s", field, dir))
	}

	// 4) Append ORDER BY if any
	if len(orders) > 0 {
		query += " ORDER BY " + strings.Join(orders, ", ")
	}

	// 5) Execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// 6) Scan results
	var teachers []models.Teacher
	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(
			&t.ID,
			&t.FirstName,
			&t.LastName,
			&t.Email,
			&t.Class,
			&t.Subject,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		teachers = append(teachers, t)
	}
	// check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return teachers, nil
}

func Delet(db *sql.DB, id int) error {
	// Execute the DELETE statement. For MySQL, use a `?` placeholder.
	result, err := db.Exec("DELETE FROM teachers WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Get the number of rows affected by the DELETE statement.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were deleted, return an error.
	if rowsAffected == 0 {
		return fmt.Errorf("teacher with id %d not found", id)
	}

	// If the deletion was successful, return nil.
	return nil
}

func PatchTeachersServices(db *sql.DB, teachers []map[string]any) error {
	transaction, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer transaction.Rollback()

	const qry = `
        UPDATE teachers
        SET first_name = COALESCE(NULLIF($1, ''), first_name),
            last_name = COALESCE(NULLIF($2, ''), last_name),
            email = COALESCE(NULLIF($3, ''), email),
            class = COALESCE(NULLIF($4, ''), class),
            subject = COALESCE(NULLIF($5, ''), subject)
        WHERE id = $6`

	for _, value := range teachers {
		id, ok := value["id"].(float64)
		if !ok {
			return errors.New("invalid id format or missing id")
		}

		// Use a helper function to safely get string values.
		getString := func(key string) string {
			if val, ok := value[key].(string); ok {
				return val
			}
			return ""
		}

		firstName := getString("first_name")
		lastName := getString("last_name")
		email := getString("email")
		class := getString("class")
		subject := getString("subject")

		result, err := transaction.Exec(
			qry,
			firstName,
			lastName,
			email,
			class,
			subject,
			int(id),
		)
		if err != nil {
			return fmt.Errorf("failed to update teacher with id %d: %w", int(id), err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for teacher with id %d: %w", int(id), err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("teacher with id %d not found", int(id))
		}
	}

	return transaction.Commit()
}
