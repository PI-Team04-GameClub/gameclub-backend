package observer

import (
	"fmt"
	"log"
)

type EmailNotifier struct {
	userEmails map[string]string
}

func NewEmailNotifier(userEmails map[string]string) *EmailNotifier {
	return &EmailNotifier{
		userEmails: userEmails,
	}
}

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

func (e *EmailNotifier) sendEmail(to, subject, body string) {
	log.Printf("Sending email to: %s", to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body: %s", body)
	log.Println("Email sent successfully")
}

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
