package cmd

import (
	"fmt"
	"os"

	"github.com/builtbyrobben/namecheap-cli/internal/namecheap"
	"github.com/builtbyrobben/namecheap-cli/internal/secrets"
)

func getNamecheapClient(sandbox bool) (*namecheap.Client, error) {
	apiKey, err := resolveCredential("NAMECHEAP_API_KEY", "api_key", "API key", "namecheap-cli auth set-key --stdin")
	if err != nil {
		return nil, err
	}

	apiUser, err := resolveCredential("NAMECHEAP_USER", "api_user", "username", "namecheap-cli auth set-user <username>")
	if err != nil {
		return nil, err
	}

	clientIP, err := resolveCredential("NAMECHEAP_CLIENT_IP", "client_ip", "client IP", "namecheap-cli auth set-ip <ip>")
	if err != nil {
		return nil, err
	}

	return namecheap.NewClient(apiKey, apiUser, clientIP, sandbox), nil
}

func resolveCredential(envVar, secretKey, label, hint string) (string, error) {
	if v := os.Getenv(envVar); v != "" {
		return v, nil
	}

	if secretKey == "api_key" {
		store, err := secrets.OpenDefault()
		if err != nil {
			return "", fmt.Errorf("open credential store: %w", err)
		}

		val, err := store.GetAPIKey()
		if err != nil {
			return "", fmt.Errorf("get %s: %w (set %s or run '%s')", label, err, envVar, hint)
		}

		return val, nil
	}

	val, err := secrets.GetSecret(secretKey)
	if err != nil || len(val) == 0 {
		return "", fmt.Errorf("missing %s (set %s or run '%s')", label, envVar, hint)
	}

	return string(val), nil
}
