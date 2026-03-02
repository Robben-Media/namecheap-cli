package namecheap

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	productionBaseURL = "https://api.namecheap.com/xml.response"
	sandboxBaseURL    = "https://api.sandbox.namecheap.com/xml.response"
)

var (
	errMissingAPIKey    = errors.New("missing API key")
	errMissingUser      = errors.New("missing API user")
	errMissingIP        = errors.New("missing client IP")
	errDomainRequired   = errors.New("domain name is required")
	errSLDRequired      = errors.New("SLD is required")
	errTLDRequired      = errors.New("TLD is required")
	errNameserversRequired = errors.New("nameservers are required")
	errAPIUnknownStatus    = errors.New("API error: unknown status")
	errAPIResponse         = errors.New("API error")
)

// Client handles Namecheap XML API requests.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiUser    string
	apiKey     string
	userName   string
	clientIP   string
}

// NewClient creates a new Namecheap API client.
func NewClient(apiKey, apiUser, clientIP string, sandbox bool) *Client {
	baseURL := productionBaseURL
	if sandbox {
		baseURL = sandboxBaseURL
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  baseURL,
		apiUser:  apiUser,
		apiKey:   apiKey,
		userName: apiUser,
		clientIP: clientIP,
	}
}

// Validate checks that all required credentials are configured.
func (c *Client) Validate() error {
	if c.apiKey == "" {
		return errMissingAPIKey
	}

	if c.apiUser == "" {
		return errMissingUser
	}

	if c.clientIP == "" {
		return errMissingIP
	}

	return nil
}

// do executes a Namecheap API request and parses the XML response.
func (c *Client) do(ctx context.Context, command string, extraParams map[string]string) (*ApiResponse, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("ApiUser", c.apiUser)
	params.Set("ApiKey", c.apiKey)
	params.Set("UserName", c.userName)
	params.Set("ClientIp", c.clientIP)
	params.Set("Command", command)

	for k, v := range extraParams {
		params.Set(k, v)
	}

	reqURL := c.baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "namecheap-cli/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var apiResp ApiResponse
	if err := xml.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse XML response: %w", err)
	}

	if apiResp.Status != "OK" {
		return &apiResp, c.buildError(&apiResp)
	}

	return &apiResp, nil
}

// buildError creates an error from API error responses.
func (c *Client) buildError(resp *ApiResponse) error {
	if len(resp.Errors.Errors) == 0 {
		return fmt.Errorf("%w: status %s", errAPIUnknownStatus, resp.Status)
	}

	msgs := make([]string, 0, len(resp.Errors.Errors))

	for _, e := range resp.Errors.Errors {
		msgs = append(msgs, fmt.Sprintf("[%s] %s", e.Number, e.Message))
	}

	return fmt.Errorf("%w: %s", errAPIResponse, strings.Join(msgs, "; "))
}

// --- Domain operations ---

// DomainsGetList returns a list of domains.
func (c *Client) DomainsGetList(ctx context.Context, listType string, page, pageSize int) (*ApiResponse, error) {
	params := map[string]string{
		"Page":     fmt.Sprintf("%d", page),
		"PageSize": fmt.Sprintf("%d", pageSize),
	}

	if listType != "" {
		params["ListType"] = listType
	}

	return c.do(ctx, "namecheap.domains.getList", params)
}

// DomainsCheck checks availability of one or more domains.
func (c *Client) DomainsCheck(ctx context.Context, domains string) (*ApiResponse, error) {
	if domains == "" {
		return nil, errDomainRequired
	}

	return c.do(ctx, "namecheap.domains.check", map[string]string{
		"DomainList": domains,
	})
}

// DomainsGetInfo returns detailed info for a domain.
func (c *Client) DomainsGetInfo(ctx context.Context, domain string) (*ApiResponse, error) {
	if domain == "" {
		return nil, errDomainRequired
	}

	return c.do(ctx, "namecheap.domains.getInfo", map[string]string{
		"DomainName": domain,
	})
}

// --- DNS operations ---

// DNSGetHosts returns DNS host records for a domain.
func (c *Client) DNSGetHosts(ctx context.Context, sld, tld string) (*ApiResponse, error) {
	if sld == "" {
		return nil, errSLDRequired
	}

	if tld == "" {
		return nil, errTLDRequired
	}

	return c.do(ctx, "namecheap.domains.dns.getHosts", map[string]string{
		"SLD": sld,
		"TLD": tld,
	})
}

// DNSSetHosts sets DNS host records for a domain.
func (c *Client) DNSSetHosts(ctx context.Context, sld, tld string, records []DNSRecordInput) (*ApiResponse, error) {
	if sld == "" {
		return nil, errSLDRequired
	}

	if tld == "" {
		return nil, errTLDRequired
	}

	params := map[string]string{
		"SLD": sld,
		"TLD": tld,
	}

	for i, rec := range records {
		idx := fmt.Sprintf("%d", i+1)
		params["HostName"+idx] = rec.HostName
		params["RecordType"+idx] = rec.RecordType
		params["Address"+idx] = rec.Address

		if rec.MXPref != "" {
			params["MXPref"+idx] = rec.MXPref
		} else {
			params["MXPref"+idx] = "10"
		}

		if rec.TTL != "" {
			params["TTL"+idx] = rec.TTL
		} else {
			params["TTL"+idx] = "1800"
		}
	}

	return c.do(ctx, "namecheap.domains.dns.setHosts", params)
}

// DNSSetCustom sets custom nameservers for a domain.
func (c *Client) DNSSetCustom(ctx context.Context, sld, tld, nameservers string) (*ApiResponse, error) {
	if sld == "" {
		return nil, errSLDRequired
	}

	if tld == "" {
		return nil, errTLDRequired
	}

	if nameservers == "" {
		return nil, errNameserversRequired
	}

	return c.do(ctx, "namecheap.domains.dns.setCustom", map[string]string{
		"SLD":         sld,
		"TLD":         tld,
		"Nameservers": nameservers,
	})
}

// --- SSL operations ---

// SSLGetList returns a list of SSL certificates.
func (c *Client) SSLGetList(ctx context.Context, listType string) (*ApiResponse, error) {
	params := map[string]string{}

	if listType != "" {
		params["ListType"] = listType
	}

	return c.do(ctx, "namecheap.ssl.getList", params)
}
