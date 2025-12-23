package repo

import (
	"sync"

	"example.com/prac11TIP/internal/core"
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
	r.mu.Lock()
	defer r.mu.Unlock()
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
