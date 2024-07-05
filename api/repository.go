package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// create db tables here
func (r *Repository) createTables() error {
	partsTable := `
	CREATE TABLE IF NOT EXISTS parts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		images JSON,
		sku VARCHAR(255),
		description TEXT,
		price DECIMAL(10, 2),
		attributes JSON,
		fitment_data JSON,
		location VARCHAR(255),
		shipment JSON,
		metadata JSON
	);`
	_, err := r.db.Exec(partsTable)
	if err != nil {
		log.Printf("Error creating parts table: %v", err)
		return fmt.Errorf("error creating parts table: %v", err)
	}

	partVersionsTable := `
	CREATE TABLE IF NOT EXISTS part_versions (
		version_id INT AUTO_INCREMENT PRIMARY KEY,
		part_id INT,
		version INT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name VARCHAR(255),
		images JSON,
		sku VARCHAR(255),
		description TEXT,
		price DECIMAL(10, 2),
		attributes JSON,
		fitment_data JSON,
		location VARCHAR(255),
		shipment JSON,
		metadata JSON,
		FOREIGN KEY (part_id) REFERENCES parts(id)
	);`
	_, err = r.db.Exec(partVersionsTable)
	if err != nil {
		log.Printf("Error creating part_versions table: %v", err)
		return fmt.Errorf("error creating part_versions table: %v", err)
	}

	return nil
}

// function to create Part and insert information into db
func (r *Repository) CreatePart(part Part) (string, error) {
	// Check if part with same SKU exists
	existingPart, err := r.findPartBySKU(part.SKU)
	if err == nil && existingPart.ID != "" {
		// Part exists, update the existing part
		return r.updateExistingPart(existingPart.ID, part)
	}

	// Marshal JSON fields
	images, err := json.Marshal(part.Images)
	if err != nil {
		return "", err
	}
	attributes, err := json.Marshal(part.Attributes)
	if err != nil {
		return "", err
	}
	fitmentData, err := json.Marshal(part.FitmentData)
	if err != nil {
		return "", err
	}
	shipment, err := json.Marshal(part.Shipment)
	if err != nil {
		return "", err
	}
	metadata, err := json.Marshal(part.Metadata)
	if err != nil {
		return "", err
	}

	// Insert part into the parts table
	query := `
		INSERT INTO parts (name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata)
	if err != nil {
		return "", err
	}

	partID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	// Insert the initial version into the part_versions table
	versionQuery := `
		INSERT INTO part_versions (part_id, version, timestamp, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(versionQuery, partID, 1, time.Now(), part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", partID), nil
}

func (r *Repository) updateExistingPart(id string, part Part) (string, error) {
	// Marshal JSON fields
	images, err := json.Marshal(part.Images)
	if err != nil {
		return "", err
	}
	attributes, err := json.Marshal(part.Attributes)
	if err != nil {
		return "", err
	}
	fitmentData, err := json.Marshal(part.FitmentData)
	if err != nil {
		return "", err
	}
	shipment, err := json.Marshal(part.Shipment)
	if err != nil {
		return "", err
	}
	metadata, err := json.Marshal(part.Metadata)
	if err != nil {
		return "", err
	}

	// Get the current version number
	var currentVersion int
	query := `SELECT COUNT(*) FROM part_versions WHERE part_id = ?`
	err = r.db.QueryRow(query, id).Scan(&currentVersion)
	if err != nil {
		return "", err
	}
	currentVersion++

	// Insert a new version in the part_versions table
	versionQuery := `
		INSERT INTO part_versions (part_id, version, timestamp, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(versionQuery, id, currentVersion, time.Now(), part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata)
	if err != nil {
		return "", err
	}

	// Update the existing part in the parts table
	updateQuery := `
		UPDATE parts SET name = ?, images = ?, sku = ?, description = ?, price = ?, attributes = ?, fitment_data = ?, location = ?, shipment = ?, metadata = ?
		WHERE id = ?
	`
	_, err = r.db.Exec(updateQuery, part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata, id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) findPartBySKU(sku string) (Part, error) {
	query := `SELECT id, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM parts WHERE sku = ?`
	row := r.db.QueryRow(query, sku)

	var part Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&part.ID, &part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, nil // Part not found
		}
		return Part{}, err
	}

	if err := json.Unmarshal(images, &part.Images); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(attributes, &part.Attributes); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(fitmentData, &part.FitmentData); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(shipment, &part.Shipment); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(metadata, &part.Metadata); err != nil {
		return Part{}, err
	}

	return part, nil
}

func (r *Repository) GetPart(id string) (Part, error) {
	query := `SELECT id, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM parts WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var part Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&part.ID, &part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, errors.New("part not found")
		}
		return Part{}, err
	}

	if err := json.Unmarshal(images, &part.Images); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(attributes, &part.Attributes); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(fitmentData, &part.FitmentData); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(shipment, &part.Shipment); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(metadata, &part.Metadata); err != nil {
		return Part{}, err
	}

	return part, nil
}

