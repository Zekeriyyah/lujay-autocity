package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendListingStatusEmail sends an HTML email to the seller about the status of their listing.
func SendListingStatusEmail(sellerEmail, status, listingTitle string) error {
	if status != "approved" && status != "rejected" {
		return fmt.Errorf("invalid status: %s, must be 'approved' or 'rejected'", status)
	}

	subject := fmt.Sprintf("Your AutoCity Listing Status: %s", listingTitle)
	from := os.Getenv("SMTP_USERNAME")
	smtpHost := os.Getenv("SMTP_HOST") 
	smtpPort := os.Getenv("SMTP_PORT") 

	var statusMessage, statusExplanation string
	switch status {
	case "approved":
		statusMessage = "Approved!"
		statusExplanation = "Your listing has passed our vetting process and is now visible to potential buyers."
	case "rejected":
		statusMessage = "Not Approved"
		statusExplanation = "Unfortunately, your listing did not meet our current standards and has been rejected."
	}

	// Email headers and body with proper MIME and HTML formatting
	body := fmt.Sprintf(`To: %s
Subject: %s
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

<p>Hello,</p>
<p>Your vehicle listing "<strong>%s</strong>" has been <strong>%s</strong>.</p>
<p>%s</p>
<p>Thanks,<br>The AutoCity Team</p>
`, sellerEmail, subject, listingTitle, statusMessage, statusExplanation)

	// Set up authentication and SMTP address
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return smtp.SendMail(addr, auth, from, []string{sellerEmail}, []byte(body))
}