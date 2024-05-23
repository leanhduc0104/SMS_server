package mail

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
	"vcs_server/cache"
	"vcs_server/controller"
	"vcs_server/helper"
	"vcs_server/service"

	"gopkg.in/gomail.v2"
)

func sendMail(reportData string) error {

	smtpServer := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT")) // Convert smtpPort to int
	emailAddress := os.Getenv("SMTP_USER")
	emailPassword := os.Getenv("SMTP_PASS")
	admin_Email := os.Getenv("ADMIN_EMAIL")
	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", admin_Email) // You can change this to any recipient email address
	m.SetHeader("Subject", "Server Report")
	m.SetBody("text/html", reportData)

	d := gomail.NewDialer(smtpServer, smtpPort, emailAddress, emailPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func getReport(rptime int) (string, error) {
	var (
		serverCache      cache.ServerCache           = cache.NewRedisCache("redis:6379", 0, 120)
		serverService    service.ServerService       = service.NewServerService()
		serverController controller.ServerController = controller.NewServerController(serverService, serverCache)
	)
	report_servers, err := serverController.ReportServerStatus(rptime)
	if err != nil {
		return "", err
	}
	data := struct {
		Servers     []helper.ReportReposne
		ServerCount int
	}{
		Servers:     report_servers,
		ServerCount: len(report_servers),
	}
	return formatEmailBody(data.Servers), nil
}

func SendMail(rptime int) {

	ticker := time.NewTicker(time.Duration(rptime) * time.Hour)
	defer ticker.Stop()

	for {
		// Get the report data
		reportData, err := getReport(rptime)
		if err != nil {
			log.Printf("Error generating report: %v", err)
		}

		// Send the email with the report data
		err = sendMail(reportData)
		if err != nil {
			log.Printf("Error sending email: %v", err)
		} else {
			log.Println("Email sent successfully!")
		}

		<-ticker.C
	}
}

func SendReport(rptime int) error {
	reportData, err := getReport(rptime)
	if err != nil {
		log.Printf("Error generating report: %v", err)

		return err
	}

	// Send the email with the report data
	err = sendMail(reportData)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	} else {
		log.Println("Email sent successfully!")
	}
	return nil
}

func formatEmailBody(servers []helper.ReportReposne) string {
	const emailTemplate = `
    <html>
    <body>
    <h1>Server Status Report</h1>
    <table border="1">
        <tr>
            <th>IP</th>
            <th>Name</th>
            <th>Uptime</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.Ipv4}}</td>
            <td>{{.Name}}</td>
            <td>{{.Uptime}}</td>
        </tr>
        {{end}}
    </table>
    </body>
    </html>
    `

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Fatal("Error creating email template:", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, servers)
	if err != nil {
		log.Fatal("Error executing email template:", err)
	}

	return body.String()
}
