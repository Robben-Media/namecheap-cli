package namecheap

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseDomainsGetListResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.getList</RequestedCommand>
  <CommandResponse Type="namecheap.domains.getList">
    <DomainGetListResult>
      <Domain ID="123" Name="example.com" User="testuser" Created="05/03/2025" Expires="05/03/2026" IsExpired="false" IsLocked="false" AutoRenew="true" WhoisGuard="ENABLED"/>
      <Domain ID="456" Name="test.net" User="testuser" Created="01/15/2024" Expires="01/15/2025" IsExpired="true" IsLocked="false" AutoRenew="false" WhoisGuard="NOTPRESENT"/>
    </DomainGetListResult>
    <Paging>
      <TotalItems>2</TotalItems>
      <CurrentPage>1</CurrentPage>
      <PageSize>100</PageSize>
    </Paging>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	if resp.Status != "OK" {
		t.Errorf("status = %q, want OK", resp.Status)
	}

	domains := resp.CommandResponse.DomainList.Domains
	if len(domains) != 2 {
		t.Fatalf("got %d domains, want 2", len(domains))
	}

	if domains[0].Name != "example.com" {
		t.Errorf("domain[0].Name = %q, want example.com", domains[0].Name)
	}

	if domains[0].ID != "123" {
		t.Errorf("domain[0].ID = %q, want 123", domains[0].ID)
	}

	if domains[0].AutoRenew != "true" {
		t.Errorf("domain[0].AutoRenew = %q, want true", domains[0].AutoRenew)
	}

	if domains[1].IsExpired != "true" {
		t.Errorf("domain[1].IsExpired = %q, want true", domains[1].IsExpired)
	}

	if resp.CommandResponse.Paging.TotalItems != 2 {
		t.Errorf("paging.TotalItems = %d, want 2", resp.CommandResponse.Paging.TotalItems)
	}
}

func TestParseDomainCheckResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.check</RequestedCommand>
  <CommandResponse Type="namecheap.domains.check">
    <DomainCheckResult Domain="available-domain.com" Available="true"/>
    <DomainCheckResult Domain="google.com" Available="false"/>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	checks := resp.CommandResponse.DomainChecks
	if len(checks) != 2 {
		t.Fatalf("got %d checks, want 2", len(checks))
	}

	if checks[0].Domain != "available-domain.com" {
		t.Errorf("check[0].Domain = %q, want available-domain.com", checks[0].Domain)
	}

	if checks[0].Available != "true" {
		t.Errorf("check[0].Available = %q, want true", checks[0].Available)
	}

	if checks[1].Available != "false" {
		t.Errorf("check[1].Available = %q, want false", checks[1].Available)
	}
}

func TestParseDomainInfoResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.getInfo</RequestedCommand>
  <CommandResponse Type="namecheap.domains.getInfo">
    <DomainGetInfoResult Status="Ok" ID="123" DomainName="example.com" OwnerName="testuser" IsOwner="true">
      <DomainDetails>
        <CreatedDate>05/03/2025</CreatedDate>
        <ExpiredDate>05/03/2026</ExpiredDate>
        <NumYears>1</NumYears>
      </DomainDetails>
      <Whoisguard Enabled="True">
        <ID>12345</ID>
        <ExpiredAt>05/03/2026</ExpiredAt>
      </Whoisguard>
      <DnsDetails ProviderType="CUSTOM" IsUsingOurDNS="false">
        <Nameserver>ns1.example.com</Nameserver>
        <Nameserver>ns2.example.com</Nameserver>
      </DnsDetails>
    </DomainGetInfoResult>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	info := resp.CommandResponse.DomainInfo
	if info.DomainName != "example.com" {
		t.Errorf("DomainName = %q, want example.com", info.DomainName)
	}

	if info.Status != "Ok" {
		t.Errorf("Status = %q, want Ok", info.Status)
	}

	if info.DomainDetails.CreatedDate != "05/03/2025" {
		t.Errorf("CreatedDate = %q, want 05/03/2025", info.DomainDetails.CreatedDate)
	}

	if info.Whoisguard.Enabled != "True" {
		t.Errorf("Whoisguard.Enabled = %q, want True", info.Whoisguard.Enabled)
	}

	if len(info.DNSDetails.Nameservers) != 2 {
		t.Fatalf("got %d nameservers, want 2", len(info.DNSDetails.Nameservers))
	}

	if info.DNSDetails.Nameservers[0] != "ns1.example.com" {
		t.Errorf("nameserver[0] = %q, want ns1.example.com", info.DNSDetails.Nameservers[0])
	}
}

func TestParseDNSHostsResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.dns.getHosts</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.getHosts">
    <DomainDNSGetHostsResult Domain="example.com" IsUsingOurDNS="true">
      <host HostId="1" Name="@" Type="A" Address="1.2.3.4" MXPref="10" TTL="1800" AssociatedAppTitle="" FriendlyName="" IsActive="true" IsDDNSEnabled="false"/>
      <host HostId="2" Name="www" Type="CNAME" Address="example.com." MXPref="10" TTL="1800" AssociatedAppTitle="" FriendlyName="" IsActive="true" IsDDNSEnabled="false"/>
      <host HostId="3" Name="@" Type="MX" Address="mail.example.com." MXPref="5" TTL="1800" AssociatedAppTitle="" FriendlyName="" IsActive="true" IsDDNSEnabled="false"/>
    </DomainDNSGetHostsResult>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	hosts := resp.CommandResponse.DNSHosts
	if hosts.Domain != "example.com" {
		t.Errorf("Domain = %q, want example.com", hosts.Domain)
	}

	if len(hosts.Hosts) != 3 {
		t.Fatalf("got %d hosts, want 3", len(hosts.Hosts))
	}

	if hosts.Hosts[0].Type != "A" {
		t.Errorf("host[0].Type = %q, want A", hosts.Hosts[0].Type)
	}

	if hosts.Hosts[0].Address != "1.2.3.4" {
		t.Errorf("host[0].Address = %q, want 1.2.3.4", hosts.Hosts[0].Address)
	}

	if hosts.Hosts[2].MXPref != "5" {
		t.Errorf("host[2].MXPref = %q, want 5", hosts.Hosts[2].MXPref)
	}
}

func TestParseDNSSetHostsResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.setHosts">
    <DomainDNSSetHostsResult Domain="example.com" IsSuccess="true"/>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	result := resp.CommandResponse.DNSSetResult
	if result.Domain != "example.com" {
		t.Errorf("Domain = %q, want example.com", result.Domain)
	}

	if result.IsSuccess != "true" {
		t.Errorf("IsSuccess = %q, want true", result.IsSuccess)
	}
}

func TestParseDNSSetCustomResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.dns.setCustom</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.setCustom">
    <DomainDNSSetCustomResult Domain="example.com" Updated="true"/>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	result := resp.CommandResponse.DNSSetCustomResult
	if result.Domain != "example.com" {
		t.Errorf("Domain = %q, want example.com", result.Domain)
	}

	if result.Updated != "true" {
		t.Errorf("Updated = %q, want true", result.Updated)
	}
}

func TestDNSSetCustomValidation(t *testing.T) {
	t.Parallel()

	client := NewClient("key", "user", "1.2.3.4", false)

	_, err := client.DNSSetCustom(context.Background(), "", "com", "ns1.example.com")
	if !errors.Is(err, errSLDRequired) {
		t.Errorf("error = %v, want %v", err, errSLDRequired)
	}

	_, err = client.DNSSetCustom(context.Background(), "example", "", "ns1.example.com")
	if !errors.Is(err, errTLDRequired) {
		t.Errorf("error = %v, want %v", err, errTLDRequired)
	}

	_, err = client.DNSSetCustom(context.Background(), "example", "com", "")
	if !errors.Is(err, errNameserversRequired) {
		t.Errorf("error = %v, want %v", err, errNameserversRequired)
	}
}

func TestParseSSLListResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.ssl.getList</RequestedCommand>
  <CommandResponse Type="namecheap.ssl.getList">
    <SSLListResult>
      <SSL CertificateID="999" HostName="example.com" SSLType="PositiveSSL" PurchaseDate="01/01/2025" ExpireDate="01/01/2026" ActivationExpireDate="02/01/2025" IsExpiredYN="false" Status="active"/>
    </SSLListResult>
  </CommandResponse>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	certs := resp.CommandResponse.SSLList.Certificates
	if len(certs) != 1 {
		t.Fatalf("got %d certs, want 1", len(certs))
	}

	if certs[0].CertificateID != "999" {
		t.Errorf("CertificateID = %q, want 999", certs[0].CertificateID)
	}

	if certs[0].SSLType != "PositiveSSL" {
		t.Errorf("SSLType = %q, want PositiveSSL", certs[0].SSLType)
	}

	if certs[0].Status != "active" {
		t.Errorf("Status = %q, want active", certs[0].Status)
	}
}

