package services

import (
	"database/sql"
	"fmt"
	"goapi/models"
)

// EnsureManfactoryTable ensures the table exists before inserting

// CreateManfactory inserts a new Manfactory into the DB
func CreateManfactory(db *sql.DB, t models.Manfactory) (models.Manfactory, error) {

	// Insert into DB only if manfactory_name does not exist
	sqlStmt := `
        INSERT INTO manfactory 
            (manfactory_name, manfactory_country, email, contact_number, license_number)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	err := db.QueryRow(
		sqlStmt,
		t.ManfactoryName,
		t.ManfactoryCountry,
		t.ContactEmail,
		t.ContactNumber,
		t.LicenseNumber,
	).Scan(&t.ID)

	if err == sql.ErrNoRows {
		// No new row inserted because manfactory_name exists
		return t, fmt.Errorf("manfactory with name '%s' already exists", t.ManfactoryName)
	}

	if err != nil {
		return models.Manfactory{}, err
	}

	return t, nil
}
