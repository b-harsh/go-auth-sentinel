
# 🛡️ Go-Auth-Sentinel

> A high-performance, containerized Two-Factor Authentication (2FA) microservice built in Go.
> **Implements RFC 6238 (TOTP) from scratch without external OTP libraries.**

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green)

## 📖 Overview

**Go-Auth-Sentinel** is a backend microservice designed to handle Time-Based One-Time Password (TOTP) verification. 

Unlike standard implementations that rely on pre-built libraries (like `pquerna/otp`), this project manually implements the **HMAC-SHA1 algorithm** and **dynamic binary truncation** logic defined in [RFC 6238](https://tools.ietf.org/html/rfc6238). This approach ensures a deep understanding of cryptographic protocols and bitwise operations.

### Key Features
* **Zero-Dependency Logic:** Core TOTP algorithm written purely in Go's standard library (`crypto/hmac`, `crypto/sha1`, `encoding/binary`).
* **Containerized:** Fully Dockerized for consistent deployment across any environment (Linux/AWS/Local).
* **High Performance:** Built on the **Gin** framework for sub-millisecond API response times.
* **Stateless:** Designed as a RESTful microservice suitable for horizontal scaling.

---

## 🛠️ Tech Stack

* **Language:** Go (Golang) 1.23
* **Framework:** Gin Web Framework
* **Containerization:** Docker (Alpine Linux base)
* **Protocols:** HTTP/1.1, RFC 6238 (TOTP)

---

## 📐 How It Works (The Logic)

This service validates a 6-digit code using the following cryptographic workflow:

1.  **Time Step Calculation:**
    * Takes the current Unix Epoch time.
    * Divides by `30` (the default time-step window) to get a moving counter.
2.  **HMAC-SHA1 Hashing:**
    * Combines the **Secret Key** (Base32 decoded) and the **Time Counter** (converted to Big-Endian bytes).
    * Generates a 160-bit SHA1 hash.
3.  **Dynamic Truncation:**
    * Extracts the last nibble of the hash to determine an `offset`.
    * Reads 4 bytes starting from that `offset`.
    * Performs a bitwise `AND` (`& 0x7fffffff`) to mask the signed bit.
4.  **Modulo Operation:**
    * Calculates `Result % 1,000,000` to generate the final 6-digit PIN.

---

## 🚀 Getting Started

### Prerequisites
* [Go 1.21+](https://go.dev/dl/) (for local dev)
* [Docker](https://www.docker.com/) (for containerized run)

### Option A: Run with Docker (Recommended)
This requires no Go installation on your machine.

1.  **Build the Image:**
    ```bash
    docker build -t go-auth-sentinel .
    ```

2.  **Run the Container:**
    ```bash
    docker run -p 8080:8080 go-auth-sentinel
    ```

### Option B: Run Locally
```bash
go mod download
go run totp.go

```

---

## 🔌 API Documentation

### Verify Token

Validates a TOTP code against a secret key.

* **Endpoint:** `POST /verify`
* **Content-Type:** `application/json`

#### Request Body

```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "token": "123456"
}

```

*Note: You can test this using the **Google Authenticator** app. Add a new account manually using the secret `JBSWY3DPEHPK3PXP` to generate valid tokens.*

#### Response (Success)

```json
{
    "status": "success",
    "message": "Authentication Successful",
    "valid": true
}

```

#### Response (Failure)

```json
{
    "status": "failed",
    "message": "Invalid Token",
    "valid": false
}

```

---

## 📂 Project Structure

```text
go-auth-sentinel/
├── Dockerfile        # Multi-stage Docker build config
├── totp.go           # Entry point and core algorithm logic
├── go.mod            # Go module definitions
├── go.sum            # Checksums for dependencies
└── README.md         # Project documentation

```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
