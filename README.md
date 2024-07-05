# Coding Challenge

This project is a Vehicle Parts Inventory Management system built with Go and React. It allows users to manage an inventory of vehicle parts, including creating, updating, deleting, and retrieving parts, as well as handling versioning of parts data.

## Features

- Add, update, delete, and list vehicle parts
- Versioning of parts data
- Retrieve specific versions of parts data


## Technologies Used

- **Backend**: Go, Gorilla Mux
- **Frontend**: React, Axios
- **Database**: In-memory storage (can be extended to persistent storage)
  

## Prerequisites

- Go 1.16+
- Node.js 14+
- npm 6+

## Getting Started

### Backend Setup

1. Clone the repository
   ```sh
   git clone https://github.com/BrandonTrapp88/pdmCodingChallenge.git
   ```
2. Install Go dependencies

  ```sh
  go mod tidy
```
3. Run the backend server

```sh
cd (FileName)/api
go run main.go repository.go handlers.go routers.go

The backend server will start on http://localhost:1710.
```
### Frontend Setup
1. Open another Terminal
   Navigate to the frontend directory

```sh

cd ../frontend

```
2. Install npm dependencies

```sh
npm install
```
3. Run the frontend server

```sh

npm start
The frontend server will start on http://localhost:3000.
```
### Project Structure
## Backend
main.go: Entry point of the application
repository.go: In-memory data storage with versioning
handlers.go: HTTP handlers for CRUD operations 
routers.go: Router configuration
##Frontend
src/
AddPartForm.js: Form for adding and editing parts
PartsList.js: Component for listing parts and handling part operations
SearchPage.js: Component for searching parts
NavBar.js: Navigation bar
App.js: Main application component
##API Endpoints
Parts
POST /parts: Create a new part
GET /parts: List all parts
GET /parts/{id}: Get a part by ID
PUT /parts/{id}: Update a part by ID
DELETE /parts/{id}: Delete a part by ID
GET /parts/{id}/version/{version}: Get a specific version of a part by ID and version
