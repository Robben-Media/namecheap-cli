package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/namecheap-cli/internal/outfmt"
)

type DomainsCmd struct {
	List  DomainsListCmd  `cmd:"" help:"List domains in your account"`
	Check DomainsCheckCmd `cmd:"" help:"Check domain availability"`
	Get   DomainsGetCmd   `cmd:"" help:"Get domain details"`
}

type DomainsListCmd struct {
	Type     string `help:"Filter: ALL, EXPIRING, or EXPIRED" default:"ALL" enum:"ALL,EXPIRING,EXPIRED"`
	Page     int    `help:"Page number" default:"1"`
	PageSize int    `help:"Results per page" default:"100" name:"page-size"`
}

func (cmd *DomainsListCmd) Run(ctx context.Context, flags *RootFlags) error {
	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.DomainsGetList(ctx, cmd.Type, cmd.Page, cmd.PageSize)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, map[string]any{
			"domains": resp.CommandResponse.DomainList.Domains,
			"paging": map[string]int{
				"total_items":  resp.CommandResponse.Paging.TotalItems,
				"current_page": resp.CommandResponse.Paging.CurrentPage,
				"page_size":    resp.CommandResponse.Paging.PageSize,
			},
		})
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"NAME", "EXPIRES", "AUTO_RENEW", "WHOIS_GUARD"}

		var rows [][]string
		for _, d := range resp.CommandResponse.DomainList.Domains {
			rows = append(rows, []string{d.Name, d.Expires, d.AutoRenew, d.WhoisGuard})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	domains := resp.CommandResponse.DomainList.Domains
	if len(domains) == 0 {
		fmt.Fprintln(os.Stderr, "No domains found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d domains (page %d)\n\n",
		len(domains), resp.CommandResponse.Paging.TotalItems, resp.CommandResponse.Paging.CurrentPage)

	for _, d := range domains {
		fmt.Printf("%-30s  Expires: %s  AutoRenew: %s  WhoisGuard: %s\n",
			d.Name, d.Expires, d.AutoRenew, d.WhoisGuard)
	}

	return nil
}

type DomainsCheckCmd struct {
	Domains string `arg:"" required:"" help:"Comma-separated domains to check (e.g., example.com,test.net)"`
}

func (cmd *DomainsCheckCmd) Run(ctx context.Context, flags *RootFlags) error {
	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.DomainsCheck(ctx, cmd.Domains)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, resp.CommandResponse.DomainChecks)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"DOMAIN", "AVAILABLE"}

		var rows [][]string
		for _, d := range resp.CommandResponse.DomainChecks {
			rows = append(rows, []string{d.Domain, d.Available})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	for _, d := range resp.CommandResponse.DomainChecks {
		status := "unavailable"
		if d.Available == "true" {
			status = "available"
		}

		fmt.Printf("%-30s  %s\n", d.Domain, status)
	}

	return nil
}

type DomainsGetCmd struct {
	Domain string `arg:"" required:"" help:"Domain name to get details for"`
}

func (cmd *DomainsGetCmd) Run(ctx context.Context, flags *RootFlags) error {
	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.DomainsGetInfo(ctx, cmd.Domain)
	if err != nil {
		return err
	}

	info := resp.CommandResponse.DomainInfo

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, info)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"DOMAIN", "STATUS", "OWNER", "CREATED", "EXPIRES", "DNS_TYPE"}
		rows := [][]string{{info.DomainName, info.Status, info.OwnerName, info.DomainDetails.CreatedDate, info.DomainDetails.ExpiredDate, info.DNSDetails.ProviderType}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Printf("Domain:     %s\n", info.DomainName)
	fmt.Printf("Status:     %s\n", info.Status)
	fmt.Printf("Owner:      %s\n", info.OwnerName)
	fmt.Printf("Created:    %s\n", info.DomainDetails.CreatedDate)
	fmt.Printf("Expires:    %s\n", info.DomainDetails.ExpiredDate)
	fmt.Printf("WhoisGuard: %s\n", info.Whoisguard.Enabled)
	fmt.Printf("DNS Type:   %s\n", info.DNSDetails.ProviderType)

	if len(info.DNSDetails.Nameservers) > 0 {
		fmt.Println("Nameservers:")

		for _, ns := range info.DNSDetails.Nameservers {
			fmt.Printf("  %s\n", ns)
		}
	}

	return nil
}