func (r *Repository) UpdatePart(id string, part Part) error {
	// Marshal JSON fields
	images, err := json.Marshal(part.Images)
	if err != nil {
		return err
	}
	attributes, err := json.Marshal(part.Attributes)
	if err != nil {
		return err
	}
	fitmentData, err := json.Marshal(part.FitmentData)
	if err != nil {
		return err
	}
	shipment, err := json.Marshal(part.Shipment)
	if err != nil {
		return err
	}
	metadata, err := json.Marshal(part.Metadata)
	if err != nil {
		return err
	}

	// Get the current version number
	var currentVersion int
	query := `SELECT COUNT(*) FROM part_versions WHERE part_id = ?`
	err = r.db.QueryRow(query, id).Scan(&currentVersion)
	if err != nil {
		return err
	}
	currentVersion++

	// Insert a new version in the part_versions table
	versionQuery := `
		INSERT INTO part_versions (part_id, version, timestamp, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.Exec(versionQuery, id, currentVersion, time.Now(), part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata)
	if err != nil {
		return err
	}

	// Update the existing part in the parts table
	updateQuery := `
		UPDATE parts SET name = ?, images = ?, sku = ?, description = ?, price = ?, attributes = ?, fitment_data = ?, location = ?, shipment = ?, metadata = ?
		WHERE id = ?
	`
	_, err = r.db.Exec(updateQuery, part.Name, images, part.SKU, part.Description, part.Price, attributes, fitmentData, part.Location, shipment, metadata, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeletePart(id string) error {
	query := `DELETE FROM parts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM part_versions WHERE part_id = ?`
	_, err = r.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListParts() ([]Part, error) {
	query := `SELECT id, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM parts`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []Part
	for rows.Next() {
		var part Part
		var images, attributes, fitmentData, shipment, metadata []byte
		if err := rows.Scan(&part.ID, &part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(images, &part.Images); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(attributes, &part.Attributes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(fitmentData, &part.FitmentData); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(shipment, &part.Shipment); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadata, &part.Metadata); err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func (r *Repository) GetPartVersion(id string, version int) (Part, error) {
	query := `SELECT name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM part_versions WHERE part_id = ? AND version = ?`
	row := r.db.QueryRow(query, id, version)

	var part Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, errors.New("version not found")
		}
		return Part{}, err
	}

	if err := json.Unmarshal(images, &part.Images); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(attributes, &part.Attributes); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(fitmentData, &part.FitmentData); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(shipment, &part.Shipment); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(metadata, &part.Metadata); err != nil {
		return Part{}, err
	}

	return part, nil
}

func (r *Repository) ListPartVersions(id string) ([]PartVersion, error) {
	query := `SELECT version, timestamp FROM part_versions WHERE part_id = ? ORDER BY version`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []PartVersion
	for rows.Next() {
		var version PartVersion
		if err := rows.Scan(&version.Version, &version.Timestamp); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}
