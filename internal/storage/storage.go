package storage

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	UserNotFound = errors.New("user_not_found")
)

type StorageInterface interface {
	CreateUser(u UserInput) (string, error)
	GetUserByID(id string) (User, error)
	GetAllUsers() (UserList, error)
	DeleteUser(id string) error
	UpdateUser(id string, u UserInput) error
}

type Storage struct {
	us          UserStore
	storagePath string
	mu          sync.Mutex
}

func New(storagePath string) *Storage {
	return &Storage{us: UserStore{}, storagePath: storagePath, mu: sync.Mutex{}}
}

func (s *Storage) CreateUser(u UserInput) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.ReadFile(s.storagePath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(f, &s.us)
	if err != nil {
		return "", err
	}
	s.us.Increment++
	user := User{
		CreatedAt:   time.Now(),
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
	id := strconv.Itoa(s.us.Increment)
	s.us.List[id] = user

	b, err := json.Marshal(&s.us)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(s.storagePath, b, fs.ModePerm)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Storage) GetUserByID(id string) (User, error) {
	f, err := os.ReadFile(s.storagePath)
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(f, &s.us)
	if err != nil {
		return User{}, err
	}
	_, ok := s.us.List[id]
	if !ok {
		return User{}, UserNotFound
	}

	return s.us.List[id], nil
}

func (s *Storage) GetAllUsers() (UserList, error) {
	f, err := os.ReadFile(s.storagePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &s.us)
	if err != nil {
		return nil, err
	}
	return s.us.List, nil
}

func (s *Storage) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.ReadFile(s.storagePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &s.us)
	if err != nil {
		return err
	}
	_, ok := s.us.List[id]
	if !ok {
		return UserNotFound
	}

	delete(s.us.List, id)
	b, err := json.Marshal(&s.us)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.storagePath, b, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil

}

func (s *Storage) UpdateUser(id string, u UserInput) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.ReadFile(s.storagePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &s.us)
	if err != nil {
		return err
	}
	_, ok := s.us.List[id]
	if !ok {
		return UserNotFound
	}
	usr := s.us.List[id]
	if u.DisplayName != "" {
		usr.DisplayName = u.DisplayName

	}
	if u.Email != "" {

		usr.Email = u.Email
	}

	s.us.List[id] = usr
	b, err := json.Marshal(&s.us)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.storagePath, b, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
