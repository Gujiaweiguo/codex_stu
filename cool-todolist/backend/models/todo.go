package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Completed bool      `json:"completed" gorm:"default:false"`
	Order     int       `json:"order" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoService struct {
	DB *gorm.DB
}

func NewTodoService(db *gorm.DB) *TodoService {
	return &TodoService{DB: db}
}

func (s *TodoService) GetAll() ([]Todo, error) {
	var todos []Todo
	err := s.DB.Order("`order` ASC").Find(&todos).Error
	return todos, err
}

func (s *TodoService) GetByID(id uint) (*Todo, error) {
	var todo Todo
	err := s.DB.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *TodoService) Create(title string, order int) (*Todo, error) {
	todo := &Todo{
		Title:     title,
		Completed: false,
		Order:     order,
	}
	err := s.DB.Create(todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) Update(todo *Todo) (*Todo, error) {
	err := s.DB.Save(todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) Toggle(id uint) (*Todo, error) {
	var todo Todo
	err := s.DB.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	todo.Completed = !todo.Completed
	err = s.DB.Save(&todo).Error
	return &todo, err
}

func (s *TodoService) Delete(id uint) error {
	return s.DB.Delete(&Todo{}, id).Error
}

func (s *TodoService) Reorder(ids []uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			err := tx.Model(&Todo{}).Where("id = ?", id).Update("order", i).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
