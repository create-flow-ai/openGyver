package network

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// jsonOut is the persistent --json/-j flag on the parent command.
var jsonOut bool

// ──────────────────────────────────────────────────────────────────────────────
// Parent command
// ──────────────────────────────────────────────────────────────────────────────

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network and web utilities",
	Long: `DNS lookups, public IP, WHOIS, CIDR calculation, URL parsing, HTTP status
codes, User-Agent parsing, and JWT decoding.

SUBCOMMANDS:

  dns          DNS lookup for a domain (A, AAAA, MX, TXT, NS, CNAME, SOA, SRV, PTR)
  ip           Show your public IP address
  whois        WHOIS lookup for a domain
  cidr         CIDR calculator (network, broadcast, host range, mask)
  urlparse     Parse a URL into its components
  httpstatus   Look up an HTTP status code meaning
  useragent    Parse a User-Agent string into components
  jwt          Decode and inspect a JWT token

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver network dns example.com
  openGyver network dns --type MX example.com
  openGyver network ip
  openGyver network whois example.com
  openGyver network cidr 192.168.1.0/24
  openGyver network urlparse "https://example.com:8080/path?q=1#frag"
  openGyver network httpstatus 404
  openGyver network useragent "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) ..."
  openGyver network jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`,
}

func init() {
	networkCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	networkCmd.AddCommand(dnsCmd)
	networkCmd.AddCommand(ipCmd)
	networkCmd.AddCommand(whoisCmd)
	networkCmd.AddCommand(cidrCmd)
	networkCmd.AddCommand(urlparseCmd)
	networkCmd.AddCommand(httpstatusCmd)
	networkCmd.AddCommand(useragentCmd)
	networkCmd.AddCommand(jwtCmd)

	cmd.Register(networkCmd)
}

// ──────────────────────────────────────────────────────────────────────────────
// dns
// ──────────────────────────────────────────────────────────────────────────────

var dnsRecordType string

var dnsCmd = &cobra.Command{
	Use:   "dns <domain>",
	Short: "DNS lookup for a domain",
	Long: `Perform a DNS lookup for the given domain and record type.

SUPPORTED TYPES:

  A        IPv4 address (default)
  AAAA     IPv6 address
  MX       Mail exchange servers
  TXT      TXT records (SPF, DKIM, etc.)
  NS       Name servers
  CNAME    Canonical name
  SOA      Start of Authority
  SRV      Service locator (query format: _service._proto.name)
  PTR      Reverse DNS (provide an IP address)

EXAMPLES:

  openGyver network dns example.com
  openGyver network dns --type AAAA example.com
  openGyver network dns --type MX gmail.com
  openGyver network dns --type TXT example.com
  openGyver network dns --type NS example.com
  openGyver network dns --type SOA example.com
  openGyver network dns --type SRV _sip._tcp.example.com
  openGyver network dns --type PTR 8.8.8.8
  openGyver network dns --type CNAME www.example.com
  openGyver network dns example.com --json`,
	Args: cobra.ExactArgs(1),
	RunE: runDNS,
}

func init() {
	dnsCmd.Flags().StringVar(&dnsRecordType, "type", "A", "record type: A, AAAA, MX, TXT, NS, CNAME, SOA, SRV, PTR")
}

