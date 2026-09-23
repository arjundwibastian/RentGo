package db

import (
	"errors"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormUserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &gormUserRepo{db: db}
}

func (r *gormUserRepo) CheckEmailExists(email string) (bool, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}

	}
	return true, nil
}

func (r *gormUserRepo) RegisterUser(user *domain.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *gormUserRepo) ValidateUserLogin(userLogin *dto.LoginRequest) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", userLogin.Email).First(&user).Error
	if err != nil {
		return nil, domain.InvalidEmailorPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		return nil, domain.InvalidEmailorPassword
	}
	return &user, nil
}

func (r *gormUserRepo) GetUserProfile(userID int) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *gormUserRepo) GetUserBalance(userID int) (*int64, error) {
	var balance int64
	var user domain.User
	err := r.db.Model(&user).Where("id = ?", userID).Pluck("balance", &balance).Error
	if err != nil {
		return nil, err
	}
	return &balance, nil
}
func (r *gormUserRepo) UpdateBalance(userID int, updatedBalance int64) (*int64, error) {
	var user domain.User
	err := r.db.Model(&user).Where("id = ?", userID).Update("balance", updatedBalance).Error
	if err != nil{
		return nil, err
	}
	return &updatedBalance, nil
}

func (r *gormUserRepo) AddUserBalanceAtomic(userID int, amount int64) (*int64, error) {
	var user domain.User
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.User{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		user.Balance += amount
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user.Balance, nil
}

func (r *gormUserRepo) GetVehicleByID(vehicleID int) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	err := r.db.Model(&vehicle).Where("id = ?", vehicleID).First(&vehicle).Error
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (r *gormUserRepo) GetVehicles() (*[]domain.Vehicle, error) {
	var vehicles []domain.Vehicle
	err := r.db.Model(&vehicles).Find(&vehicles).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}
	return &vehicles, nil
}
func (r *gormUserRepo)GetVehiclesAvailableByDate(reqStart, reqEnd time.Time) (*[]dto.VehicleAvailableResponse, error) {
	var vehicles []dto.VehicleAvailableResponse
	err := r.db.Table("vehicles").
		Select(`
			vehicles.id, 
			vehicles.name,
			vehicles.description,
			(vehicles.quantity - count(bookings.id)) as available_quantity,
			vehicles.daily_rate,
			vehicles.Category
		`).
		Joins(`
			LEFT JOIN bookings 
			ON vehicles.id = bookings.vehicle_id 
			AND bookings.status = ?
			AND bookings.booking_start < ? 
			AND bookings.booking_end > ?
		`, "confirmed", reqEnd, reqStart).
		Group("vehicles.id, vehicles.name, vehicles.description, vehicles.quantity, vehicles.daily_rate, vehicles.category").
		Scan(&vehicles).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		} else {
			return nil, err
		}
	}
	return &vehicles, nil
}

func (r *gormUserRepo) GetAvailableVehicle(vehicleID int, start, end time.Time) (*int, error) {
	var count int
	err := r.db.Table("bookings").
		Where("vehicle_id = ?", vehicleID).
		Where("status = ?", "confirmed").
		Where("booking_start < ? AND booking_end > ?", end, start).
		Select("Count(bookings.id)").
		Scan(&count).Error
	if err != nil {
		return nil, err
	}
	return &count, nil
}

func (r *gormUserRepo) CreateBooking(booking *domain.Booking) (*domain.Booking, error) {
	err := r.db.Create(booking).Error
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *gormUserRepo) CreateBookingAtomic(booking *domain.Booking) (*domain.Booking, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var vehicle domain.Vehicle
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", booking.VehicleID).First(&vehicle).Error; err != nil {
			return err
		}
		var user domain.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", booking.UserID).First(&user).Error; err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&domain.Booking{}).
			Where("vehicle_id = ? AND status = ?", booking.VehicleID, "confirmed").
			Where("booking_start < ? AND booking_end > ?", booking.BookingEnd, booking.BookingStart).
			Count(&count).Error; err != nil {
			return err
		}
		if vehicle.Quantity-int(count) < 1 {
			return domain.ErrvehicleUnavailable
		}
		if user.Balance < booking.TotalPrice {
			return domain.ErrNotEnoughBalance
		}

		if err := tx.Create(booking).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.User{}).Where("id = ?", user.ID).
			Update("balance", gorm.Expr("balance - ?", booking.TotalPrice)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *gormUserRepo) GetUserBookingHistory(userID int) (*[]domain.Booking, error) {
	var booking []domain.Booking
	err := r.db.Where("user_id = ?", userID).Find(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
func (r *gormUserRepo) GetBookingByID(bookingID int)(*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.Where("id = ?", bookingID).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *gormUserRepo) CheckBookingUser(userID, bookingID int)(bool, error) {
	var booking domain.Booking
	err := r.db.Model(&booking).Where("id = ? AND user_id = ?", bookingID, userID).First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
func (r *gormUserRepo) CancelUserBooking(booking *domain.Booking) (*domain.Booking, error) {
	err := r.db.Model(&booking).Where("id = ? AND user_id = ?", booking.ID, booking.UserID).Update("status", "cancelled").Error
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *gormUserRepo) CancelBookingAtomic(bookingID, userID int) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", bookingID).First(&booking).Error; err != nil {
			return err
		}
		if booking.UserID != userID {
			return domain.ErrInvalidBookingUserID
		}
		if booking.Status != "confirmed" {
			return domain.ErrInvalidBookingCancel
		}
		if err := tx.Model(&domain.Booking{}).Where("id = ?", bookingID).
			Update("status", "cancelled").Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.User{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", booking.TotalPrice)).Error; err != nil {
			return err
		}
		booking.Status = "cancelled"
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func(r *gormUserRepo) CreateNewVehicles(vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	err := r.db.Create(vehicle).Error
	if err != nil {
		return nil, err
	}
	return vehicle, nil
}
func(r *gormUserRepo)UpdateVehicles(vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	result := r.db.Model(vehicle).Updates(vehicle)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}
	return vehicle, nil
}

func(r *gormUserRepo) GetRevenue()(*dto.RevenueReportResponse, error) {
	var report dto.RevenueReportResponse
	err := r.db.Table("bookings").
		Select("COALESCE(SUM(total_price), 0) as total_revenue, COUNT(id) as total_bookings").
		Where("status != ?", "cancelled"). 
		Scan(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
func(r *gormUserRepo) GetTopVehicle()(*[]dto.TopVehicle, error) {
	var topVehicles []dto.TopVehicle
	err := r.db.Table("vehicles").
		Select("vehicles.id as vehicle_id, vehicles.name, COUNT(bookings.id) as total_bookings").
		Joins("JOIN bookings ON vehicles.id = bookings.vehicle_id").
		Where("bookings.status != ?", "cancelled").
		Group("vehicles.id, vehicles.name").
		Order("total_bookings DESC"). 
		Scan(&topVehicles).Error
	if err != nil {
		return nil, err
	}
	return &topVehicles, nil
}