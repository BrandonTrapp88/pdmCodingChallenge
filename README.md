# Coding Challenge

This project is a Vehicle Parts Inventory Management system built with Go and React. It allows users to manage an inventory of vehicle parts, including creating, updating, deleting, and retrieving parts, as well as handling versioning of parts data.

## Features

- Add, update, delete, and list vehicle parts
- Versioning of parts data
- Retrieve specific versions of parts data

## Technologies Used

- **Backend**: Go, Gorilla Mux
- **Frontend**: React, Axios
- **Database**: MySQL (vehicle_parts_db)

## Prerequisites

- Go 1.16+
- Node.js 14+
- npm 6+
- MySQL Server

## Getting Started

### Backend Setup

1. Clone the repository:
   ``` sh
   git clone https://github.com/BrandonTrapp88/pdmCodingChallenge.git
   ```

Install Go dependencies:

``` sh

go mod tidy
```
Set up your MySQL database and create a database named vehicle_parts_db.

Update the main.go file with your MySQL database credentials.

Run the backend server:

``` sh

cd pdmCodingChallenge/api
go run main.go repository.go handlers.go routers.go
```


Frontend Setup
Open another terminal and navigate to the frontend directory:

``` sh
cd ../frontend
``` 
Install npm dependencies:

``` sh

npm install
``` 
Run the frontend server:

``` sh

npm start

```


### Project Structure
# Backend
- main.go: Entry point of the application
- repository.go: Data storage with versioning
- handlers.go: HTTP handlers for CRUD operations
- routers.go: Router configuration
# Frontend
- src/
- AddPartForm.js: Form for adding and editing parts
- PartsList.js: Component for listing parts and handling part operations
- SearchPage.js: Component for searching parts
- NavBar.js: Navigation bar
- App.js: Main application component
# API Endpoints
- Parts
- POST /parts: Create a new part
- GET /parts: List all parts
- GET /parts/{id}: Get a part by ID
- PUT /parts/{id}: Update a part by ID
- DELETE /parts/{id}: Delete a part by ID
- GET /parts/{id}/version/{version}: Get a specific version of a part by ID and version
### Makefile Commands
To simplify the process of running the API and frontend servers, use the provided Makefile.

Commands
Run the API server

``` sh

make api DB_USER=USERNAME DB_PASSWORD=PASSWORD
``` 
Build the API server

``` sh

make api-build
``` 
Clean API build artifacts

``` sh

make api-clean
```
Start the frontend development server

``` sh

make frontend
``` 
Build the frontend for production

``` sh

make frontend-build
``` 
Clean frontend build artifacts

``` sh

make frontend-clean
```
Run both API and frontend

``` sh

make all DB_USER=USERNAME DB_PASSWORD=PASSWORD
``` 
Build both API and frontend

``` sh

make build
```
Clean both API and frontend

``` sh
make clean
```
Display help message

``` sh
make help
```
