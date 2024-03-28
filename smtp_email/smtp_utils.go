package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

func ConnectSimpleSMTP(host, port, username, password string) (*smtp.Client, error) {
	servername := fmt.Sprintf("%v:%v", host, port)
	// Dial the SMTP server with the timeout
	c, err := smtp.Dial(servername)
	if err != nil {
		fmt.Println("Failed to connect to SMTP server:", err)
		return nil, err
	}

	return c, nil
}

func ConnectSecureSMTP(host, port, username, password string) (*smtp.Client, error) {
	servername := fmt.Sprintf("%v:%v", host, port)
	hostPort, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", username, password, hostPort)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: SMTP_CONFIG.InsecureSkipVerify,
		ServerName:         host,
	}

	// TODO: set timeout
	dialer := &net.Dialer{
		Timeout: 15 * time.Second, // 30-second timeout
	}

	// Dial and establish TLS connection
	conn, err := tls.DialWithDialer(dialer, "tcp", servername, tlsconfig)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Connection timed out")
			errInfo := errors.New("connection timed out")
			return nil, errInfo
		} else {
			fmt.Println("Error connecting to SMTP server:", err)
			errInfo := errors.New("error connecting to SMTP server")
			// return nil, errInfo
			_ = errInfo
			return nil, err
		}
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}

	if SMTP_CONFIG.StartTLS {
		if err := c.StartTLS(tlsconfig); err != nil {
			return nil, err
		}
	}
	if SMTP_CONFIG.IsAuthRequired {
		// Auth
		if err = c.Auth(auth); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func sendEmailHelper(conn *smtp.Client, from, to, msg string) error {
	if err := conn.Mail(from); err != nil {
		return err
	}
	if err := conn.Rcpt(to); err != nil {
		return err
	}

	wc, err := conn.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}

func SendEmail(from, to, subject, msg string) error {
	host := SMTP_CONFIG.Host
	port := SMTP_CONFIG.Port
	username := SMTP_CONFIG.Username
	password := SMTP_CONFIG.Password

	// c, err := ConnectSMTP(host, port, username, password)
	var err error
	var c *smtp.Client

	if SMTP_CONFIG.InsecureSkipVerify && SMTP_CONFIG.IsAuthRequired {
		c, err = ConnectSecureSMTP(host, port, username, password)
		if err != nil {
			return err
		}
		defer c.Close()
	} else {
		c, err = ConnectSimpleSMTP(host, port, username, password)
		if err != nil {
			return err
		}
		defer c.Close()
	}

	// Setup headers
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + msg

	err = sendEmailHelper(c, from, to, string(message))
	if err != nil {
		return err
	}

	return err
}
