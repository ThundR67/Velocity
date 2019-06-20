package handler

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
)

func generateRandomString(length int) string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = config.UserDataConfigUserIDCharset[seededRand.Intn(len(config.UserDataConfigUserIDCharset))]
	}
	return string(b)
}

//checks if response has given access and refresh token and not thrown error
func testForAccessAndRefresh(body []byte, assert *assert.Assertions) (string, string) {
	var data map[string]string
	json.Unmarshal(body, &data)

	val, isErr := data[config.AuthServerConfigErrField]
	assert.Falsef(isErr, "Page Returned Error %s", val)

	access, isAccessToken := data[config.AuthServerConfigAccessTokenField]
	refresh, isRefreshToken := data[config.AuthServerConfigRefreshTokenField]
	assert.True(isAccessToken, "Sign Up Did Not Return Access Token")
	assert.True(isRefreshToken, "Sign Up Did Not Return Refresh Token")
	return access, refresh
}

//used to test sign up handler
func testSignUp(
	main config.UserMain, extra config.UserExtra, handler Handler, assert *assert.Assertions) {

	mainQuery, _ := query.Values(main)
	extraQuery, _ := query.Values(extra)
	url := "/sign-up?" + mainQuery.Encode() + "&" + extraQuery.Encode() + "&scopes=read"

	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(err, "Http new request returned error")

	rr := httptest.NewRecorder()
	handleFunc := http.HandlerFunc(handler.SignUpHandler)
	handleFunc.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code, "Sign Up should return status code 200")
	testForAccessAndRefresh(rr.Body.Bytes(), assert)
}

//used to test sign in handler
func testSignIn(
	username, password string, handler Handler, assert *assert.Assertions) (string, string) {

	user := config.UserMain{Username: username, Password: password}
	query, _ := query.Values(user)
	url := "/sign-in?" + query.Encode() + "&scopes=read"

	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(err, "Http new request returned error")

	rr := httptest.NewRecorder()
	handleFunc := http.HandlerFunc(handler.SignInHandler)
	handleFunc.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code, "Sign In should return status code 200")
	return testForAccessAndRefresh(rr.Body.Bytes(), assert)
}

//used to test sign in fresh
func testRefreshTokens(refreshToken string, handler Handler, assert *assert.Assertions) {
	url := "/refresh?" + config.AuthServerConfigRefreshTokenField + "=" + refreshToken
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(err, "Http new request returned error")
	rr := httptest.NewRecorder()
	handleFunc := http.HandlerFunc(handler.RefreshHandler)
	handleFunc.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code, "Refresh should return status code 200")
	testForAccessAndRefresh(rr.Body.Bytes(), assert)
}

//used to test sign in fresh handler
func testSignInFresh(
	username, password string, handler Handler, assert *assert.Assertions) {

	user := config.UserMain{Username: username, Password: password}
	query, _ := query.Values(user)
	url := "/sign-in-fresh?" + query.Encode()

	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(err, "Http new request returned error")

	rr := httptest.NewRecorder()
	handleFunc := http.HandlerFunc(handler.SignInFreshHandler)
	handleFunc.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code, "Sign In Fresh should return status code 200")
	var data map[string]string
	json.Unmarshal(rr.Body.Bytes(), &data)
	_, isAccessToken := data[config.AuthServerConfigAccessTokenField]
	assert.True(isAccessToken, "Sign In Fresh Did Not Return Fresh Access Token")
}

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	handler := Handler{}
	handler.Init()

	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)

	main := config.UserMain{
		Username: mockUsername,
		Password: mockPassword,
		Email:    mockUsername + "@gmail.com",
	}

	extra := config.UserExtra{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: int64(864466669), //A Timestamp of year 1997
	}

	testSignUp(main, extra, handler, assert)
	_, refresh := testSignIn(main.Username, main.Password, handler, assert)
	testRefreshTokens(refresh, handler, assert)
	testSignInFresh(main.Username, main.Password, handler, assert)
}
