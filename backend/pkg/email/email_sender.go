package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type EmailSender struct {
	config EmailConfig
}

type EmailData struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

func NewEmailSender(config EmailConfig) *EmailSender {
	return &EmailSender{config: config}
}

func (s *EmailSender) Send(data EmailData) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	var body string
	if data.IsHTML {
		body = data.Body
	} else {
		body = data.Body
	}

	msg := s.buildMessage(data.To, data.Subject, body, data.IsHTML)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	return smtp.SendMail(addr, auth, s.config.From, []string{data.To}, []byte(msg))
}

func (s *EmailSender) buildMessage(to, subject, body string, isHTML bool) string {
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	msg := fmt.Sprintf("From: %s\r\n", s.config.From)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("MIME-Version: 1.0\r\n")
	msg += fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", contentType)
	msg += fmt.Sprintf("\r\n%s\r\n", body)

	return msg
}

// Template methods
func (s *EmailSender) SendWelcomeEmail(to, name string) error {
	subject := "Selamat Datang di ErpCosmetics!"
	body := fmt.Sprintf(`
		<h1>Halo %s!</h1>
		<p>Selamat datang di ErpCosmetics. Kami senang Anda bergabung!</p>
		<p>Nikmati pengalaman belanja produk kosmetik dan herbal terbaik.</p>
		<br>
		<p>Salam,<br>Tim ErpCosmetics</p>
	`, name)

	return s.Send(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailSender) SendOrderConfirmation(to, orderNumber string, totalAmount float64) error {
	subject := fmt.Sprintf("Konfirmasi Pesanan #%s", orderNumber)
	body := fmt.Sprintf(`
		<h1>Pesanan Berhasil!</h1>
		<p>Terima kasih telah berbelanja di ErpCosmetics.</p>
		<p>Detail pesanan:</p>
		<ul>
			<li>Nomor Pesanan: %s</li>
			<li>Total Pembayaran: Rp %.0f</li>
		</ul>
		<p>Pesanan Anda akan segera diproses.</p>
		<br>
		<p>Salam,<br>Tim ErpCosmetics</p>
	`, orderNumber, totalAmount)

	return s.Send(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailSender) SendResetPasswordEmail(to, token string) error {
	subject := "Reset Password ErpCosmetics"
	resetLink := fmt.Sprintf("https://erpcosmetics.com/reset-password?token=%s", token)
	
	body := fmt.Sprintf(`
		<h1>Reset Password</h1>
		<p>Kami menerima permintaan untuk mereset password akun Anda.</p>
		<p>Klik link berikut untuk mereset password Anda:</p>
		<p><a href="%s">%s</a></p>
		<p>Link ini akan kadaluarsa dalam 1 jam.</p>
		<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
		<br>
		<p>Salam,<br>Tim ErpCosmetics</p>
	`, resetLink, resetLink)

	return s.Send(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailSender) SendPaymentConfirmation(to, orderNumber string, amount float64) error {
	subject := fmt.Sprintf("Pembayaran Diterima - Pesanan #%s", orderNumber)
	body := fmt.Sprintf(`
		<h1>Pembayaran Berhasil!</h1>
		<p>Kami telah menerima pembayaran Anda sebesar Rp %.0f untuk pesanan #%s.</p>
		<p>Pesanan Anda sedang kami proses.</p>
		<p>Anda dapat melacak status pesanan di dashboard akun Anda.</p>
		<br>
		<p>Salam,<br>Tim ErpCosmetics</p>
	`, amount, orderNumber)

	return s.Send(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailSender) SendShippingNotification(to, orderNumber, trackingNumber, courier string) error {
	subject := fmt.Sprintf("Pesanan #%s Telah Dikirim", orderNumber)
	body := fmt.Sprintf(`
		<h1>Pesanan Dikirim!</h1>
		<p>Pesanan Anda #%s telah dikirim melalui %s.</p>
		<p>Nomor Resi: <strong>%s</strong></p>
		<p>Anda dapat melacak pengiriman menggunakan nomor resi di website ekspedisi terkait.</p>
		<br>
		<p>Salam,<br>Tim ErpCosmetics</p>
	`, orderNumber, courier, trackingNumber)

	return s.Send(EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailSender) renderTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}