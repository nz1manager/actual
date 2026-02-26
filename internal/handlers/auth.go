package handlers

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "net/http"
    "os"
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "github.com/nz1manager/actual/internal/models"
    "github.com/nz1manager/actual/internal/repository"
    "github.com/nz1manager/actual/internal/utils"
)

type AuthHandler struct {
    userRepo *repository.UserRepository
    oauthConfig *oauth2.Config
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
    return &AuthHandler{
        userRepo: userRepo,
        oauthConfig: &oauth2.Config{
            RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
            ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
            ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
            Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
            Endpoint:     google.Endpoint,
        },
    }
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
    // Generate random state
    b := make([]byte, 32)
    rand.Read(b)
    state := base64.StdEncoding.EncodeToString(b)
    
    // Store state in cookie
    c.SetCookie("oauthstate", state, 3600, "/", "", false, true)
    
    url := h.oauthConfig.AuthCodeURL(state)
    c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
    // Verify state
    state, err := c.Cookie("oauthstate")
    if err != nil || c.Query("state") != state {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state"})
        return
    }

    // Exchange code for token
    token, err := h.oauthConfig.Exchange(c, c.Query("code"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
        return
    }

    // Get user info from Google
    client := h.oauthConfig.Client(c, token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
        return
    }
    defer resp.Body.Close()

    var userInfo struct {
        ID      string `json:"id"`
        Email   string `json:"email"`
        Name    string `json:"name"`
        Picture string `json:"picture"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
        return
    }

    // Find or create user
    user, err := h.userRepo.FindByGoogleID(userInfo.ID)
    if err != nil {
        // Create new user
        user = &models.User{
            GoogleID:  userInfo.ID,
            Email:     userInfo.Email,
            AvatarURL: userInfo.Picture,
        }
        
        if err := h.userRepo.Create(user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
            return
        }
        
        // New user needs to complete profile
        user.IsAdmin = false
    }

    // Generate JWT
    jwtToken, err := utils.GenerateToken(user.ID.String(), user.IsAdmin)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Set cookie
    c.SetCookie("token", jwtToken, 3600*24*7, "/", "", false, true)

    // Check if profile is complete
    if user.FirstName == "" || user.LastName == "" {
        c.JSON(http.StatusOK, gin.H{
            "requires_profile_update": true,
            "token": jwtToken,
            "user": user,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": jwtToken,
        "user": user,
    })
}
