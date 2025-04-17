# 🛠 Go API with PostgreSQL, Gorilla Mux, and Middleware

This is a RESTful API built in Go, using Gorilla Mux for routing and PostgreSQL for data storage. It includes user registration, login, and project CRUD operations.

---

## 📦 Dependencies

- `github.com/gorilla/mux` – Router
- `github.com/joho/godotenv` – Load environment variables from `.env`
- `github.com/justinas/alice` – Middleware chaining
- `github.com/lib/pq` – PostgreSQL driver
- `golang.org/x/crypto/bcrypt` – Password hashing

---

## 📁 Environment Variables

Create a `.env` file in the root of your project with the following:

```env
XATA_PSQL_URL=your_postgres_connection_url
```

## 🏁 Getting Started

### Run the server

```bash
go run main.go

