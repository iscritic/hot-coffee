package service

import (
	"errors"
	"github.com/iscritic/hot-coffee/internal/repository"
	"github.com/iscritic/hot-coffee/models"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrders() ([]*models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	// Здесь потребуются menuRepo и inventoryRepo для полной реализации
	// menuRepo repository.MenuRepository
	// inventoryRepo repository.InventoryRepository
}

func NewOrderService(or repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: or,
		// menuRepo: mr,
		// inventoryRepo: ir,
	}
}

func (s *orderService) CreateOrder(order *models.Order) error {
	// Генерируем уникальный ID и метку времени
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now().Format(time.RFC3339)
	order.Status = "open"

	// Валидация элементов заказа
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}
		if item.ProductID == "" {
			return errors.New("product ID is required")
		}
		// Здесь мы бы проверяли, существует ли продукт в меню
		// menuItem, err := s.menuRepo.GetByID(item.ProductID)
		// if err != nil {
		//     return errors.New("product not found in menu")
		// }
		// Дополнительная логика для проверки инвентаря
	}

	// Вычитание количества из инвентаря
	// Это требует inventoryRepo и menuRepo
	// Пока оставим комментарии, где будет эта логика

	// Сохраняем заказ
	err := s.orderRepo.Create(order)
	if err != nil {
		slog.Error("Failed to create order", "error", err)
		return err
	}
	slog.Info("Order created", "orderID", order.ID)
	return nil
}

func (s *orderService) GetOrders() ([]*models.Order, error) {
	return s.orderRepo.GetAll()
}

func (s *orderService) GetOrderByID(id string) (*models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *orderService) UpdateOrder(order *models.Order) error {
	// Дополнительная валидация может быть добавлена здесь
	err := s.orderRepo.Update(order)
	if err != nil {
		slog.Error("Failed to update order", "error", err)
		return err
	}
	slog.Info("Order updated", "orderID", order.ID)
	return nil
}

func (s *orderService) DeleteOrder(id string) error {
	err := s.orderRepo.Delete(id)
	if err != nil {
		slog.Error("Failed to delete order", "error", err)
		return err
	}
	slog.Info("Order deleted", "orderID", id)
	return nil
}

func (s *orderService) CloseOrder(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		slog.Error("Order not found", "orderID", id)
		return err
	}
	if order.Status == "closed" {
		return errors.New("order is already closed")
	}
	order.Status = "closed"
	err = s.orderRepo.Update(order)
	if err != nil {
		slog.Error("Failed to close order", "error", err)
		return err
	}
	slog.Info("Order closed", "orderID", id)
	return nil
}
