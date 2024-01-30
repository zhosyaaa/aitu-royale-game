# Game Project

This project involves building a gaming application using Golang, Gin, and a PostgreSQL database.
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

POST: http://localhost:8080/app/auth/register
```
{
"username":"test",
"email":"zhake361@gmail.com",
"password":"test123"
}
```

- Authentication and Login
POST: http://localhost:8080/app/auth/register
```
{
    "username":"test",
    "password":"test123"
}
```
- Password Reset
POST: http://localhost:8080/app/auth/forgotPassword
```
{
    "email":"zhake361@gmail.com"
}
```
POST: http://localhost:8080/app/auth/checkVerificationCode
// the code from the mail
```
{
    "email":"zhake361@gmail.com",
    "code":"636738" 
}
```
POST: http://localhost:8080/app/auth/resetPassword
```
{
    "email": "zhake361@gmail.com",
    "password": "testtest",
    "passwordConfirm": "testtest"
}
```
- Fetching User Profile
- Deleting an Account

- Create Hero
POST: http://localhost:8080/app/game/create-hero
```{
    "Name": "Iron Golem",
    "Description": "A colossal golem crafted from sturdy iron, resistant to damage.",
    "Rarity": "Rare",
    "DamageType": "Melee",
    "Effect": "High defense",
    "Hitpoint": 350,
    "Damage": 40,
    "CostElixir": 8,
    "DamageTower": 15,
    "Speed": 1,
    "Price": 450
  }
```
- Create Spell
POST: http://localhost:8080/app/game/create-spell
```
{
    "Name": "Invisibility Cloak",
    "Description": "Envelops the caster in an invisible cloak, making them undetectable to enemies.",
    "Area": 0,
    "DamageType": "None",
    "Damage": 0,
    "Duration": 15,
    "Effect": "Invisibility",
    "Price": 300
}
```
- Getting hero and spell
GET: http://localhost:8080/app/game/get-all-heros
GET: http://localhost:8080/app/game/get-all-spells
GET: http://localhost:8080/app/game/hero/6
GET: http://localhost:8080/app/game/spell/4
GET: http://localhost:8080/app/game/get-my-spells
GET: http://localhost:8080/app/game/get-my-heros

GET: http://localhost:8080/app/game/get-all-heros?sortBy=speed&filterName=F&sortOrder=asc&page=1&pageSize=10
- Buy hero and spell
POST:http://localhost:8080/app/game/hero/6
POST:http://localhost:8080/app/game/spell/2
- 
## Key Features and Considerations
- Dependency injection for database interactions.
- RESTful API endpoints for managing spells, heroes, and decks.
- Database migrations for easy setup and versioning
- CRUD operations for spells, heroes, and decks.
