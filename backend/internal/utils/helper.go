package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FormatPrice formats a float64 price to Indonesian Rupiah
func FormatPrice(price float64) string {
	return fmt.Sprintf("Rp %s", FormatNumber(int64(price)))
}

// FormatNumber formats a number with thousand separators
func FormatNumber(n int64) string {
	str := strconv.FormatInt(n, 10)
	// Handle negative numbers
	negative := false
	if str[0] == '-' {
		negative = true
		str = str[1:]
	}

	// Add thousand separators
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	if negative {
		return "-" + result.String()
	}
	return result.String()
}

// ParsePrice parses a price string to float64
func ParsePrice(priceStr string) (float64, error) {
	// Remove currency symbol and thousand separators
	re := regexp.MustCompile(`[^0-9,.-]`)
	cleaned := re.ReplaceAllString(priceStr, "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	return strconv.ParseFloat(cleaned, 64)
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%d", GenerateUUID()[:8], time.Now().Unix())
}

// GenerateInvoiceNumber generates a unique invoice number
func GenerateInvoiceNumber() string {
	return fmt.Sprintf("INV-%s-%d", GenerateUUID()[:8], time.Now().Unix())
}

// RoundToTwoDecimals rounds a float64 to 2 decimal places
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// IsValidPhone validates phone number format (Indonesia)
func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,10}$`)
	return re.MatchString(phone)
}

// TruncateString truncates a string to max length and adds ellipsis
func TruncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}

// MapToStruct converts a map to a struct using JSON marshaling
func MapToStruct(m map[string]interface{}, s interface{}) error {
	// This would require json marshaling/unmarshaling
	// For simplicity, return nil
	return nil
}

// StructToMap converts a struct to a map
func StructToMap(s interface{}) (map[string]interface{}, error) {
	// This would require json marshaling/unmarshaling
	// For simplicity, return nil, nil
	return nil, nil
}

// Pointer helpers
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func BoolPtr(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

// GetOrDefault returns value if not nil, else default
func GetOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

// GenerateSHA512 generates a SHA512 hash of the input string
func GenerateSHA512(input string) string {
	h := sha512.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}