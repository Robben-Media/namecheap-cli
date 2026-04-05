package namecheap

import "encoding/xml"

// ApiResponse is the top-level XML response wrapper for all Namecheap API calls.
type ApiResponse struct {
	XMLName          xml.Name        `xml:"ApiResponse"`
	Status           string          `xml:"Status,attr"`
	Errors           ApiErrors       `xml:"Errors"`
	Warnings         ApiWarnings     `xml:"Warnings"`
	RequestedCommand string          `xml:"RequestedCommand"`
	CommandResponse  CommandResponse `xml:"CommandResponse"`
}

// ApiErrors holds error entries from the API response.
type ApiErrors struct {
	Errors []ApiError `xml:"Error"`
}

// ApiError represents a single API error.
type ApiError struct {
	Number  string `xml:"Number,attr"`
	Message string `xml:",chardata"`
}

// ApiWarnings holds warning entries from the API response.
type ApiWarnings struct {
	Warnings []ApiWarning `xml:"Warning"`
}

// ApiWarning represents a single API warning.
type ApiWarning struct {
	Number  string `xml:"Number,attr"`
	Message string `xml:",chardata"`
}

// CommandResponse contains the typed response data. Only one inner element
// will be populated per call; the rest stay zero-valued.
type CommandResponse struct {
	Type               string              `xml:"Type,attr"`
	DomainList         DomainListResult    `xml:"DomainGetListResult"`
	DomainChecks       []DomainCheckResult `xml:"DomainCheckResult"`
	DomainInfo         DomainInfoResult    `xml:"DomainGetInfoResult"`
	DNSHosts           DNSHostsResult      `xml:"DomainDNSGetHostsResult"`
	DNSSetResult       DNSSetResult        `xml:"DomainDNSSetHostsResult"`
	DNSSetCustomResult DNSSetCustomResult  `xml:"DomainDNSSetCustomResult"`
	SSLList            SSLListResult       `xml:"SSLListResult"`
	Paging             Paging              `xml:"Paging"`
}

// Paging contains pagination info.
type Paging struct {
	TotalItems  int `xml:"TotalItems"`
	CurrentPage int `xml:"CurrentPage"`
	PageSize    int `xml:"PageSize"`
}

// --- Domains ---

// DomainListResult holds the list of domains.
type DomainListResult struct {
	Domains []Domain `xml:"Domain"`
}

// Domain represents a single domain in the getList response.
type Domain struct {
	ID         string `xml:"ID,attr"         json:"id"`
	Name       string `xml:"Name,attr"       json:"name"`
	User       string `xml:"User,attr"       json:"user"`
	Created    string `xml:"Created,attr"    json:"created"`
	Expires    string `xml:"Expires,attr"    json:"expires"`
	IsExpired  string `xml:"IsExpired,attr"  json:"is_expired"`
	IsLocked   string `xml:"IsLocked,attr"   json:"is_locked"`
	AutoRenew  string `xml:"AutoRenew,attr"  json:"auto_renew"`
	WhoisGuard string `xml:"WhoisGuard,attr" json:"whois_guard"`
}

// DomainCheckResult represents a single domain availability check result.
// The Namecheap API returns one <DomainCheckResult> element per domain checked.
type DomainCheckResult struct {
	Domain    string `xml:"Domain,attr"    json:"domain"`
	Available string `xml:"Available,attr" json:"available"`
}

// DomainInfoResult holds detailed domain info.
type DomainInfoResult struct {
	Status        string        `xml:"Status,attr"     json:"status"`
	ID            string        `xml:"ID,attr"         json:"id"`
	DomainName    string        `xml:"DomainName,attr" json:"domain_name"`
	OwnerName     string        `xml:"OwnerName,attr"  json:"owner_name"`
	IsOwner       string        `xml:"IsOwner,attr"    json:"is_owner"`
	DomainDetails DomainDetails `xml:"DomainDetails"  json:"domain_details"`
	Whoisguard    Whoisguard    `xml:"Whoisguard"     json:"whois_guard"`
	DNSDetails    DNSDetails    `xml:"DnsDetails"     json:"dns_details"`
}

