package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"os"
)

func GenerateTOTP(secret string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	counter := time.Now().Unix() / 30

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf
	binaryCode := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	otp := binaryCode % 1000000

	return fmt.Sprintf("%06d", otp), nil
}

type VerifyRequest struct {
	Secret string `json:"secret" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

func verifyToken(c *gin.Context) {
	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	expectedToken, err := GenerateTOTP(req.Secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid secret"})
		return
	}

	if req.Token == expectedToken {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"valid":  true,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"valid":  false,
		})
	}
}

func main() {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "online",
            "message": "Go-Auth-Sentinel is running securely 🛡️",
            "usage":   "Send a POST request to /verify with {'secret': '...', 'token': '...'}",
        })
    })
    r.POST("/verify", verifyToken)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Println(">> Go-Auth-Sentinel running on port " + port)
    r.Run(":" + port)
}