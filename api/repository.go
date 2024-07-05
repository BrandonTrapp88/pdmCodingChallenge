package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Unable to reach database: %v", err)
		return nil, fmt.Errorf("unable to reach database: %v", err)
	}
	repo := &Repository{db: db}
	if err := repo.createTables(); err != nil {
		log.Fatalf("Unable to create tables: %v", err)
		return nil, fmt.Errorf("unable to create tables: %v", err)
	}
	return repo, nil
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
	// Check if part with same details exists
	existingPart, err := r.findPartByDetails(part)
	if err == nil && existingPart.ID != "" {
		if err := r.DeletePart(existingPart.ID); err != nil {
			return "", err
		}
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

func (r *Repository) findPartByDetails(part Part) (Part, error) {
	query := `SELECT id, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM parts WHERE name = ? AND sku = ? AND price = ?`
	row := r.db.QueryRow(query, part.Name, part.SKU, part.Price)

	var existingPart Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&existingPart.ID, &existingPart.Name, &images, &existingPart.SKU, &existingPart.Description, &existingPart.Price, &attributes, &fitmentData, &existingPart.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, nil // Part not found
		}
		return Part{}, err
	}

	if err := json.Unmarshal(images, &existingPart.Images); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(attributes, &existingPart.Attributes); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(fitmentData, &existingPart.FitmentData); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(shipment, &existingPart.Shipment); err != nil {
		return Part{}, err
	}
	if err := json.Unmarshal(metadata, &existingPart.Metadata); err != nil {
		return Part{}, err
	}

	return existingPart, nil
}

// get part from db
func (r *Repository) GetPart(id string) (Part, error) {
	query := `SELECT id, name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM parts WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var part Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&part.ID, &part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, fmt.Errorf("part not found")
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

// update part in db
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

// Delete part function
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

// List Part Function
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

// Get Part Version
func (r *Repository) GetPartVersion(id string, version int) (Part, error) {
	query := `SELECT name, images, sku, description, price, attributes, fitment_data, location, shipment, metadata FROM part_versions WHERE part_id = ? AND version = ?`
	row := r.db.QueryRow(query, id, version)

	var part Part
	var images, attributes, fitmentData, shipment, metadata []byte
	if err := row.Scan(&part.Name, &images, &part.SKU, &part.Description, &part.Price, &attributes, &fitmentData, &part.Location, &shipment, &metadata); err != nil {
		if err == sql.ErrNoRows {
			return Part{}, fmt.Errorf("version not found")
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

// list Part Version
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
