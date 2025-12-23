# Практическое занятие №11. Проектирование REST API (CRUD для заметок). Разработка структуры. Колчин Степан Сергеевич. ЭФМО-02-25

----
**1.	Цель и краткое описание задания.**

- Освоить принципы проектирования REST API

- Научиться разрабатывать структуру проекта backend-приложения на Go

- Освоить применение слоистой архитектуры (handler → service → repository)

- Подготовить основу для интеграции с базой данных и JWT-аутентификацией в следующих занятиях
  

**2.	Теоретические положения REST API и CRUD.**

- `REST (Representational State Transfer)` — это архитектурный стиль взаимодействия клиентских и серверных приложений через протокол HTTP.

`CRUD` и соответствие с `HTTP`

- Create - POST

- Read - GET

- Update - PATCH/PUT

- Delete - DELETE


**3.	Структуру созданного проекта (в виде дерева каталогов).**

```
011-practice/
├── api/
│   └── openapi.yaml
├── cmd/
│   └── api/
│       └── main.go                 
├── internal/
│   ├── core/
│   │   └── note.go                 
│   ├── http/
│   │   ├── handlers/
│   │   │   └── notes.go            
│   │   └── router.go                
│   └── repo/
│       └── note_mem.go              
├── go.mod                          
├── go.sum
└── README.md      
```

**4.	Примеры кода основных файлов (main.go, note_mem.go, handlers/notes.go).**

- [main.go](./cmd/api/main.go)
```go
package main

import (
    "log"
    "net/http"
    "example.com/notes-api/internal/http"     
    "example.com/notes-api/internal/http/handlers" 
    "example.com/notes-api/internal/repo" 
)

func main() {
  repo := repo.NewNoteRepoMem()
  h := &handlers.Handler{Repo: repo}
  r := httpx.NewRouter(h)

  log.Println("Server started at :8080")
  log.Fatal(http.ListenAndServe(":8080", r))
}

```

- [note_mem.go](./internal/repo/note_mem.go)
```go
package repo

import (
  "sync"
  "example.com/notes-api/internal/core"
)

type NoteRepoMem struct {
  mu    sync.Mutex
  notes map[int64]*core.Note
  next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
  return &NoteRepoMem{notes: make(map[int64]*core.Note)}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
  r.mu.Lock(); defer r.mu.Unlock()
  r.next++
  n.ID = r.next
  r.notes[n.ID] = &n
  return n.ID, nil
}


func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    note, exists := r.notes[id]
    if !exists {
        return nil, nil 
    }
    return note, nil
}

func (r *NoteRepoMem) GetAll() ([]*core.Note, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    notes := make([]*core.Note, 0, len(r.notes))
    for _, note := range r.notes {
        notes = append(notes, note)
    }
    return notes, nil
}

func (r *NoteRepoMem) Update(id int64, updated core.Note) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.notes[id]; !exists {
        return nil 
    }
    
    updated.ID = id
    r.notes[id] = &updated
    return nil
}

func (r *NoteRepoMem) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.notes[id]; !exists {
        return nil 
    }
    
    delete(r.notes, id)
    return nil
}
```

[notes.go](./internal/http/handlers/notes.go)
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/go-chi/chi/v5"
    "example.com/notes-api/internal/core"  
    "example.com/notes-api/internal/repo"
)

type Handler struct {
  Repo *repo.NoteRepoMem
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
  var n core.Note
  if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
    http.Error(w, "Invalid input", http.StatusBadRequest)
    return
  }
  id, _ := h.Repo.Create(n)
  n.ID = id
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(n)
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
    notes, _ := h.Repo.GetAll()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notes)
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    note, _ := h.Repo.GetByID(id)
    if note == nil {
        http.Error(w, "Note not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(note)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var n core.Note
    if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    
    // Проверяем существование
    existing, _ := h.Repo.GetByID(id)
    if existing == nil {
        http.Error(w, "Note not found", http.StatusNotFound)
        return
    }
    
    _ = h.Repo.Update(id, n)
    n.ID = id
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(n)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    existing, _ := h.Repo.GetByID(id)
    if existing == nil {
        http.Error(w, "Note not found", http.StatusNotFound)
        return
    }
    
    _ = h.Repo.Delete(id)
    w.WriteHeader(http.StatusNoContent)
}
```
5.	Скриншоты работы API (Postman или curl).

- `Create - POST`

<img width="561" height="80" alt="image" src="https://github.com/user-attachments/assets/e7327fb3-d2af-4b27-8cf3-569953e4498f" />


- `Read - GET`

<img width="641" height="60" alt="image" src="https://github.com/user-attachments/assets/cc16254a-676c-48ce-97ce-145481540cea" />


- `Read - GET (GET by ID)`

<img width="625" height="56" alt="image" src="https://github.com/user-attachments/assets/b8665434-e379-4bfc-8401-1bce7c716c1a" />

- `Update - PUT`

<img width="778" height="67" alt="image" src="https://github.com/user-attachments/assets/d88b1418-08f0-467c-858b-8f71eb874715" />

- `Delete - delete`

<img width="539" height="89" alt="image" src="https://github.com/user-attachments/assets/54bc6c63-ad84-498f-b3e7-a057d0945471" />


6.	Выводы о проделанной работе.

-   Цели работы были выполнены:

    - Освоить принципы проектирования REST API

    - Научиться разрабатывать структуру проекта backend-приложения на Go

    - Освоить применение слоистой архитектуры (handler → service → repository)

- `API` работает по стандартам `REST` - используется HTTP методы: `POST`, `GET`, `PUT`, `DELETE`

- Маршрутизация через `Chi` - настроен `RESTful` endpoints с версионированием `API` (/api/v1/)
