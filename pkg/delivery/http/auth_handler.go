package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"forum.bbilisbe/internal/cookies"
	"forum.bbilisbe/internal/models"
)

var (
	RedirectURI = "http://localhost:7070/callback"
	ClientID = ""
	ClientSecret = "-"

	GitClientID = ""
	GitRedirectURI = "http://localhost:7070/github/callback"
	GitClientSecret = ""
)

func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile", ClientID, RedirectURI)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	tokenURL := "https://accounts.google.com/o/oauth2/token"
	clientData := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, ClientID, ClientSecret, RedirectURI)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(clientData))
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		h.serverError(w, err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		h.serverError(w, err)
		return
	}
	accessToken := tokenResponse["access_token"].(string)

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	userInfoResp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		h.serverError(w, err)
		return
	}
	defer userInfoResp.Body.Close()

	form := models.UserSignupForm{}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&form); err != nil {
		http.Error(w, "Failed to decode user info response", http.StatusInternalServerError)
		h.serverError(w, err)
		return
	}

	form.Password, err = generateRandomPassword(8)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		h.serverError(w, err)
		return
	}

	h.loggedinHandler(w, r, form)
}

func (h *Handler) githubLogin(w http.ResponseWriter, r *http.Request) {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		GitClientID, GitRedirectURI,
	)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}


func (h *Handler) githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := getGitHubAccessToken(code)
	githubData := getGitHubData(githubAccessToken)
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	// Set return type JSON
	w.Header().Set("Content-type", "application/json")

	form := models.UserSignupForm{}
	json.Unmarshal([]byte(githubData), &form)

	h.loggedinHandler(w, r, form)
}

func getGitHubAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     GitClientID,
		"client_secret": GitClientSecret,
		"code":          code,
	}

	requestJSON, _ := json.Marshal(requestBodyMap)
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)

	if resperr != nil {
		log.Panic("Request failed by get the response")
	}

	// Response body converted to stringified JSON
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	// Return the access token (as the rest of the
	// details are relatively unnecessary for us)
	return ghresp.AccessToken
}

func getGitHubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed by Make the request")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}

func (h *Handler) loggedinHandler(w http.ResponseWriter, r *http.Request, form models.UserSignupForm) {

	err := h.UUsecase.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) || errors.Is(err, models.ErrDuplicateUsername) {
				userID, err := h.UUsecase.GetUserInfo(form.Email, form.Name)
				// Если пользователь уже существует, создаем сессию и устанавливаем куки
				token := cookies.SetCookie(w, userID)
				err = h.UUsecase.AddToken(userID, token)
				if err != nil {
					h.serverError(w, err)
					return
				}
				// Перенаправляем пользователя на страницу "/post/create"
				http.Redirect(w, r, "/post/create", http.StatusSeeOther)
				return
			} else {
				h.serverError(w, err)
			}
			return
		}
	userID, _ := h.UUsecase.Authenticate(form.Email, form.Password)
	token := cookies.SetCookie(w, userID)
	err = h.UUsecase.AddToken(userID, token)
	if err != nil {
		h.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/post/create", http.StatusSeeOther)
}