func runDNS(_ *cobra.Command, args []string) error {
	domain := args[0]
	qtype := strings.ToUpper(dnsRecordType)

	var records []string
	var err error

	switch qtype {
	case "A":
		var ips []net.IP
		ips, err = net.LookupIP(domain)
		if err == nil {
			for _, ip := range ips {
				if ip.To4() != nil {
					records = append(records, ip.String())
				}
			}
		}
	case "AAAA":
		var ips []net.IP
		ips, err = net.LookupIP(domain)
		if err == nil {
			for _, ip := range ips {
				if ip.To4() == nil {
					records = append(records, ip.String())
				}
			}
		}
	case "MX":
		var mxs []*net.MX
		mxs, err = net.LookupMX(domain)
		if err == nil {
			for _, mx := range mxs {
				records = append(records, fmt.Sprintf("%s (priority %d)", mx.Host, mx.Pref))
			}
		}
	case "TXT":
		records, err = net.LookupTXT(domain)
	case "NS":
		var nss []*net.NS
		nss, err = net.LookupNS(domain)
		if err == nil {
			for _, ns := range nss {
				records = append(records, ns.Host)
			}
		}
	case "CNAME":
		var cname string
		cname, err = net.LookupCNAME(domain)
		if err == nil {
			records = append(records, cname)
		}
	case "SOA":
		// Go stdlib doesn't have a direct LookupSOA; show primary NS as a hint.
		var nss []*net.NS
		nss, err = net.LookupNS(domain)
		if err == nil {
			if len(nss) > 0 {
				records = append(records, fmt.Sprintf("Primary NS: %s", nss[0].Host))
			}
			records = append(records, "(full SOA record requires a raw DNS query; showing primary NS)")
		}
	case "SRV":
		var cname string
		var srvs []*net.SRV
		cname, srvs, err = net.LookupSRV("", "", domain)
		if err == nil {
			if cname != "" {
				records = append(records, fmt.Sprintf("CNAME: %s", cname))
			}
			for _, s := range srvs {
				records = append(records, fmt.Sprintf("%s:%d (priority %d, weight %d)",
					s.Target, s.Port, s.Priority, s.Weight))
			}
		}
	case "PTR":
		var names []string
		names, err = net.LookupAddr(domain)
		if err == nil {
			records = names
		}
	default:
		return fmt.Errorf("unsupported record type: %s (supported: A, AAAA, MX, TXT, NS, CNAME, SOA, SRV, PTR)", qtype)
	}

	if err != nil {
		return fmt.Errorf("DNS lookup failed: %w", err)
	}

	if records == nil {
		records = []string{}
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"domain":  domain,
			"type":    qtype,
			"records": records,
		})
	}

	fmt.Printf("DNS %s records for %s:\n", qtype, domain)
	if len(records) == 0 {
		fmt.Println("  (no records found)")
	} else {
		for _, r := range records {
			fmt.Printf("  %s\n", r)
		}
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// ip
// ──────────────────────────────────────────────────────────────────────────────

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show your public IP address",
	Long: `Query an external service to determine your public-facing IP address.

Tries https://api.ipify.org first, falls back to https://ifconfig.me/ip.

EXAMPLES:

  openGyver network ip
  openGyver network ip --json`,
	Args: cobra.NoArgs,
	RunE: runIP,
}

func runIP(_ *cobra.Command, _ []string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	ip, err := fetchIP(client, "https://api.ipify.org")
	if err != nil {
		ip, err = fetchIP(client, "https://ifconfig.me/ip")
		if err != nil {
			return fmt.Errorf("could not determine public IP: %w", err)
		}
	}

	ip = strings.TrimSpace(ip)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"ip": ip,
		})
	}

	fmt.Println(ip)
	return nil
}

