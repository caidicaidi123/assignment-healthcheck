# Assignment-Healthcheck

## Frontend (React)
1. Navigate to /frontend  
2. npm install  
3. npm start  
The app is running at http://localhost:8080  
Refresh the page if new website does not appear, since some website has relatively long response time . 

## Backend (Go)
1. Navigate to /backend
2. go get github.com/gorilla/mux github.com/gorilla/handlers  
3. go run main.go app.go model.go  
Server is running at http://localhost:8000  

## Run API Tests
1. Navigate to /backend  
2. go get github.com/gorilla/mux github.com/gorilla/handlers  
3. go test -v

## Endpoints
POST /api/healthcheck  
GET /api/healthcheck  
DELETE /api/healthcheck  
