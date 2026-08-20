package email

import (
	"fmt"
	"github.com/wneessen/go-mail"
)

func sendEmail(to string, subject, body string, isHtml bool) {
	fmt.Println("Sending email to " + to)
	from := "info@cron_job.com"
	password := "-"
	m := mail.NewMsg()
	if err := m.From(fmt.Sprintf("\"%s\" <%s>", "cron_job", from)); err != nil {
		fmt.Printf("failed to set From address: %s\n", err)
	}
	if err := m.To(to); err != nil {
		fmt.Printf("failed to set To address: %s\n", err)
	}
	m.SetOrganization("cron_job")
	m.Subject(subject)

	if isHtml {
		m.SetBodyString(mail.TypeTextHTML, body)
	} else {
		m.SetBodyString(mail.TypeTextPlain, body)
	}

	c, err := mail.NewClient("smtp.zoho.com",
		mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(from), mail.WithPassword(password))
	if err != nil {
		fmt.Printf("failed to create mail client: %s\n", err)
	}
	if err := c.DialAndSend(m); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
	}
	fmt.Printf("email sent to %s\n", to)

}

func SendJobFailureEmail(to string, jobId uint, jobName string) {
	sendEmail(to, "job failed", fmt.Sprintf("job id : %v ,job name : %s", jobId, jobName), false)
}

func SendJobSuccessExecutedEmail(to string, jobId uint, jobName string) {
	sendEmail(to, "job executed successful", fmt.Sprintf("job id : %v ,job name : %s", jobId, jobName), false)
}