// DomainDetails holds creation/expiration dates.
type DomainDetails struct {
	CreatedDate string `xml:"CreatedDate"  json:"created_date"`
	ExpiredDate string `xml:"ExpiredDate"  json:"expired_date"`
	NumYears    int    `xml:"NumYears"     json:"num_years"`
}

// Whoisguard holds WhoisGuard status.
type Whoisguard struct {
	Enabled   string `xml:"Enabled,attr" json:"enabled"`
	ID        string `xml:"ID"           json:"id"`
	ExpiredAt string `xml:"ExpiredAt"    json:"expired_at"`
}

// DNSDetails holds DNS provider info.
type DNSDetails struct {
	ProviderType  string   `xml:"ProviderType,attr"  json:"provider_type"`
	IsUsingOurDNS string   `xml:"IsUsingOurDNS,attr" json:"is_using_our_dns"`
	Nameservers   []string `xml:"Nameserver"         json:"nameservers"`
}

// --- DNS ---

// DNSHostsResult holds the DNS host records.
type DNSHostsResult struct {
	Domain        string    `xml:"Domain,attr"    json:"domain"`
	IsUsingOurDNS string    `xml:"IsUsingOurDNS,attr" json:"is_using_our_dns"`
	Hosts         []DNSHost `xml:"host"           json:"hosts"`
}

// DNSHost represents a single DNS record.
type DNSHost struct {
	HostID        string `xml:"HostId,attr"     json:"host_id"`
	Name          string `xml:"Name,attr"       json:"name"`
	Type          string `xml:"Type,attr"       json:"type"`
	Address       string `xml:"Address,attr"    json:"address"`
	MXPref        string `xml:"MXPref,attr"     json:"mx_pref"`
	TTL           string `xml:"TTL,attr"        json:"ttl"`
	AssocAppTitle string `xml:"AssociatedAppTitle,attr" json:"assoc_app_title,omitempty"`
	FriendlyName  string `xml:"FriendlyName,attr"      json:"friendly_name,omitempty"`
	IsActive      string `xml:"IsActive,attr"   json:"is_active"`
	IsDDNSEnabled string `xml:"IsDDNSEnabled,attr" json:"is_ddns_enabled"`
}

// DNSSetResult holds the result of setting DNS hosts.
type DNSSetResult struct {
	Domain    string `xml:"Domain,attr"    json:"domain"`
	IsSuccess string `xml:"IsSuccess,attr" json:"is_success"`
}

// DNSSetCustomResult holds the result of setting custom nameservers.
type DNSSetCustomResult struct {
	Domain  string `xml:"Domain,attr"  json:"domain"`
	Updated string `xml:"Updated,attr" json:"updated"`
}

// DNSRecordInput is user-supplied JSON for setting DNS records.
type DNSRecordInput struct {
	HostName   string `json:"host_name"`
	RecordType string `json:"record_type"`
	Address    string `json:"address"`
	MXPref     string `json:"mx_pref,omitempty"`
	TTL        string `json:"ttl,omitempty"`
}

// --- SSL ---

// SSLListResult holds the list of SSL certificates.
type SSLListResult struct {
	Certificates []SSLCertificate `xml:"SSL"`
}

// SSLCertificate represents a single SSL certificate.
type SSLCertificate struct {
	CertificateID        string `xml:"CertificateID,attr" json:"certificate_id"`
	HostName             string `xml:"HostName,attr"      json:"host_name"`
	SSLType              string `xml:"SSLType,attr"        json:"ssl_type"`
	PurchaseDate         string `xml:"PurchaseDate,attr"   json:"purchase_date"`
	ExpireDate           string `xml:"ExpireDate,attr"     json:"expire_date"`
	ActivationExpireDate string `xml:"ActivationExpireDate,attr" json:"activation_expire_date"`
	IsExpired            string `xml:"IsExpiredYN,attr"    json:"is_expired"`
	Status               string `xml:"Status,attr"         json:"status"`
}