func TestParseAPIErrorResponse(t *testing.T) {
	t.Parallel()

	xmlData := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
  <Errors>
    <Error Number="1010102">Parameter APIKey is missing</Error>
  </Errors>
  <Warnings/>
  <RequestedCommand/>
  <CommandResponse/>
</ApiResponse>`

	var resp ApiResponse
	if err := xml.Unmarshal([]byte(xmlData), &resp); err != nil {
		t.Fatalf("unmarshal XML: %v", err)
	}

	if resp.Status != "ERROR" {
		t.Errorf("status = %q, want ERROR", resp.Status)
	}

	if len(resp.Errors.Errors) != 1 {
		t.Fatalf("got %d errors, want 1", len(resp.Errors.Errors))
	}

	if resp.Errors.Errors[0].Number != "1010102" {
		t.Errorf("error number = %q, want 1010102", resp.Errors.Errors[0].Number)
	}

	if resp.Errors.Errors[0].Message != "Parameter APIKey is missing" {
		t.Errorf("error message = %q, want 'Parameter APIKey is missing'", resp.Errors.Errors[0].Message)
	}
}

func TestClientValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		apiKey   string
		apiUser  string
		clientIP string
		wantErr  error
	}{
		{
			name:     "missing API key",
			apiKey:   "",
			apiUser:  "user",
			clientIP: "1.2.3.4",
			wantErr:  errMissingAPIKey,
		},
		{
			name:     "missing user",
			apiKey:   "key",
			apiUser:  "",
			clientIP: "1.2.3.4",
			wantErr:  errMissingUser,
		},
		{
			name:     "missing IP",
			apiKey:   "key",
			apiUser:  "user",
			clientIP: "",
			wantErr:  errMissingIP,
		},
		{
			name:     "all present",
			apiKey:   "key",
			apiUser:  "user",
			clientIP: "1.2.3.4",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewClient(tt.apiKey, tt.apiUser, tt.clientIP, false)
			err := client.Validate()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				return
			}

			if err == nil {
				t.Error("expected error, got nil")
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomainsCheckValidation(t *testing.T) {
	t.Parallel()

	client := NewClient("key", "user", "1.2.3.4", false)

	_, err := client.DomainsCheck(context.Background(), "")
	if !errors.Is(err, errDomainRequired) {
		t.Errorf("error = %v, want %v", err, errDomainRequired)
	}
}

func TestDNSGetHostsValidation(t *testing.T) {
	t.Parallel()

	client := NewClient("key", "user", "1.2.3.4", false)

	_, err := client.DNSGetHosts(context.Background(), "", "com")
	if !errors.Is(err, errSLDRequired) {
		t.Errorf("error = %v, want %v", err, errSLDRequired)
	}

	_, err = client.DNSGetHosts(context.Background(), "example", "")
	if !errors.Is(err, errTLDRequired) {
		t.Errorf("error = %v, want %v", err, errTLDRequired)
	}
}

func TestSandboxURL(t *testing.T) {
	t.Parallel()

	prod := NewClient("key", "user", "1.2.3.4", false)
	if prod.baseURL != productionBaseURL {
		t.Errorf("production baseURL = %q, want %q", prod.baseURL, productionBaseURL)
	}

	sandbox := NewClient("key", "user", "1.2.3.4", true)
	if sandbox.baseURL != sandboxBaseURL {
		t.Errorf("sandbox baseURL = %q, want %q", sandbox.baseURL, sandboxBaseURL)
	}
}

func TestDomainsGetListHTTP(t *testing.T) {
	t.Parallel()

	xmlResp := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors/>
  <Warnings/>
  <RequestedCommand>namecheap.domains.getList</RequestedCommand>
  <CommandResponse Type="namecheap.domains.getList">
    <DomainGetListResult>
      <Domain ID="1" Name="test.com" User="u" Created="01/01/2025" Expires="01/01/2026" IsExpired="false" IsLocked="false" AutoRenew="true" WhoisGuard="ENABLED"/>
    </DomainGetListResult>
    <Paging><TotalItems>1</TotalItems><CurrentPage>1</CurrentPage><PageSize>100</PageSize></Paging>
  </CommandResponse>
</ApiResponse>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("Command") != "namecheap.domains.getList" {
			t.Errorf("command = %q, want namecheap.domains.getList", r.URL.Query().Get("Command"))
		}

		if r.URL.Query().Get("ApiUser") != "testuser" {
			t.Errorf("ApiUser = %q, want testuser", r.URL.Query().Get("ApiUser"))
		}

		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(xmlResp))
	}))
	defer srv.Close()

	client := NewClient("testkey", "testuser", "127.0.0.1", false)
	client.baseURL = srv.URL

	resp, err := client.DomainsGetList(context.Background(), "ALL", 1, 100)
	if err != nil {
		t.Fatalf("DomainsGetList: %v", err)
	}

	domains := resp.CommandResponse.DomainList.Domains
	if len(domains) != 1 {
		t.Fatalf("got %d domains, want 1", len(domains))
	}

	if domains[0].Name != "test.com" {
		t.Errorf("domain name = %q, want test.com", domains[0].Name)
	}
}

func TestAPIErrorHandling(t *testing.T) {
	t.Parallel()

	xmlResp := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
  <Errors>
    <Error Number="2030280">Domain name not found</Error>
  </Errors>
  <Warnings/>
  <RequestedCommand>namecheap.domains.getInfo</RequestedCommand>
  <CommandResponse/>
</ApiResponse>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(xmlResp))
	}))
	defer srv.Close()

	client := NewClient("testkey", "testuser", "127.0.0.1", false)
	client.baseURL = srv.URL

	_, err := client.DomainsGetInfo(context.Background(), "notfound.com")
	if err == nil {
		t.Fatal("expected error for ERROR response, got nil")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}
}
