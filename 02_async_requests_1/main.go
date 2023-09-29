package main

import (
	"errors"
	"log"
	"time"
)

// User struct represents a user in the system.
type User struct {
	Email string
}

// UserRepository interface represents a component responsible for creating user accounts.
type UserRepository interface {
	CreateUserAccount(u User) error
}

// NotificationsClient interface represents a component responsible for sending notifications.
type NotificationsClient interface {
	SendNotification(u User) error
}

// NewsletterClient interface represents a component responsible for adding users to a newsletter.
type NewsletterClient interface {
	AddToNewsletter(u User) error
}

// Handler struct is used to handle the sign-up process by interacting with the corresponding components.
type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

// NewHandler is a constructor for the Handler struct.
func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

// SignUp method is responsible for the sign-up process.
func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	// Asynchronously add the user to the newsletter with retries
	go func() {
		for {
			if err := h.newsletterClient.AddToNewsletter(u); err != nil {
				log.Printf("failed to add user %s to the newsletter: %v", u.Email, err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}()

	// Asynchronously send a notification to the user with retries
	go func() {
		for {
			if err := h.notificationsClient.SendNotification(u); err != nil {
				log.Printf("failed to send notification to user %s: %v", u.Email, err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}()

	return nil
}

// Mock implementations for the interfaces
type MockRepository struct{}

func (m MockRepository) CreateUserAccount(u User) error {
	// Mock implementation
	return nil
}

type MockNotificationsClient struct{}

func (m MockNotificationsClient) SendNotification(u User) error {
	// Mock implementation with simulated failure
	return errors.New("notification failed")
}

type MockNewsletterClient struct{}

func (m MockNewsletterClient) AddToNewsletter(u User) error {
	// Mock implementation with simulated failure
	return errors.New("newsletter subscription failed")
}

func main() {
	handler := NewHandler(MockRepository{}, MockNewsletterClient{}, MockNotificationsClient{})
	user := User{Email: "test@example.com"}

	if err := handler.SignUp(user); err != nil {
		log.Fatalf("Failed to sign up user: %v", err)
	}

	// Keep the main goroutine running to allow asynchronous goroutines to execute.
	time.Sleep(10 * time.Second)
}
