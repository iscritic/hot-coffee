package repository

import (
	"encoding/json"
	"errors"
	"github.com/iscritic/hot-coffee/models"
	"io/ioutil"
	"os"
	"sync"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetAll() ([]*models.Order, error)
	GetByID(id string) (*models.Order, error)
	Update(order *models.Order) error
	Delete(id string) error
}

type orderRepository struct {
	filePath string
	mu       sync.Mutex
}

func NewOrderRepository(filePath string) OrderRepository {
	return &orderRepository{
		filePath: filePath,
	}
}

func (r *orderRepository) readData() ([]*models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var orders []*models.Order

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		// Пустой файл, возвращаем пустой срез
		return orders, nil
	}

	err = json.NewDecoder(file).Decode(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) writeData(orders []*models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.filePath, data, 0644)
}

func (r *orderRepository) Create(order *models.Order) error {
	orders, err := r.readData()
	if err != nil {
		return err
	}
	orders = append(orders, order)
	return r.writeData(orders)
}

func (r *orderRepository) GetAll() ([]*models.Order, error) {
	return r.readData()
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	orders, err := r.readData()
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return nil, errors.New("order not found")
}

func (r *orderRepository) Update(order *models.Order) error {
	orders, err := r.readData()
	if err != nil {
		return err
	}
	for i, o := range orders {
		if o.ID == order.ID {
			orders[i] = order
			return r.writeData(orders)
		}
	}
	return errors.New("order not found")
}

func (r *orderRepository) Delete(id string) error {
	orders, err := r.readData()
	if err != nil {
		return err
	}
	for i, order := range orders {
		if order.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			return r.writeData(orders)
		}
	}
	return errors.New("order not found")
}
