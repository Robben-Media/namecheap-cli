# namecheap-cli

A CLI tool for the [Namecheap API](https://www.namecheap.com/support/api/intro/) built with Go. Manage domains, DNS records, and SSL certificates from the command line.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap builtbyrobben/tap
brew install namecheap-cli
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/builtbyrobben/namecheap-cli/releases).

### Build from Source

```bash
git clone https://github.com/builtbyrobben/namecheap-cli.git
cd namecheap-cli
make build
```

## Configuration

namecheap-cli requires three credentials: an API key, API username, and your whitelisted client IP address. Credentials can be stored in the system keyring or provided via environment variables.

### Environment Variables

| Variable | Description |
|----------|-------------|
| `NAMECHEAP_API_KEY` | Namecheap API key |
| `NAMECHEAP_USER` | Namecheap API username |
| `NAMECHEAP_CLIENT_IP` | Whitelisted client IP address |

### Store Credentials in Keyring

```bash
# Set API key (interactive prompt, recommended)
namecheap-cli auth set-key --stdin

# Set API username
namecheap-cli auth set-user myusername

# Set client IP
namecheap-cli auth set-ip 203.0.113.50

# Check credential status
namecheap-cli auth status

# Remove all stored credentials
namecheap-cli auth remove
```

## Commands

### auth -- Credential management

```bash
namecheap-cli auth set-key --stdin         # Set API key (secure prompt)
namecheap-cli auth set-user <username>     # Set API username
namecheap-cli auth set-ip <ip>             # Set whitelisted client IP
namecheap-cli auth status                  # Show authentication status
namecheap-cli auth remove                  # Remove all stored credentials
```

### domains -- Domain management

```bash
# List all domains
namecheap-cli domains list

# List expiring domains
namecheap-cli domains list --type EXPIRING

# Paginate results
namecheap-cli domains list --page 2 --page-size 50

# Check domain availability
namecheap-cli domains check "example.com,example.net"

# Get domain details
namecheap-cli domains get example.com
```

### dns -- DNS record management

```bash
# List DNS records for a domain
namecheap-cli dns list example com

# Set DNS records (replaces all records)
namecheap-cli dns set example com --records '[{"host_name":"@","record_type":"A","address":"1.2.3.4","ttl":"1800"}]'
```

### ssl -- SSL certificate management

```bash
# List all SSL certificates
namecheap-cli ssl list

# Filter by status
namecheap-cli ssl list --type Active
```

### version

```bash
namecheap-cli version
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output JSON to stdout (for scripting) |
| `--plain` | Output stable TSV text (no colors) |
| `--verbose` | Enable verbose logging |
| `--force` | Skip confirmation prompts |
| `--no-input` | Never prompt; fail instead (CI mode) |
| `--sandbox` | Use Namecheap sandbox API endpoint |
| `--color` | Color output: `auto`, `always`, or `never` |

## License

MIT
