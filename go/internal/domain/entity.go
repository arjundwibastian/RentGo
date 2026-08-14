package domain

import "time"

type User struct {
	ID        	int `json:"id" gorm:"primaryKey"`
	Email     	string
	FullName  	string
	Password  	string `json:"-"`
	Balance	  	int
	PhoneNumber	string
	Address 	string
	Role		string	`json:"-"`
	CreatedAt 	time.Time 
	UpdatedAt 	time.Time 
}

type Vehicle struct {
	ID          int	`json:"id" gorm:"primaryKey"`
	Name        string
	Description string
	Quantity    int
	DailyRate	int
	Category	string
	CreatedAt   time.Time `json:"-"`
	UpdatedAt 	time.Time `json:"-"`
}

type Booking struct {
	ID          	int `json:"id" gorm:"primaryKey"`
	UserID   		int
	VehicleID		int
	BookingStart	time.Time
	BookingEnd		time.Time
	TotalPrice		int
	Status			string
	CreatedAt   	time.Time
	UpdatedAt 		time.Time
}