func fetchIP(client *http.Client, rawURL string) (string, error) {
	resp, err := client.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ──────────────────────────────────────────────────────────────────────────────
// whois
// ──────────────────────────────────────────────────────────────────────────────

var whoisCmd = &cobra.Command{
	Use:   "whois <domain>",
	Short: "WHOIS lookup for a domain",
	Long: `Query the WHOIS database for domain registration information.

Connects to whois.iana.org to find the authoritative WHOIS server for the
TLD, then follows the referral to get full registration details.

EXAMPLES:

  openGyver network whois example.com
  openGyver network whois google.com --json`,
	Args: cobra.ExactArgs(1),
	RunE: runWhois,
}

func runWhois(_ *cobra.Command, args []string) error {
	domain := args[0]

	// Step 1: Query IANA to find the referral WHOIS server.
	ianaResp, err := queryWhoisServer("whois.iana.org", domain)
	if err != nil {
		return fmt.Errorf("IANA WHOIS query failed: %w", err)
	}

	referral := extractWhoisReferral(ianaResp)

	var result string
	if referral != "" {
		// Step 2: Query the referral server.
		result, err = queryWhoisServer(referral, domain)
		if err != nil {
			return fmt.Errorf("referral WHOIS query to %s failed: %w", referral, err)
		}
	} else {
		result = ianaResp
	}

	if jsonOut {
		data := map[string]interface{}{
			"domain": domain,
			"raw":    result,
		}
		if referral != "" {
			data["whois_server"] = referral
		}
		return cmd.PrintJSON(data)
	}

	if referral != "" {
		fmt.Printf("WHOIS server: %s\n\n", referral)
	}
	fmt.Println(result)
	return nil
}

func queryWhoisServer(server, query string) (string, error) {
	conn, err := net.DialTimeout("tcp", server+":43", 10*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	_, err = fmt.Fprintf(conn, "%s\r\n", query)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func extractWhoisReferral(response string) string {
	for _, line := range strings.Split(response, "\n") {
		lower := strings.ToLower(strings.TrimSpace(line))
		if strings.HasPrefix(lower, "refer:") || strings.HasPrefix(lower, "whois:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				server := strings.TrimSpace(parts[1])
				if server != "" {
					return server
				}
			}
		}
	}
	return ""
}

// ──────────────────────────────────────────────────────────────────────────────
// cidr
// ──────────────────────────────────────────────────────────────────────────────

var cidrCmd = &cobra.Command{
	Use:   "cidr <CIDR>",
	Short: "CIDR calculator — network details from CIDR notation",
	Long: `Calculate network information from a CIDR notation address.

Outputs the network address, broadcast address, first and last usable IPs,
total number of hosts, subnet mask, and wildcard mask.

Supports both IPv4 and IPv6 CIDR notation.

EXAMPLES:

  openGyver network cidr 192.168.1.0/24
  openGyver network cidr 10.0.0.0/8
  openGyver network cidr 172.16.0.0/12 --json
  openGyver network cidr 192.168.1.128/25`,
	Args: cobra.ExactArgs(1),
	RunE: runCIDR,
}

func runCIDR(_ *cobra.Command, args []string) error {
	_, ipNet, err := net.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("invalid CIDR notation: %w", err)
	}

	ones, bits := ipNet.Mask.Size()
	if bits == 0 {
		return fmt.Errorf("could not determine mask size")
	}

	networkAddr := ipNet.IP
	mask := ipNet.Mask

	// Broadcast address: network OR (NOT mask).
	broadcast := make(net.IP, len(networkAddr))
	for i := range networkAddr {
		broadcast[i] = networkAddr[i] | ^mask[i]
	}

	// Wildcard mask: NOT mask.
	wildcard := make(net.IP, len(mask))
	for i := range mask {
		wildcard[i] = ^mask[i]
	}

	// Total hosts.
	hostBits := bits - ones
	var totalHosts uint64
	if hostBits >= 64 {
		totalHosts = math.MaxUint64
	} else {
		totalHosts = 1 << uint(hostBits)
	}

	// Usable hosts (subtract network and broadcast for IPv4 /30 and larger).
	var usableHosts uint64
	var firstUsable, lastUsable string

	if bits == 32 { // IPv4
		if ones == 32 {
			// /32: single host.
			usableHosts = 1
			firstUsable = networkAddr.String()
			lastUsable = networkAddr.String()
		} else if ones == 31 {
			// /31: point-to-point link (RFC 3021), 2 usable.
			usableHosts = 2
			firstUsable = networkAddr.String()
			lastUsable = broadcast.String()
		} else {
			usableHosts = totalHosts - 2
			first := make(net.IP, len(networkAddr))
			copy(first, networkAddr)
			first[len(first)-1]++
			firstUsable = first.String()

			last := make(net.IP, len(broadcast))
			copy(last, broadcast)
			last[len(last)-1]--
			lastUsable = last.String()
		}
	} else { // IPv6
		usableHosts = totalHosts
		firstUsable = networkAddr.String()
		lastUsable = broadcast.String()
	}

	subnetMask := net.IP(mask).String()
	wildcardMask := wildcard.String()

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"cidr":          args[0],
			"network":       networkAddr.String(),
			"broadcast":     broadcast.String(),
			"first_usable":  firstUsable,
			"last_usable":   lastUsable,
			"total_hosts":   totalHosts,
			"usable_hosts":  usableHosts,
			"subnet_mask":   subnetMask,
			"wildcard_mask": wildcardMask,
			"prefix_length": ones,
		})
	}

	fmt.Printf("CIDR:           %s\n", args[0])
	fmt.Printf("Network:        %s\n", networkAddr)
	fmt.Printf("Broadcast:      %s\n", broadcast)
	fmt.Printf("First usable:   %s\n", firstUsable)
	fmt.Printf("Last usable:    %s\n", lastUsable)
	fmt.Printf("Total hosts:    %d\n", totalHosts)
	fmt.Printf("Usable hosts:   %d\n", usableHosts)
	fmt.Printf("Subnet mask:    %s\n", subnetMask)
	fmt.Printf("Wildcard mask:  %s\n", wildcardMask)
	fmt.Printf("Prefix length:  /%d\n", ones)
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// urlparse
// ──────────────────────────────────────────────────────────────────────────────

