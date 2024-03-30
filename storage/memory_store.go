package storage

import (
	"errors"
	"sync"
)

type InMemoryUserStore struct {
	mu     sync.RWMutex // guards the following
	users  map[string]User
	nextID int64
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:  make(map[string]User),
		nextID: 1,
	}
}

func (s *InMemoryUserStore) CreateUser(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Email] = user

	return nil
}

func (s *InMemoryUserStore) ReadUser(email string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[email]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

type InMemoryRecipeStore struct {
	mu      sync.RWMutex // guards the following
	recipes map[int64]Recipe
	nextID  int64
}

func NewInMemoryRecipeStore() *InMemoryRecipeStore {
	return &InMemoryRecipeStore{
		recipes: make(map[int64]Recipe),
		nextID:  1,
	}
}

func (s *InMemoryRecipeStore) CreateRecipe(recipe Recipe) (Recipe, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	recipe.ID = s.nextID
	s.nextID++

	s.recipes[recipe.ID] = recipe

	return recipe, nil
}

func (s *InMemoryRecipeStore) ReadRecipe(id int64) (Recipe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	recipe, ok := s.recipes[id]
	if !ok {
		return Recipe{}, errors.New("recipe not found")
	}

	return recipe, nil
}

func (s *InMemoryRecipeStore) ListRecipes() ([]Recipe, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var recipes []Recipe
	for _, recipe := range s.recipes {
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (s *InMemoryRecipeStore) UpdateRecipe(id int64, recipe Recipe) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.recipes[id]
	if !ok {
		return errors.New("recipe not found")
	}

	recipe.ID = id
	s.recipes[id] = recipe

	return nil
}

func (s *InMemoryRecipeStore) DeleteRecipe(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.recipes[id]
	if !ok {
		return errors.New("recipe not found")
	}

	delete(s.recipes, id)

	return nil
}
