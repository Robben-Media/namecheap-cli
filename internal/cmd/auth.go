package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/builtbyrobben/namecheap-cli/internal/outfmt"
	"github.com/builtbyrobben/namecheap-cli/internal/secrets"
)

type AuthCmd struct {
	SetKey  AuthSetKeyCmd  `cmd:"" name:"set-key" help:"Set API key (uses --stdin by default)"`
	SetUser AuthSetUserCmd `cmd:"" name:"set-user" help:"Set API username"`
	SetIP   AuthSetIPCmd   `cmd:"" name:"set-ip" help:"Set client IP address"`
	Status  AuthStatusCmd  `cmd:"" help:"Show authentication status"`
	Remove  AuthRemoveCmd  `cmd:"" help:"Remove all stored credentials"`
}

// --- set-key ---

type AuthSetKeyCmd struct {
	Stdin bool   `help:"Read API key from stdin (default: true)" default:"true"`
	Key   string `arg:"" optional:"" help:"API key (discouraged; exposes in shell history)"`
}

func (cmd *AuthSetKeyCmd) Run(ctx context.Context) error {
	apiKey, err := readSecret(cmd.Key, "Enter API key: ")
	if err != nil {
		return err
	}

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.SetAPIKey(apiKey); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "API key stored in keyring",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "API key stored in keyring"}})
	}

	fmt.Fprintln(os.Stderr, "API key stored in keyring")

	return nil
}

// --- set-user ---

type AuthSetUserCmd struct {
	Username string `arg:"" required:"" help:"Namecheap API username"`
}

func (cmd *AuthSetUserCmd) Run(ctx context.Context) error {
	username := strings.TrimSpace(cmd.Username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if err := secrets.SetSecret("api_user", []byte(username)); err != nil {
		return fmt.Errorf("store username: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Username stored in keyring",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "Username stored in keyring"}})
	}

	fmt.Fprintln(os.Stderr, "Username stored in keyring")

	return nil
}

// --- set-ip ---

type AuthSetIPCmd struct {
	IP string `arg:"" required:"" help:"Client IP address for API access"`
}

func (cmd *AuthSetIPCmd) Run(ctx context.Context) error {
	ip := strings.TrimSpace(cmd.IP)
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	if err := secrets.SetSecret("client_ip", []byte(ip)); err != nil {
		return fmt.Errorf("store client IP: %w", err)
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "Client IP stored in keyring",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "Client IP stored in keyring"}})
	}

	fmt.Fprintln(os.Stderr, "Client IP stored in keyring")

	return nil
}

// --- status ---

type AuthStatusCmd struct{}

func (cmd *AuthStatusCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	hasKey, _ := store.HasKey()
	hasUser := hasSecret("api_user")
	hasIP := hasSecret("client_ip")

	envKey := os.Getenv("NAMECHEAP_API_KEY")
	envUser := os.Getenv("NAMECHEAP_USER")
	envIP := os.Getenv("NAMECHEAP_CLIENT_IP")

	status := map[string]any{
		"api_key":   credStatus(hasKey, envKey != ""),
		"api_user":  credStatus(hasUser, envUser != ""),
		"client_ip": credStatus(hasIP, envIP != ""),
		"storage":   "keyring",
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, status)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"API_KEY", "API_USER", "CLIENT_IP", "STORAGE"}
		rows := [][]string{{
			credStatus(hasKey, envKey != ""),
			credStatus(hasUser, envUser != ""),
			credStatus(hasIP, envIP != ""),
			"keyring",
		}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Fprintf(os.Stdout, "Storage: keyring\n\n")
	printCredLine("API Key", hasKey, envKey != "", "NAMECHEAP_API_KEY", redactKey(store))
	printCredLine("Username", hasUser, envUser != "", "NAMECHEAP_USER", readSecretValue("api_user"))
	printCredLine("Client IP", hasIP, envIP != "", "NAMECHEAP_CLIENT_IP", readSecretValue("client_ip"))

	if !hasKey && envKey == "" {
		fmt.Fprintln(os.Stderr, "\nRun: namecheap-cli auth set-key --stdin")
	}

	if !hasUser && envUser == "" {
		fmt.Fprintln(os.Stderr, "Run: namecheap-cli auth set-user <username>")
	}

	if !hasIP && envIP == "" {
		fmt.Fprintln(os.Stderr, "Run: namecheap-cli auth set-ip <ip>")
	}

	return nil
}

// --- remove ---

type AuthRemoveCmd struct{}

func (cmd *AuthRemoveCmd) Run(ctx context.Context) error {
	store, err := secrets.OpenDefault()
	if err != nil {
		return fmt.Errorf("open credential store: %w", err)
	}

	if err := store.DeleteAPIKey(); err != nil {
		return fmt.Errorf("remove API key: %w", err)
	}

	_ = secrets.DeleteSecret("api_user")
	_ = secrets.DeleteSecret("client_ip")

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]string{
			"status":  "success",
			"message": "All credentials removed",
		})
	}

	if outfmt.IsPlain(ctx) {
		return outfmt.WritePlain(os.Stdout, []string{"STATUS", "MESSAGE"}, [][]string{{"success", "All credentials removed"}})
	}

	fmt.Fprintln(os.Stderr, "All credentials removed")

	return nil
}

// --- helpers ---

func readSecret(argValue, prompt string) (string, error) {
	if argValue != "" {
		fmt.Fprintln(os.Stderr, "Warning: passing keys as arguments exposes them in shell history. Use --stdin instead.")
		return strings.TrimSpace(argValue), nil
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprint(os.Stderr, prompt)

		byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)

		if err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}

		return strings.TrimSpace(string(byteKey)), nil
	}

	byteKey, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read from stdin: %w", err)
	}

	return strings.TrimSpace(string(byteKey)), nil
}

func hasSecret(key string) bool {
	val, err := secrets.GetSecret(key)
	return err == nil && len(val) > 0
}

func readSecretValue(key string) string {
	val, err := secrets.GetSecret(key)
	if err != nil || len(val) == 0 {
		return ""
	}

	return string(val)
}

func redactKey(store secrets.Store) string {
	key, err := store.GetAPIKey()
	if err != nil || len(key) < 8 {
		return ""
	}

	return key[:4] + "..." + key[len(key)-4:]
}

func credStatus(stored, envOverride bool) string {
	if envOverride {
		return "env"
	}

	if stored {
		return "stored"
	}

	return "missing"
}

func printCredLine(label string, stored, envOverride bool, envVar, value string) {
	prefix := label + ":"

	switch {
	case envOverride:
		fmt.Fprintf(os.Stdout, "%-10s Using %s environment variable\n", prefix, envVar)
	case stored && value != "":
		fmt.Fprintf(os.Stdout, "%-10s %s\n", prefix, value)
	case stored:
		fmt.Fprintf(os.Stdout, "%-10s Stored\n", prefix)
	default:
		fmt.Fprintf(os.Stdout, "%-10s Not configured\n", prefix)
	}
}
