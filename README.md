# Authentication Project

This project represents an authentication system built using Golang, Gin, GORM, and Redis.

## Installation and Running the Project

### Requirements
- Go 1.17+
- Redis
- PostgreSQL

### Installation Steps
1. Clone the repository:``` git clone https://github.com/zhosyaaa/auth-app.git```
2. Install dependencies: ```go mod tidy```

### Setting up Database and Redis
- Set up and configure PostgreSQL and Redis according to the .env file.
- Create postgres container: ```make postgres```
- Create the database: ```make createdb```
- Create the redis: ```make redis```

### Migrations
- Apply migrations: ```make migrateup```
- Rollback migrations: ```make migratedown```

### Running the Server
- Start the server:``` go run main.go```

## Project Structure
- internal/rest: Code related to handling HTTP requests.
- internal/db: Configuration and interaction with the database.
- internal/repository: Implementation of methods for database interaction.
- pkg: Utility packages and helper functions.

## Main Functionalities
- User Registration

url: http://localhost:8080/auth/register
body:
```
{
"username":"test",
"email":"zhake361@gmail.com",
"password":"test123"
}
```

- Authentication and Login
url: http://localhost:8080/auth/register
body:
```
{
    "username":"test",
    "password":"test123"
}
```
- Password Reset
url: http://localhost:8080/auth/forgotPassword
body:
```
{
    "email":"zhake361@gmail.com"
}
```
url: http://localhost:8080/auth/checkVerificationCode
// the code from the mail
```
{
    "email":"zhake361@gmail.com",
    "code":"636738" 
}
```
url: http://localhost:8080/auth/resetPassword
```
{
    "email": "zhake361@gmail.com",
    "password": "testtest",
    "passwordConfirm": "testtest"
}
```
- Fetching User Profile
- Deleting an Account

## Key Features and Considerations
- Dependency injection for database and Redis interactions.
- Use of JWT for authentication.
- Request handling through Gin and routing via Routers.
- Handling hashed passwords.
