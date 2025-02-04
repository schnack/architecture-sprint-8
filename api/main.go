package main

import (
	"context"
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization")

	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>Главная страница</title></head><body>Hello</body></html`))
	w.WriteHeader(http.StatusOK)
}

type Report struct {
	Id      string `json:"id"`
	Sensor1 int    `json:"sensor1"`
	Sensor2 int    `json:"sensor2"`
	Sensor3 int    `json:"sensor3"`
}

func reports(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization")
	if r.Method == "OPTIONS" {
		return
	}

	parts := strings.Split(r.Header.Get("Authorization"), " ")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO Првоерять роль!!!
	t, err := verifier.Verify(context.Background(), parts[1])
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := new(JwtJson)
	if err := t.Claims(user); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isAccess := false
	for _, role := range user.RealmAccess.Roles {
		if role == "prothetic_user" {
			isAccess = true
			break
		}
	}
	if isAccess == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	reports := []Report{
		{
			Id:      uuid.New().String(),
			Sensor1: rand.Intn(100),
			Sensor2: rand.Intn(200),
			Sensor3: rand.Intn(300),
		},
	}

	if err := json.NewEncoder(w).Encode(reports); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

var verifier *oidc.IDTokenVerifier

func main() {

	configURL := "http://keycloak:8080/realms/reports-realm"

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}

	// SkipIssuerCheck = true, потому что для обращения к keycloak внутри контейнера не можем использовать localhost
	verifier = provider.Verifier(&oidc.Config{ClientID: "reports-api", SkipClientIDCheck: true, SkipIssuerCheck: true})

	slog.InfoContext(ctx, "Starting API")

	http.HandleFunc("/", index)
	http.HandleFunc("/reports", reports)

	log.Fatal(http.ListenAndServe(":8000", nil))

}

type JwtJson struct {
	Exp            int      `json:"exp"`
	Iat            int      `json:"iat"`
	AuthTime       int      `json:"auth_time"`
	Jti            string   `json:"jti"`
	Iss            string   `json:"iss"`
	Sub            string   `json:"sub"`
	Typ            string   `json:"typ"`
	Azp            string   `json:"azp"`
	Nonce          string   `json:"nonce"`
	SessionState   string   `json:"session_state"`
	Acr            string   `json:"acr"`
	AllowedOrigins []string `json:"allowed-origins"`
	RealmAccess    struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	Scope             string `json:"scope"`
	Sid               string `json:"sid"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
	Header            struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
		Kid string `json:"kid"`
	} `json:"header"`
}
