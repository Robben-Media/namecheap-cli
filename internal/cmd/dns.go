package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/builtbyrobben/namecheap-cli/internal/namecheap"
	"github.com/builtbyrobben/namecheap-cli/internal/outfmt"
)

type DNSCmd struct {
	List DNSListCmd `cmd:"" help:"List DNS records for a domain"`
	Set  DNSSetCmd  `cmd:"" help:"Set DNS records for a domain"`
}

type DNSListCmd struct {
	SLD string `arg:"" required:"" help:"Second-level domain (e.g., example)"`
	TLD string `arg:"" required:"" help:"Top-level domain (e.g., com)"`
}

func (cmd *DNSListCmd) Run(ctx context.Context, flags *RootFlags) error {
	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.DNSGetHosts(ctx, cmd.SLD, cmd.TLD)
	if err != nil {
		return err
	}

	hosts := resp.CommandResponse.DNSHosts

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, hosts)
	}

	if len(hosts.Hosts) == 0 {
		fmt.Fprintf(os.Stderr, "No DNS records found for %s.%s\n", cmd.SLD, cmd.TLD)
		return nil
	}

	fmt.Fprintf(os.Stderr, "DNS records for %s.%s\n\n", cmd.SLD, cmd.TLD)

	for _, h := range hosts.Hosts {
		fmt.Printf("%-20s  %-8s  %-40s  TTL: %s\n", h.Name, h.Type, h.Address, h.TTL)
	}

	return nil
}

type DNSSetCmd struct {
	SLD     string `arg:"" required:"" help:"Second-level domain (e.g., example)"`
	TLD     string `arg:"" required:"" help:"Top-level domain (e.g., com)"`
	Records string `required:"" help:"JSON array of records: [{\"host_name\":\"@\",\"record_type\":\"A\",\"address\":\"1.2.3.4\"}]"`
}

func (cmd *DNSSetCmd) Run(ctx context.Context, flags *RootFlags) error {
	var records []namecheap.DNSRecordInput

	if err := json.Unmarshal([]byte(cmd.Records), &records); err != nil {
		return fmt.Errorf("parse records JSON: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("at least one DNS record is required")
	}

	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.DNSSetHosts(ctx, cmd.SLD, cmd.TLD, records)
	if err != nil {
		return err
	}

	result := resp.CommandResponse.DNSSetResult

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if result.IsSuccess == "true" {
		fmt.Fprintf(os.Stderr, "DNS records updated for %s.%s\n", cmd.SLD, cmd.TLD)
	} else {
		fmt.Fprintf(os.Stderr, "Failed to update DNS records for %s.%s\n", cmd.SLD, cmd.TLD)
	}

	return nil
}
