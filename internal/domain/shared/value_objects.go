package shared

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ID represents a domain identifier
type ID struct {
	value uuid.UUID
}

// NewID creates a new ID
func NewID() ID {
	return ID{value: uuid.New()}
}

// NewIDFromString creates an ID from a string
func NewIDFromString(s string) (ID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID{value: u}, nil
}

// String returns the string representation of the ID
func (id ID) String() string {
	return id.value.String()
}

// Value returns the underlying UUID
func (id ID) Value() uuid.UUID {
	return id.value
}

// IsZero checks if the ID is zero value
func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}

// Email represents an email address
type Email struct {
	value string
}

// NewEmail creates a new email value object
func NewEmail(email string) (Email, error) {
	// Basic email validation
	if email == "" {
		return Email{}, NewDomainError(ErrCodeValidation, "email cannot be empty", "")
	}
	// Add more sophisticated email validation here
	return Email{value: email}, nil
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// Value returns the underlying email string
func (e Email) Value() string {
	return e.value
}

// Timestamp represents a domain timestamp
type Timestamp struct {
	value time.Time
}

// NewTimestamp creates a new timestamp
func NewTimestamp() Timestamp {
	return Timestamp{value: time.Now()}
}

// NewTimestampFromTime creates a timestamp from a time.Time
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp{value: t}
}

// Time returns the underlying time.Time
func (t Timestamp) Time() time.Time {
	return t.value
}

// Unix returns the Unix timestamp
func (t Timestamp) Unix() int64 {
	return t.value.Unix()
}

// String returns the RFC3339 string representation
func (t Timestamp) String() string {
	return t.value.Format(time.RFC3339)
}

// IsZero checks if the timestamp is zero value
func (t Timestamp) IsZero() bool {
	return t.value.IsZero()
}

// Money represents a monetary value
type Money struct {
	amount   int64  // Amount in cents
	currency string // Currency code (e.g., "USD", "EUR")
}

// NewMoney creates a new money value object
func NewMoney(amount int64, currency string) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}

// Amount returns the amount in cents
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// Add adds another money amount (must be same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, NewDomainError(ErrCodeValidation, "cannot add different currencies", "")
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// String returns the string representation
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.amount)/100, m.currency)
}
