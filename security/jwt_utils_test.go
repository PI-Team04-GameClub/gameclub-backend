package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_ValidInput(t *testing.T) {
	// Given: Valid user data
	id := uint(1)
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// When: Generating a token
	token, err := GenerateToken(id, email, firstName, lastName)

	// Then: The token should be generated without error
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_TokenIsValidJWT(t *testing.T) {
	// Given: Valid user data
	id := uint(42)
	email := "user@test.com"
	firstName := "Jane"
	lastName := "Smith"

	// When: Generating a token
	token, err := GenerateToken(id, email, firstName, lastName)

	// Then: The token should be a valid JWT that can be parsed
	assert.NoError(t, err)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestGenerateToken_ContainsCorrectClaims(t *testing.T) {
	// Given: Valid user data
	id := uint(123)
	email := "claims@test.com"
	firstName := "Alice"
	lastName := "Wonder"

	// When: Generating a token and extracting claims
	token, _ := GenerateToken(id, email, firstName, lastName)
	claims, err := ExtractClaims(token)

	// Then: The claims should contain the correct user data
	assert.NoError(t, err)
	assert.Equal(t, float64(id), claims["id"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, firstName, claims["first_name"])
	assert.Equal(t, lastName, claims["last_name"])
}

func TestGenerateToken_ContainsExpirationClaim(t *testing.T) {
	// Given: Valid user data
	id := uint(1)
	email := "exp@test.com"
	firstName := "Bob"
	lastName := "Builder"

	// When: Generating a token and extracting claims
	token, _ := GenerateToken(id, email, firstName, lastName)
	claims, err := ExtractClaims(token)

	// Then: The claims should contain an expiration time in the future
	assert.NoError(t, err)
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	expTime := time.Unix(int64(exp), 0)
	assert.True(t, expTime.After(time.Now()))
}

func TestGenerateToken_ExpiresIn72Hours(t *testing.T) {
	// Given: Valid user data
	id := uint(1)
	email := "72hours@test.com"
	firstName := "Time"
	lastName := "Traveler"

	// When: Generating a token
	beforeGeneration := time.Now()
	token, _ := GenerateToken(id, email, firstName, lastName)
	claims, _ := ExtractClaims(token)

	// Then: The expiration should be approximately 72 hours from now
	exp := claims["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	expectedExp := beforeGeneration.Add(72 * time.Hour)
	diff := expTime.Sub(expectedExp)
	assert.True(t, diff < time.Minute && diff > -time.Minute)
}

func TestGenerateToken_EmptyEmail(t *testing.T) {
	// Given: User data with empty email
	id := uint(1)
	email := ""
	firstName := "Empty"
	lastName := "Email"

	// When: Generating a token
	token, err := GenerateToken(id, email, firstName, lastName)

	// Then: The token should still be generated (no validation on empty values)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_ZeroID(t *testing.T) {
	// Given: User data with zero ID
	id := uint(0)
	email := "zero@test.com"
	firstName := "Zero"
	lastName := "ID"

	// When: Generating a token
	token, err := GenerateToken(id, email, firstName, lastName)

	// Then: The token should still be generated
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestExtractClaims_ValidToken(t *testing.T) {
	// Given: A valid token
	id := uint(99)
	email := "extract@test.com"
	firstName := "Extract"
	lastName := "Claims"
	token, _ := GenerateToken(id, email, firstName, lastName)

	// When: Extracting claims from the token
	claims, err := ExtractClaims(token)

	// Then: The claims should be extracted successfully
	assert.NoError(t, err)
	assert.NotNil(t, claims)
}

func TestExtractClaims_InvalidToken(t *testing.T) {
	// Given: An invalid token string
	invalidToken := "invalid.token.string"

	// When: Trying to extract claims
	claims, err := ExtractClaims(invalidToken)

	// Then: An error should be returned
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestExtractClaims_EmptyToken(t *testing.T) {
	// Given: An empty token string
	emptyToken := ""

	// When: Trying to extract claims
	claims, err := ExtractClaims(emptyToken)

	// Then: An error should be returned
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestExtractClaims_MalformedToken(t *testing.T) {
	// Given: A malformed token (wrong format)
	malformedToken := "not-a-jwt-token"

	// When: Trying to extract claims
	claims, err := ExtractClaims(malformedToken)

	// Then: An error should be returned
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestExtractClaims_WrongSignature(t *testing.T) {
	// Given: A token signed with a different secret
	claims := jwt.MapClaims{
		"id":    1,
		"email": "wrong@secret.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	wrongSecret := []byte("wrong-secret-key")
	tokenString, _ := token.SignedString(wrongSecret)

	// When: Trying to extract claims
	extractedClaims, err := ExtractClaims(tokenString)

	// Then: An error should be returned due to signature mismatch
	assert.Error(t, err)
	assert.Nil(t, extractedClaims)
}

func TestExtractClaims_ExpiredToken(t *testing.T) {
	// Given: An expired token
	claims := jwt.MapClaims{
		"id":    1,
		"email": "expired@test.com",
		"exp":   time.Now().Add(-time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JwtSecret)

	// When: Trying to extract claims
	extractedClaims, err := ExtractClaims(tokenString)

	// Then: An error should be returned due to expiration
	assert.Error(t, err)
	assert.Nil(t, extractedClaims)
}

func TestJwtSecret_IsNotEmpty(t *testing.T) {
	// Given: The JwtSecret variable

	// When: Checking its value

	// Then: It should not be empty
	assert.NotEmpty(t, JwtSecret)
}

func TestGenerateToken_DifferentUsersGetDifferentTokens(t *testing.T) {
	// Given: Two different users
	token1, _ := GenerateToken(1, "user1@test.com", "User", "One")
	token2, _ := GenerateToken(2, "user2@test.com", "User", "Two")

	// When: Comparing the tokens

	// Then: The tokens should be different
	assert.NotEqual(t, token1, token2)
}

func TestGenerateToken_SpecialCharactersInNames(t *testing.T) {
	// Given: User data with special characters
	id := uint(1)
	email := "special@test.com"
	firstName := "José"
	lastName := "O'Brien"

	// When: Generating a token
	token, err := GenerateToken(id, email, firstName, lastName)

	// Then: The token should be generated successfully with special characters preserved
	assert.NoError(t, err)
	claims, _ := ExtractClaims(token)
	assert.Equal(t, firstName, claims["first_name"])
	assert.Equal(t, lastName, claims["last_name"])
}

func TestExtractClaims_ReturnsAllExpectedFields(t *testing.T) {
	// Given: A valid token with all fields
	id := uint(500)
	email := "allfields@test.com"
	firstName := "All"
	lastName := "Fields"
	token, _ := GenerateToken(id, email, firstName, lastName)

	// When: Extracting claims
	claims, err := ExtractClaims(token)

	// Then: All expected fields should be present
	assert.NoError(t, err)
	assert.Contains(t, claims, "id")
	assert.Contains(t, claims, "email")
	assert.Contains(t, claims, "first_name")
	assert.Contains(t, claims, "last_name")
	assert.Contains(t, claims, "exp")
}
