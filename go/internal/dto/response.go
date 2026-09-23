package dto

import "time"

type HandlerResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	ResponseData    any    `json:"responseData"`
}

type VehicleAvailableResponse struct {
	ID          int `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string	`json:"desc"`
	AvailableQuantity    int `json:"availableQuantity"`
	DailyRate   int64		`json:"dailyRate"`
	Category    string	`json:"category"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type RevenueReportResponse struct {
	TotalRevenue int64 `json:"total_revenue"`
	TotalBookings int `json:"total_bookings"`
}

type TopVehicle struct {
	VehicleID     int    `json:"vehicle_id"`
	Name          string `json:"name"`
	TotalBookings int    `json:"total_bookings"`
}