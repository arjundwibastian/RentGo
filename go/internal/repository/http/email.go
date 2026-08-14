package http

import (
	"fmt"
	"milestone-02/internal/domain"

	"github.com/resend/resend-go/v2"
)

type SendGridEmailService struct {
	apiKey string
}

func NewSendGridEmailService(apiKey string) *SendGridEmailService {
	return &SendGridEmailService{apiKey: apiKey}
}
func (s *SendGridEmailService) SendBookingConfirmation(booking *domain.Booking, vehicle *domain.Vehicle, user *domain.User) error {
	client := resend.NewClient(s.apiKey)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Booking Confirmed! 🎉</h2>
			<p>Hello <strong>%s</strong>,</p>
			<p>Thank you for your payment. Your booking for the <strong>%s</strong> is now confirmed.</p>
			
			<h3>Booking Details:</h3>
			<ul>
				<li><strong>Start:</strong> %s</li>
				<li><strong>End:</strong> %s</li>
				<li><strong>Total Price:</strong> Rp %d</li>
			</ul>
			<p>Have a safe trip!</p>
		</div>
	`, 
		user.FullName, 
		vehicle.Name, 
		booking.BookingStart.Format("02 Jan 2006 at 15:04"), // Formats to "15 Aug 2026 at 08:00"
		booking.BookingEnd.Format("02 Jan 2006 at 15:04"),
		booking.TotalPrice,
	)

	params := &resend.SendEmailRequest{
        From:    "arjunRental@resend.dev", 
        To:      []string{user.Email},
        Subject: "Booking Confirmed",
        Html:    htmlContent,
    }
    // Send it!
    _, err := client.Emails.Send(params)
    if err != nil {
        return err
    }
    
	return nil
}


func (s *SendGridEmailService) SendBookingCancellation(booking *domain.Booking, user *domain.User) error {
	client := resend.NewClient(s.apiKey)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Booking Cancelled!</h2>
			<p>Hello <strong>%s</strong>,</p>
			<p>Your booking with booking id: %d has been cancelled! </p>
			
			<h3>Booking Details:</h3>
			<ul>
				<li><strong>Start:</strong> %s</li>
				<li><strong>End:</strong> %s</li>
				<li><strong>Total Price:</strong> Rp %d</li>
			</ul>
			<p>your money has been refunded to your balance! Thank you for using our service</p>
		</div>
	`, 
		user.FullName, 
		booking.ID,
		booking.BookingStart.Format("02 Jan 2006 at 15:04"), 
		booking.BookingEnd.Format("02 Jan 2006 at 15:04"),
		booking.TotalPrice,
	)

	params := &resend.SendEmailRequest{
        From:    "arjunRental@resend.dev", 
        To:      []string{user.Email},
        Subject: "Booking Confirmed",
        Html:    htmlContent,
    }
    // Send it!
    _, err := client.Emails.Send(params)
    if err != nil {
        return err
    }
    
	return nil
}