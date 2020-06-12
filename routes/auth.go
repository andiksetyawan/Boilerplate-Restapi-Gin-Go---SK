package routes

import (
    "os"
    "fmt"
    "time"
	"net/http"
    "../models"
    "../config"
	"github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "github.com/danilopolani/gocialite/structs"
)

// Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	// Retrieve provider from route
	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
            "clientID":     os.Getenv("CLIENT_ID_G"),
            "clientSecret": os.Getenv("CLIENT_SECRET_G"),
            "redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github":  []string{},
        "google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(c *gin.Context) {
	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

    var newUser = getOrRegisterUser(provider, user)
    var jwtToken = createToken(&newUser)

    c.JSON(200, gin.H{
        "data": newUser,
        "token" : jwtToken,
        "message": "berhasil login",
    })
}

func getOrRegisterUser(provider string, user *structs.User) models.User{
    var userData models.User

    config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

    if userData.ID == 0 {
        newUser := models.User{
            FullName: user.FullName,
            Email   : user.Email,
            SocialId: user.ID,
            Provider: provider,
            Avatar  : user.Avatar,
        }
        config.DB.Create(&newUser)
        return newUser
    } else {
        return userData
    }
}

func createToken(user *models.User) string {
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "user_role": user.Role,
        "exp": time.Now().AddDate(0, 0, 7).Unix(),
        "iat": time.Now().Unix(),
    })

    // Sign and get the complete encoded token as a string using the secret
    tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        fmt.Println(err)
    }

    return tokenString
}