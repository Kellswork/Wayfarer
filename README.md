# Wayfarer - Golang

## Project Overview:

This project implements a public bus transportation booking API server.

## Technologies Used:

### Language

- Golang

### Database

- PostgreSQL

## Features

### Implemented:

#### User functionality:

- Sign up
- Login
- View all trips

#### Admin functionality:
- Login
- Add bus for a Trip
- See all buses stored in the db
- Create a trip
- Cancel a trip
- View all trips

### Planned

#### Admin functionality:

- View bookings for all trips

#### User functionality:

- Book a seat on a trip
- View bookings
- Delete their bookings
- Filter trips by origin
- Filter trips by destination
- Specify seat numbers when making a booking

### Installation and Setup

- Clone the repository `https://github.com/Kellswork/Wayfarer.git`.
- Navigate to the project directory.
- cd into the folder
- Create a `.env` file with the required configurations.
- In your terminal, run `go get` to install dependencies.
- Run `go run cmd/main.go` to get the server running on your local machine.

### Run Migrations
`goose postgres "<db-connection-string>" up`
