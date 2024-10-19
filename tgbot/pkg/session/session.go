package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type Session struct {
	UserID    string
	State     string
	ExpiresAt time.Time
}

type SessionRepository struct {
	db *badger.DB
}

func NewSessionRepository(db *badger.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession creates a new session in the database
func (r *SessionRepository) CreateSession(session *Session) error {
	return r.db.Update(func(txn *badger.Txn) error {
		key := []byte("session_" + session.UserID)
		data, err := json.Marshal(session)
		if err != nil {
			return fmt.Errorf("failed to marshal session: %w", err)
		}

		err = txn.Set(key, data)
		if err != nil {
			return fmt.Errorf("failed to set session in db: %w", err)
		}

		return nil
	})
}

// DeleteSession deletes an existing session from the database
func (r *SessionRepository) DeleteSession(userID string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		key := []byte("session_" + userID)

		err := txn.Delete(key)
		if err != nil {
			return fmt.Errorf("failed to delete session from db: %w", err)
		}

		return nil
	})
}

// GetSession retrieves a session by userID from the database
func (r *SessionRepository) GetSession(userID string) (*Session, error) {
	var session Session
	err := r.db.View(func(txn *badger.Txn) error {
		key := []byte("session_" + userID)
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return fmt.Errorf("session not found")
			}
			return fmt.Errorf("failed to get session: %w", err)
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &session)
		})
		if err != nil {
			return fmt.Errorf("failed to unmarshal session: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ChangeUserState updates the state of the session for a given user
func (r *SessionRepository) ChangeUserState(userID string, newState string) error {
	session, err := r.GetSession(userID)
	if err != nil {
		return err
	}

	session.State = newState
	return r.CreateSession(session) // Re-save the updated session
}
