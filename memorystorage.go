package cookiejarx

import (
	"sort"
	"sync"
	"time"
)

type inMemoryEntry struct {
	*Entry

	// seqNum is a sequence number so that Cookies returns cookies in a
	// deterministic order, even for cookies that have equal Path length and
	// equal Creation time. This simplifies testing.
	seqNum uint64
}

// InMemoryStorage provides thread-safe in-memory entry storage with predictable entry sorting
type InMemoryStorage struct {
	// mu locks the remaining fields.
	mu sync.Mutex

	// entries is a set of entries, keyed by their eTLD+1 and subkeyed by
	// their name/domain/path.
	entries map[string]map[string]inMemoryEntry

	// nextSeqNum is the next sequence number assigned to a new cookie
	// created SetCookies.
	nextSeqNum uint64
}

// NewInMemoryStorage returns new InMemoryStorage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		entries: make(map[string]map[string]inMemoryEntry),
	}
}

// EntriesDump returns all entries persisted in in-memory storage
func (s *InMemoryStorage) EntriesDump() (entries []*Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, submap := range s.entries {
		for _, e := range submap {
			entries = append(entries, e.Entry)
		}
	}

	return entries
}

// EntriesRestore adds provide entries to current in-memory storage
func (s *InMemoryStorage) EntriesRestore(entries []*Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range entries {
		s.saveEntry(e)
	}
}

// EntriesClear empties current in-memory storage
func (s *InMemoryStorage) EntriesClear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make(map[string]map[string]inMemoryEntry)
}

// SaveEntry in-memory implementation of Storage.SaveEntry
func (s *InMemoryStorage) SaveEntry(entry *Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.saveEntry(entry)
}

func (s *InMemoryStorage) saveEntry(entry *Entry) {
	submap := s.entries[entry.Key]

	if submap == nil {
		submap = make(map[string]inMemoryEntry)
	}

	e := inMemoryEntry{
		Entry: entry,
	}

	id := entry.ID

	if old, ok := submap[id]; ok {
		e.Creation = old.Creation
		e.seqNum = old.seqNum
	} else {
		e.seqNum = s.nextSeqNum
		s.nextSeqNum++
	}

	submap[id] = e

	s.entries[entry.Key] = submap
}

// RemoveEntry in-memory implementation of Storage.RemoveEntry
func (s *InMemoryStorage) RemoveEntry(key, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	submap := s.entries[key]

	var modified bool

	if submap != nil {
		if _, ok := submap[id]; ok {
			delete(submap, id)
			modified = true
		}
	}

	if modified && len(submap) == 0 {
		delete(s.entries, key)
	}
}


// Entries in-memory implementation of Storage.Entries
func (s *InMemoryStorage) Entries(https bool, host, path, key string, now time.Time) (entries []*Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	submap := s.entries[key]
	if submap == nil {
		return entries
	}

	modified := false
	var selected []inMemoryEntry
	for id, e := range submap {
		if e.Persistent && !e.Expires.After(now) {
			delete(submap, id)
			modified = true
			continue
		}

		if !e.ShouldSend(https, host, path) {
			continue
		}
		e.LastAccess = now
		submap[id] = e
		selected = append(selected, e)
		modified = true
	}
	if modified {
		if len(submap) == 0 {
			delete(s.entries, key)
		} else {
			s.entries[key] = submap
		}
	}

	// sort according to RFC 6265 section 5.4 point 2: by longest
	// path and then by earliest creation time.
	sort.Slice(selected, func(i, j int) bool {
		sel := selected
		if len(sel[i].Path) != len(sel[j].Path) {
			return len(sel[i].Path) > len(sel[j].Path)
		}
		if !sel[i].Creation.Equal(sel[j].Creation) {
			return sel[i].Creation.Before(sel[j].Creation)
		}
		return sel[i].seqNum < sel[j].seqNum
	})

	if len(selected) == 0 {
		return entries
	}

	entries = make([]*Entry, len(selected))

	for i, sel := range selected {
		entries[i] = sel.Entry
	}

	return entries
}
