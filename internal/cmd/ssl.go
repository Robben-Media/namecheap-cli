package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/namecheap-cli/internal/outfmt"
)

type SSLCmd struct {
	List SSLListCmd `cmd:"" help:"List SSL certificates"`
}

type SSLListCmd struct {
	Type string `help:"Filter: All, Processing, or Active" default:"" enum:",All,Processing,Active"`
}

func (cmd *SSLListCmd) Run(ctx context.Context, flags *RootFlags) error {
	client, err := getNamecheapClient(flags.Sandbox)
	if err != nil {
		return err
	}

	resp, err := client.SSLGetList(ctx, cmd.Type)
	if err != nil {
		return err
	}

	certs := resp.CommandResponse.SSLList.Certificates

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, certs)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"CERT_ID", "HOSTNAME", "TYPE", "STATUS", "EXPIRES"}

		var rows [][]string
		for _, c := range certs {
			rows = append(rows, []string{c.CertificateID, c.HostName, c.SSLType, c.Status, c.ExpireDate})
		}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	if len(certs) == 0 {
		fmt.Fprintln(os.Stderr, "No SSL certificates found")
		return nil
	}

	for _, c := range certs {
		fmt.Printf("%-12s  %-30s  %-12s  Status: %s  Expires: %s\n",
			c.CertificateID, c.HostName, c.SSLType, c.Status, c.ExpireDate)
	}

	return nil
}
