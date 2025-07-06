# 📝 Go Notes API

A simple RESTful Notes API built in Go using MongoDB as the database. This project allows you to create, read, update, and delete notes.

---

## 🚀 Features

- Create a new note
- Get all notes
- Get a note by ID
- Update a note
- Delete a note
- Auto-incrementing `noteId` for easy retrieval

---

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Framework**: net/http, Gorilla Mux
- **Database**: MongoDB (Atlas)
- **ORM**: MongoDB Go Driver

---


## 🧪 API Endpoints

| Method | Endpoint               | Description          |
|--------|------------------------|----------------------|
| POST   | `/notes/create`        | Create a new note    |
| GET    | `/notes`               | Get all notes        |
| GET    | `/notes/{id}`          | Get a note by ID     |
| POST   | `/notes/update/{id}`   | Update a note by ID  |
| DELETE | `/notes/delete/{id}`   | Delete a note by ID  |

---


---

## ⚙️ Environment Setup

Create a `.env` file at the root:

```
PORT=8080
MONGO_URL=mongodb+srv://<your-user>:<your-pass>@cluster.mongodb.net/
DB_NAME=goNotes
```
> ⚠️ Replace with your own MongoDB Atlas URI.

---

## ▶️ Run the App

```bash
# Clone the project
git clone https://github.com/your-username/go-notes.git
cd go-notes

# Install dependencies
go mod tidy

# Run the API
go run cmd/api/main.go

```