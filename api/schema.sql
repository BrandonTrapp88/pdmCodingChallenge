CREATE DATABASE vehicle_parts_db;
USE vehicle_parts_db;

CREATE TABLE parts (
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
);

CREATE TABLE part_versions (
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
);