var urlparseCmd = &cobra.Command{
	Use:   "urlparse <url>",
	Short: "Parse a URL into its components",
	Long: `Break a URL into its constituent parts: scheme, user, host, port, path,
query parameters, and fragment.

EXAMPLES:

  openGyver network urlparse "https://example.com:8080/path?q=search&lang=en#results"
  openGyver network urlparse "ftp://user:pass@files.example.com/pub/file.txt"
  openGyver network urlparse "https://example.com/path?a=1&a=2" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runURLParse,
}

func runURLParse(_ *cobra.Command, args []string) error {
	u, err := url.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	host := u.Hostname()
	port := u.Port()

	// Build query param map (preserving multiple values).
	queryParams := make(map[string]interface{})
	for k, v := range u.Query() {
		if len(v) == 1 {
			queryParams[k] = v[0]
		} else {
			queryParams[k] = v
		}
	}

	user := ""
	password := ""
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}

	if jsonOut {
		data := map[string]interface{}{
			"scheme":   u.Scheme,
			"host":     host,
			"port":     port,
			"path":     u.Path,
			"query":    queryParams,
			"fragment": u.Fragment,
			"raw":      args[0],
		}
		if user != "" {
			data["user"] = user
		}
		if password != "" {
			data["password"] = password
		}
		return cmd.PrintJSON(data)
	}

	fmt.Printf("Scheme:     %s\n", u.Scheme)
	if user != "" {
		fmt.Printf("User:       %s\n", user)
	}
	if password != "" {
		fmt.Printf("Password:   %s\n", password)
	}
	fmt.Printf("Host:       %s\n", host)
	fmt.Printf("Port:       %s\n", port)
	fmt.Printf("Path:       %s\n", u.Path)
	if len(queryParams) > 0 {
		fmt.Println("Query params:")
		for k, v := range u.Query() {
			fmt.Printf("  %s = %s\n", k, strings.Join(v, ", "))
		}
	}
	fmt.Printf("Fragment:   %s\n", u.Fragment)
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// httpstatus
// ──────────────────────────────────────────────────────────────────────────────

var httpstatusCmd = &cobra.Command{
	Use:   "httpstatus <code>",
	Short: "Look up an HTTP status code meaning",
	Long: `Display the name and description for a given HTTP status code.

Covers all standard codes from RFC 9110 and common extensions.

EXAMPLES:

  openGyver network httpstatus 200
  openGyver network httpstatus 404
  openGyver network httpstatus 503 --json
  openGyver network httpstatus 418`,
	Args: cobra.ExactArgs(1),
	RunE: runHTTPStatus,
}

type httpStatusInfo struct {
	Name        string
	Description string
}

var httpStatusCodes = map[int]httpStatusInfo{
	// 1xx Informational
	100: {"Continue", "The server has received the request headers and the client should proceed to send the request body."},
	101: {"Switching Protocols", "The server is switching protocols as requested by the client via the Upgrade header."},
	102: {"Processing", "The server has received and is processing the request, but no response is available yet (WebDAV)."},
	103: {"Early Hints", "Used to return some response headers before the final HTTP message."},

	// 2xx Success
	200: {"OK", "The request has succeeded. The response meaning depends on the request method."},
	201: {"Created", "The request has been fulfilled and a new resource has been created."},
	202: {"Accepted", "The request has been accepted for processing, but processing has not been completed."},
	203: {"Non-Authoritative Information", "The returned metadata is not exactly the same as available from the origin server."},
	204: {"No Content", "The server has fulfilled the request but there is no content to send in the response body."},
	205: {"Reset Content", "The server has fulfilled the request and the client should reset the document view."},
	206: {"Partial Content", "The server is delivering only part of the resource due to a Range header sent by the client."},
	207: {"Multi-Status", "Provides status for multiple independent operations (WebDAV)."},
	208: {"Already Reported", "The members of a DAV binding have already been enumerated in a previous reply (WebDAV)."},
	226: {"IM Used", "The server has fulfilled a GET request and the response is a representation of one or more instance-manipulations applied to the current instance."},

	// 3xx Redirection
	300: {"Multiple Choices", "There are multiple options for the resource, and the client may choose one."},
	301: {"Moved Permanently", "The resource has been permanently moved to a new URL. Future requests should use the new URL."},
	302: {"Found", "The resource resides temporarily under a different URL. The client should continue to use the original URL."},
	303: {"See Other", "The response to the request can be found under a different URL using a GET method."},
	304: {"Not Modified", "The resource has not been modified since the version specified by the request headers."},
	305: {"Use Proxy", "The requested resource must be accessed through the proxy given in the Location header (deprecated)."},
	307: {"Temporary Redirect", "The resource resides temporarily under a different URL. The method and body should not change."},
	308: {"Permanent Redirect", "The resource has been permanently moved. The method and body should not change."},

	// 4xx Client Errors
	400: {"Bad Request", "The server cannot process the request due to malformed syntax or invalid request framing."},
	401: {"Unauthorized", "The request requires user authentication. The client must provide valid credentials."},
	402: {"Payment Required", "Reserved for future use. Some APIs use this for rate-limited or paid-tier endpoints."},
	403: {"Forbidden", "The server understood the request but refuses to authorize it."},
	404: {"Not Found", "The server cannot find the requested resource. The URL may be incorrect or the resource may have been removed."},
	405: {"Method Not Allowed", "The request method is not supported for the requested resource."},
	406: {"Not Acceptable", "The resource cannot generate a response matching the Accept headers sent by the client."},
	407: {"Proxy Authentication Required", "The client must first authenticate itself with the proxy."},
	408: {"Request Timeout", "The server timed out waiting for the request from the client."},
	409: {"Conflict", "The request conflicts with the current state of the resource."},
	410: {"Gone", "The resource is no longer available and no forwarding address is known. This is permanent."},
	411: {"Length Required", "The server requires a Content-Length header to be specified."},
	412: {"Precondition Failed", "One or more conditions in the request header fields evaluated to false on the server."},
	413: {"Content Too Large", "The request body is larger than the server is willing or able to process."},
	414: {"URI Too Long", "The URI provided was too long for the server to process."},
	415: {"Unsupported Media Type", "The media type of the request body is not supported by the server."},
	416: {"Range Not Satisfiable", "The range specified in the Range header cannot be fulfilled."},
	417: {"Expectation Failed", "The Expect header value could not be met by the server."},
	418: {"I'm a Teapot", "The server refuses to brew coffee because it is, permanently, a teapot (RFC 2324)."},
	421: {"Misdirected Request", "The request was directed at a server that cannot produce a response."},
	422: {"Unprocessable Content", "The request was well-formed but contained semantic errors (WebDAV)."},
	423: {"Locked", "The resource that is being accessed is locked (WebDAV)."},
	424: {"Failed Dependency", "The request failed because it depended on another request which failed (WebDAV)."},
	425: {"Too Early", "The server is unwilling to risk processing a request that might be replayed."},
	426: {"Upgrade Required", "The client should switch to a different protocol as indicated in the Upgrade header."},
	428: {"Precondition Required", "The origin server requires the request to be conditional."},
	429: {"Too Many Requests", "The user has sent too many requests in a given amount of time (rate limiting)."},
	431: {"Request Header Fields Too Large", "The server is unwilling to process the request because its header fields are too large."},
	451: {"Unavailable For Legal Reasons", "The resource is unavailable due to legal demands (e.g., censorship or government-mandated blocking)."},

	// 5xx Server Errors
	500: {"Internal Server Error", "The server encountered an unexpected condition that prevented it from fulfilling the request."},
	501: {"Not Implemented", "The server does not support the functionality required to fulfill the request."},
	502: {"Bad Gateway", "The server, while acting as a gateway or proxy, received an invalid response from the upstream server."},
	503: {"Service Unavailable", "The server is currently unable to handle the request due to overload or maintenance."},
	504: {"Gateway Timeout", "The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server."},
	505: {"HTTP Version Not Supported", "The HTTP version used in the request is not supported by the server."},
	506: {"Variant Also Negotiates", "The server has an internal configuration error: transparent content negotiation results in a circular reference."},
	507: {"Insufficient Storage", "The server is unable to store the representation needed to complete the request (WebDAV)."},
	508: {"Loop Detected", "The server detected an infinite loop while processing the request (WebDAV)."},
	510: {"Not Extended", "Further extensions to the request are required for the server to fulfill it."},
	511: {"Network Authentication Required", "The client needs to authenticate to gain network access (e.g., captive portal)."},
}

func runHTTPStatus(_ *cobra.Command, args []string) error {
	code, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid status code: %q (must be a number)", args[0])
	}

	info, ok := httpStatusCodes[code]
	if !ok {
		// Try Go's built-in StatusText as a fallback.
		text := http.StatusText(code)
		if text != "" {
			info = httpStatusInfo{Name: text, Description: "(no extended description available)"}
		} else {
			return fmt.Errorf("unknown HTTP status code: %d", code)
		}
	}

	category := ""
	switch {
	case code >= 100 && code < 200:
		category = "Informational"
	case code >= 200 && code < 300:
		category = "Success"
	case code >= 300 && code < 400:
		category = "Redirection"
	case code >= 400 && code < 500:
		category = "Client Error"
	case code >= 500 && code < 600:
		category = "Server Error"
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"code":        code,
			"name":        info.Name,
			"description": info.Description,
			"category":    category,
		})
	}

	fmt.Printf("%d %s\n", code, info.Name)
	fmt.Printf("Category:    %s\n", category)
	fmt.Printf("Description: %s\n", info.Description)
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// useragent
// ──────────────────────────────────────────────────────────────────────────────

var useragentCmd = &cobra.Command{
	Use:   "useragent <user-agent-string>",
	Short: "Parse a User-Agent string into components",
	Long: `Extract browser, version, operating system, and device type from a
User-Agent string using pattern matching.

Recognizes common browsers (Chrome, Firefox, Safari, Edge, Opera, etc.),
operating systems (Windows, macOS, Linux, Android, iOS), and device types.

EXAMPLES:

  openGyver network useragent "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
  openGyver network useragent "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
  openGyver network useragent "curl/8.1.2" --json`,
	Args: cobra.ExactArgs(1),
	RunE: runUserAgent,
}

// Precompiled regexes for User-Agent parsing.
var (
	reBrowserEdge      = regexp.MustCompile(`Edg(?:e|A|iOS)?/(\S+)`)
	reBrowserOpera     = regexp.MustCompile(`(?:OPR|Opera)[/ ](\S+)`)
	reBrowserSamsung   = regexp.MustCompile(`SamsungBrowser/(\S+)`)
	reBrowserUC        = regexp.MustCompile(`UCBrowser/(\S+)`)
	reBrowserBrave     = regexp.MustCompile(`Brave/(\S+)`)
	reBrowserVivaldi   = regexp.MustCompile(`Vivaldi/(\S+)`)
	reBrowserFirefox   = regexp.MustCompile(`Firefox/(\S+)`)
	reBrowserChrome    = regexp.MustCompile(`Chrome/(\S+)`)
	reBrowserSafari    = regexp.MustCompile(`Version/(\S+).*Safari`)
	reBrowserCurl      = regexp.MustCompile(`curl/(\S+)`)
	reBrowserWget      = regexp.MustCompile(`(?i)Wget/(\S+)`)
	reBrowserGooglebot = regexp.MustCompile(`Googlebot/(\S+)`)
	reBrowserBingbot   = regexp.MustCompile(`bingbot/(\S+)`)

	reOSiDevice   = regexp.MustCompile(`iPhone|iPad|iPod`)
	reOSiOS       = regexp.MustCompile(`CPU (?:iPhone )?OS (\d+[_.\d]*)`)
	reOSAndroid   = regexp.MustCompile(`Android`)
	reOSAndroidV  = regexp.MustCompile(`Android (\d+[.\d]*)`)
	reOSMacOSX    = regexp.MustCompile(`Mac OS X`)
	reOSMacOSXV   = regexp.MustCompile(`Mac OS X (\d+[_.\d]*)`)
	reOSWindowsNT = regexp.MustCompile(`Windows NT`)
	reOSWindowsV  = regexp.MustCompile(`Windows NT (\d+\.\d+)`)
	reOSChromeOS  = regexp.MustCompile(`CrOS`)
	reOSLinux     = regexp.MustCompile(`Linux`)
	reOSFreeBSD   = regexp.MustCompile(`FreeBSD`)

	reDeviceMobile  = regexp.MustCompile(`(?i)Mobi|iPhone|iPod|Android.*Mobile|Opera Mini|IEMobile`)
	reDeviceTablet  = regexp.MustCompile(`(?i)iPad|Tablet|Kindle|Silk`)
	reDeviceAndroid = regexp.MustCompile(`(?i)Android`)
	reDeviceMobileW = regexp.MustCompile(`(?i)Mobile`)
	reDeviceBot     = regexp.MustCompile(`(?i)bot|crawl|spider|slurp|Googlebot|Bingbot`)
	reDeviceCLI     = regexp.MustCompile(`(?i)curl|wget|Postman|httpie`)
)

func runUserAgent(_ *cobra.Command, args []string) error {
	ua := args[0]

	browser, version := parseBrowser(ua)
	osName := parseOS(ua)
	device := parseDevice(ua)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"raw":             ua,
			"browser":         browser,
			"browser_version": version,
			"os":              osName,
			"device":          device,
		})
	}

	fmt.Printf("Browser:   %s\n", browser)
	fmt.Printf("Version:   %s\n", version)
	fmt.Printf("OS:        %s\n", osName)
	fmt.Printf("Device:    %s\n", device)
	return nil
}

func parseBrowser(ua string) (string, string) {
	// Order matters: more specific patterns first.
	patterns := []struct {
		name string
		re   *regexp.Regexp
	}{
		{"Edge", reBrowserEdge},
		{"Opera", reBrowserOpera},
		{"Samsung Internet", reBrowserSamsung},
		{"UC Browser", reBrowserUC},
		{"Brave", reBrowserBrave},
		{"Vivaldi", reBrowserVivaldi},
		{"Firefox", reBrowserFirefox},
		{"Chrome", reBrowserChrome},
		{"Safari", reBrowserSafari},
		{"curl", reBrowserCurl},
		{"wget", reBrowserWget},
		{"Googlebot", reBrowserGooglebot},
		{"Bingbot", reBrowserBingbot},
	}

	for _, p := range patterns {
		m := p.re.FindStringSubmatch(ua)
		if m != nil {
			return p.name, m[1]
		}
	}

	return "Unknown", ""
}

func parseOS(ua string) string {
	switch {
	case reOSiDevice.MatchString(ua):
		m := reOSiOS.FindStringSubmatch(ua)
		if m != nil {
			return "iOS " + strings.ReplaceAll(m[1], "_", ".")
		}
		return "iOS"
	case reOSAndroid.MatchString(ua):
		m := reOSAndroidV.FindStringSubmatch(ua)
		if m != nil {
			return "Android " + m[1]
		}
		return "Android"
	case reOSMacOSX.MatchString(ua):
		m := reOSMacOSXV.FindStringSubmatch(ua)
		if m != nil {
			return "macOS " + strings.ReplaceAll(m[1], "_", ".")
		}
		return "macOS"
	case reOSWindowsNT.MatchString(ua):
		m := reOSWindowsV.FindStringSubmatch(ua)
		if m != nil {
			return "Windows " + windowsVersion(m[1])
		}
		return "Windows"
	case reOSChromeOS.MatchString(ua):
		return "Chrome OS"
	case reOSLinux.MatchString(ua):
		return "Linux"
	case reOSFreeBSD.MatchString(ua):
		return "FreeBSD"
	default:
		return "Unknown"
	}
}

func windowsVersion(nt string) string {
	versions := map[string]string{
		"10.0": "10/11",
		"6.3":  "8.1",
		"6.2":  "8",
		"6.1":  "7",
		"6.0":  "Vista",
		"5.1":  "XP",
		"5.0":  "2000",
	}
	if v, ok := versions[nt]; ok {
		return v
	}
	return "NT " + nt
}

func parseDevice(ua string) string {
	switch {
	case reDeviceMobile.MatchString(ua):
		return "Mobile"
	case reDeviceTablet.MatchString(ua):
		return "Tablet"
	case reDeviceAndroid.MatchString(ua) && !reDeviceMobileW.MatchString(ua):
		// Android without "Mobile" is typically a tablet.
		return "Tablet"
	case reDeviceBot.MatchString(ua):
		return "Bot"
	case reDeviceCLI.MatchString(ua):
		return "CLI/Tool"
	default:
		return "Desktop"
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// jwt
// ──────────────────────────────────────────────────────────────────────────────

var jwtCmd = &cobra.Command{
	Use:   "jwt <token>",
	Short: "Decode and inspect a JWT token",
	Long: `Decode a JSON Web Token (JWT) by splitting it into its three parts
(header, payload, signature), base64-decoding the header and payload,
and displaying them as formatted JSON.

NOTE: This does NOT verify the signature. It is an inspection/debugging tool.

EXAMPLES:

  openGyver network jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
  openGyver network jwt <token> --json`,
	Args: cobra.ExactArgs(1),
	RunE: runJWT,
}

func runJWT(_ *cobra.Command, args []string) error {
	token := args[0]
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid JWT: expected 3 dot-separated parts, got %d", len(parts))
	}

	header, err := base64URLDecode(parts[0])
	if err != nil {
		return fmt.Errorf("failed to decode JWT header: %w", err)
	}

	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	signature := parts[2]

	// Check for common claims in payload to display timestamps nicely.
	payloadHints := extractTimestampHints(payload)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"header":    jsonRawOrString(header),
			"payload":   jsonRawOrString(payload),
			"signature": signature,
		})
	}

	fmt.Println("=== Header ===")
	fmt.Println(prettyJSON(header))
	fmt.Println()
	fmt.Println("=== Payload ===")
	fmt.Println(prettyJSON(payload))
	if len(payloadHints) > 0 {
		fmt.Println()
		fmt.Println("=== Timestamps ===")
		for _, h := range payloadHints {
			fmt.Println(h)
		}
	}
	fmt.Println()
	fmt.Printf("=== Signature ===\n%s\n", signature)
	return nil
}

func base64URLDecode(s string) (string, error) {
	// JWT uses base64url encoding without padding.
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// prettyJSON attempts to format a JSON string with indentation.
// If parsing fails, returns the raw string as-is.
func prettyJSON(s string) string {
	var m interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return s
	}
	return string(b)
}

// jsonRawOrString returns the decoded JSON as a raw object if valid,
// otherwise as a plain string.
func jsonRawOrString(s string) interface{} {
	var m interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}
	return m
}

// extractTimestampHints looks for common JWT timestamp claims (iat, exp, nbf)
// and returns human-readable descriptions.
func extractTimestampHints(payload string) []string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		return nil
	}

	names := map[string]string{
		"iat": "Issued At",
		"exp": "Expires At",
		"nbf": "Not Before",
	}

	var hints []string
	for claim, label := range names {
		if v, ok := m[claim]; ok {
			if ts, ok := v.(float64); ok {
				t := time.Unix(int64(ts), 0).UTC()
				hints = append(hints, fmt.Sprintf("  %s (%s): %s", label, claim, t.Format(time.RFC3339)))
			}
		}
	}
	return hints
}
