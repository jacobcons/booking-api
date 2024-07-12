# booking-api 

REST API that handles bookings (note: the folder structure is a bit convoluted because I wanted to test how I'd structure a larger project)

### Running locally -
1. Install go modules `go mod download`
2. Install goose (for migrations) `go install github.com/pressly/goose/v3/cmd/goose@latest`
3. Set env variables
    ```
    APP_ENV development

    DBSTRING <your-db-string>?sslmode=disable
    
    JWT_SECRET <your-jwt-secret>
    
    GOOSE_DRIVER postgres
    GOOSE_DBSTRING <your-db-string>
    GOOSE_MIGRATION_DIR ./db/migrations
    ```
4. Run migrations `goose up`
5. Seed db with users (script will also output valid jwts for those users) `go run ./cmd/seed.go`
6. Run server `go run ./cmd/server.go`
7. (optional) Run server with live reloading and debugging (config only works on windows) 
   1. `go install github.com/air-verse/air@latest` 
   2. `air`

### Endpoints
- GET /bookings?startDatetime=2024-08-13T09:00:00Z&endDatetime=2024-08-13T17:00:00Z
  - Get bookings between startDatetime and endDatetime
- POST /bookings 
  - Create booking for logged-in user, datetimes must be on a weekday, between 9am-5pm, on the hour and in the future
  - Body `{ startDatetime: "2024-11-13T15:00:00Z", endDatetime: "2024-11-13T16:00:00Z" }`
  - Include jwt as bearer token
  
- DELETE /bookings/:id
  - Delete booking with given id belonging to logged-in user
  - Include jwt as bearer token