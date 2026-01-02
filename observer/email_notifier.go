package observer

import (
	"fmt"
	"log"
)

// EmailNotifier sends email notifications to users
type EmailNotifier struct {
	userEmails map[string]string // map[email]name
}

// NewEmailNotifier creates a new email notifier with user data
func NewEmailNotifier(userEmails map[string]string) *EmailNotifier {
	return &EmailNotifier{
		userEmails: userEmails,
	}
}

// OnTournamentCreated sends email notification to all users
func (e *EmailNotifier) OnTournamentCreated(tournament TournamentData) {
	for email, name := range e.userEmails {
		e.sendEmail(
			email,
			"New Tournament Created!",
			e.formatEmail(name, tournament),
		)
	}

	log.Printf("Sent tournament creation emails to %d users", len(e.userEmails))
}

// sendEmail simulates sending an email
func (e *EmailNotifier) sendEmail(to, subject, body string) {
	log.Printf("Sending email to: %s", to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body: %s", body)
	log.Println("Email sent successfully")
}

// formatEmail creates the email body
func (e *EmailNotifier) formatEmail(userName string, tournament TournamentData) string {
	return fmt.Sprintf(`
Hi %s,

A new tournament has been created!

Tournament: %s
Prize Pool: $%.2f
Start Date: %s

Don't miss out! Register your team now.

Best regards,
GameClub Team
	`,
		userName,
		tournament.Name,
		tournament.PrizePool,
		tournament.StartDate,
	)
}
