package network

import (
	"net"
	"net/url"
	"strings"
	"testing"
)

// ── command metadata ───────────────────────────────────────────────────────

func TestNetworkCmd_Metadata(t *testing.T) {
	if networkCmd.Use == "" {
		t.Error("networkCmd.Use must not be empty")
	}
	if networkCmd.Short == "" {
		t.Error("networkCmd.Short must not be empty")
	}
}

func TestSubcommands_Metadata(t *testing.T) {
	cmds := []struct {
		name  string
		use   string
		short string
	}{
		{"dnsCmd", dnsCmd.Use, dnsCmd.Short},
		{"ipCmd", ipCmd.Use, ipCmd.Short},
		{"whoisCmd", whoisCmd.Use, whoisCmd.Short},
		{"cidrCmd", cidrCmd.Use, cidrCmd.Short},
		{"urlparseCmd", urlparseCmd.Use, urlparseCmd.Short},
		{"httpstatusCmd", httpstatusCmd.Use, httpstatusCmd.Short},
		{"useragentCmd", useragentCmd.Use, useragentCmd.Short},
		{"jwtCmd", jwtCmd.Use, jwtCmd.Short},
	}
	for _, c := range cmds {
		if c.use == "" {
			t.Errorf("%s.Use must not be empty", c.name)
		}
		if c.short == "" {
			t.Errorf("%s.Short must not be empty", c.name)
		}
	}
}

// ── flag existence ─────────────────────────────────────────────────────────

func TestNetworkCmd_PersistentFlags(t *testing.T) {
	f := networkCmd.PersistentFlags()
	if f.Lookup("json") == nil {
		t.Error("expected persistent flag --json")
	}
}

func TestDnsCmd_Flags(t *testing.T) {
	f := dnsCmd.Flags()
	if f.Lookup("type") == nil {
		t.Error("expected flag --type on dnsCmd")
	}
}

// ── CIDR calculation ───────────────────────────────────────────────────────

func TestCIDR_192_168_1_0_24(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	_, ipNet, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		t.Fatalf("ParseCIDR error: %v", err)
	}

	ones, bits := ipNet.Mask.Size()
	if ones != 24 {
		t.Errorf("prefix length = %d, want 24", ones)
	}
	if bits != 32 {
		t.Errorf("total bits = %d, want 32", bits)
	}

	// Network address
	networkAddr := ipNet.IP.String()
	if networkAddr != "192.168.1.0" {
		t.Errorf("network address = %q, want %q", networkAddr, "192.168.1.0")
	}

	// Broadcast address
	broadcast := make(net.IP, len(ipNet.IP))
	for i := range ipNet.IP {
		broadcast[i] = ipNet.IP[i] | ^ipNet.Mask[i]
	}
	if broadcast.String() != "192.168.1.255" {
		t.Errorf("broadcast address = %q, want %q", broadcast.String(), "192.168.1.255")
	}

	// Total hosts
	hostBits := bits - ones
	totalHosts := uint64(1) << uint(hostBits)
	if totalHosts != 256 {
		t.Errorf("total hosts = %d, want 256", totalHosts)
	}

	// Usable hosts
	usableHosts := totalHosts - 2
	if usableHosts != 254 {
		t.Errorf("usable hosts = %d, want 254", usableHosts)
	}

	// First usable
	first := make(net.IP, len(ipNet.IP))
	copy(first, ipNet.IP)
	first[len(first)-1]++
	if first.String() != "192.168.1.1" {
		t.Errorf("first usable = %q, want %q", first.String(), "192.168.1.1")
	}

	// Last usable
	last := make(net.IP, len(broadcast))
	copy(last, broadcast)
	last[len(last)-1]--
	if last.String() != "192.168.1.254" {
		t.Errorf("last usable = %q, want %q", last.String(), "192.168.1.254")
	}

	// Subnet mask
	subnetMask := net.IP(ipNet.Mask).String()
	if subnetMask != "255.255.255.0" {
		t.Errorf("subnet mask = %q, want %q", subnetMask, "255.255.255.0")
	}

	// Wildcard mask
	wildcard := make(net.IP, len(ipNet.Mask))
	for i := range ipNet.Mask {
		wildcard[i] = ^ipNet.Mask[i]
	}
	if wildcard.String() != "0.0.0.255" {
		t.Errorf("wildcard mask = %q, want %q", wildcard.String(), "0.0.0.255")
	}
}

func TestCIDR_10_0_0_0_8(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	_, ipNet, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		t.Fatalf("ParseCIDR error: %v", err)
	}

	ones, _ := ipNet.Mask.Size()
	if ones != 8 {
		t.Errorf("prefix length = %d, want 8", ones)
	}

	if ipNet.IP.String() != "10.0.0.0" {
		t.Errorf("network = %q, want %q", ipNet.IP.String(), "10.0.0.0")
	}
}

func TestCIDR_Slash32(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	_, ipNet, err := net.ParseCIDR("192.168.1.1/32")
	if err != nil {
		t.Fatalf("ParseCIDR error: %v", err)
	}

	ones, bits := ipNet.Mask.Size()
	if ones != 32 || bits != 32 {
		t.Errorf("expected /32, got /%d (bits=%d)", ones, bits)
	}
}

func TestCIDR_Invalid(t *testing.T) {
	_, _, err := net.ParseCIDR("not-a-cidr")
	if err == nil {
		t.Error("expected error for invalid CIDR")
	}
}

// ── URL parsing ────────────────────────────────────────────────────────────

func TestURLParse_Full(t *testing.T) {
	u, err := url.Parse("https://example.com:8080/path?q=search&lang=en#results")
	if err != nil {
		t.Fatalf("url.Parse error: %v", err)
	}

	if u.Scheme != "https" {
		t.Errorf("scheme = %q, want %q", u.Scheme, "https")
	}
	if u.Hostname() != "example.com" {
		t.Errorf("hostname = %q, want %q", u.Hostname(), "example.com")
	}
	if u.Port() != "8080" {
		t.Errorf("port = %q, want %q", u.Port(), "8080")
	}
	if u.Path != "/path" {
		t.Errorf("path = %q, want %q", u.Path, "/path")
	}
	if u.Fragment != "results" {
		t.Errorf("fragment = %q, want %q", u.Fragment, "results")
	}

	q := u.Query()
	if q.Get("q") != "search" {
		t.Errorf("query param 'q' = %q, want %q", q.Get("q"), "search")
	}
	if q.Get("lang") != "en" {
		t.Errorf("query param 'lang' = %q, want %q", q.Get("lang"), "en")
	}
}

func TestURLParse_WithUserInfo(t *testing.T) {
	u, err := url.Parse("ftp://user:pass@files.example.com/pub/file.txt")
	if err != nil {
		t.Fatalf("url.Parse error: %v", err)
	}

	if u.Scheme != "ftp" {
		t.Errorf("scheme = %q, want %q", u.Scheme, "ftp")
	}
	if u.User.Username() != "user" {
		t.Errorf("username = %q, want %q", u.User.Username(), "user")
	}
	pw, _ := u.User.Password()
	if pw != "pass" {
		t.Errorf("password = %q, want %q", pw, "pass")
	}
	if u.Hostname() != "files.example.com" {
		t.Errorf("hostname = %q, want %q", u.Hostname(), "files.example.com")
	}
}

func TestURLParse_MultipleQueryValues(t *testing.T) {
	u, err := url.Parse("https://example.com/path?a=1&a=2")
	if err != nil {
		t.Fatalf("url.Parse error: %v", err)
	}
	vals := u.Query()["a"]
	if len(vals) != 2 {
		t.Errorf("expected 2 values for 'a', got %d", len(vals))
	}
}

func TestURLParse_Minimal(t *testing.T) {
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("url.Parse error: %v", err)
	}
	if u.Scheme != "https" {
		t.Errorf("scheme = %q, want %q", u.Scheme, "https")
	}
	if u.Hostname() != "example.com" {
		t.Errorf("hostname = %q, want %q", u.Hostname(), "example.com")
	}
}

// ── HTTP status code lookup ────────────────────────────────────────────────

func TestHTTPStatus_200(t *testing.T) {
	info, ok := httpStatusCodes[200]
	if !ok {
		t.Fatal("expected 200 in httpStatusCodes")
	}
	if info.Name != "OK" {
		t.Errorf("200 name = %q, want %q", info.Name, "OK")
	}
	if info.Description == "" {
		t.Error("200 description should not be empty")
	}
}

func TestHTTPStatus_404(t *testing.T) {
	info, ok := httpStatusCodes[404]
	if !ok {
		t.Fatal("expected 404 in httpStatusCodes")
	}
	if info.Name != "Not Found" {
		t.Errorf("404 name = %q, want %q", info.Name, "Not Found")
	}
}

func TestHTTPStatus_500(t *testing.T) {
	info, ok := httpStatusCodes[500]
	if !ok {
		t.Fatal("expected 500 in httpStatusCodes")
	}
	if info.Name != "Internal Server Error" {
		t.Errorf("500 name = %q, want %q", info.Name, "Internal Server Error")
	}
}

func TestHTTPStatus_418(t *testing.T) {
	info, ok := httpStatusCodes[418]
	if !ok {
		t.Fatal("expected 418 in httpStatusCodes")
	}
	if info.Name != "I'm a Teapot" {
		t.Errorf("418 name = %q, want %q", info.Name, "I'm a Teapot")
	}
}

func TestHTTPStatus_AllCategories(t *testing.T) {
	categories := map[int]string{
		100: "Informational",
		200: "Success",
		301: "Redirection",
		404: "Client Error",
		500: "Server Error",
	}
	for code, wantCat := range categories {
		var category string
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
		if category != wantCat {
			t.Errorf("code %d: category = %q, want %q", code, category, wantCat)
		}
	}
}

func TestHTTPStatus_Unknown(t *testing.T) {
	_, ok := httpStatusCodes[999]
	if ok {
		t.Error("did not expect 999 in httpStatusCodes")
	}
}

// ── User-Agent parsing ─────────────────────────────────────────────────────

func TestParseBrowser_Chrome(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	browser, version := parseBrowser(ua)
	if browser != "Chrome" {
		t.Errorf("browser = %q, want %q", browser, "Chrome")
	}
	if version == "" {
		t.Error("expected non-empty version for Chrome")
	}
	if !strings.HasPrefix(version, "120") {
		t.Errorf("version = %q, expected to start with '120'", version)
	}
}

func TestParseBrowser_Firefox(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/109.0"
	browser, version := parseBrowser(ua)
	if browser != "Firefox" {
		t.Errorf("browser = %q, want %q", browser, "Firefox")
	}
	if version != "109.0" {
		t.Errorf("version = %q, want %q", version, "109.0")
	}
}

func TestParseBrowser_Safari(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"
	browser, _ := parseBrowser(ua)
	if browser != "Safari" {
		t.Errorf("browser = %q, want %q", browser, "Safari")
	}
}

func TestParseBrowser_Curl(t *testing.T) {
	ua := "curl/8.1.2"
	browser, version := parseBrowser(ua)
	if browser != "curl" {
		t.Errorf("browser = %q, want %q", browser, "curl")
	}
	if version != "8.1.2" {
		t.Errorf("version = %q, want %q", version, "8.1.2")
	}
}

func TestParseBrowser_Unknown(t *testing.T) {
	ua := "SomeRandomBot/1.0"
	browser, _ := parseBrowser(ua)
	if browser != "Unknown" {
		t.Errorf("browser = %q, want %q", browser, "Unknown")
	}
}

func TestParseOS_MacOS(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
	os := parseOS(ua)
	if !strings.HasPrefix(os, "macOS") {
		t.Errorf("os = %q, want prefix 'macOS'", os)
	}
}

func TestParseOS_Windows(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	os := parseOS(ua)
	if !strings.HasPrefix(os, "Windows") {
		t.Errorf("os = %q, want prefix 'Windows'", os)
	}
}

func TestParseOS_iOS(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15"
	os := parseOS(ua)
	if !strings.HasPrefix(os, "iOS") {
		t.Errorf("os = %q, want prefix 'iOS'", os)
	}
}

func TestParseOS_Android(t *testing.T) {
	ua := "Mozilla/5.0 (Linux; Android 13; Pixel 7)"
	os := parseOS(ua)
	if !strings.HasPrefix(os, "Android") {
		t.Errorf("os = %q, want prefix 'Android'", os)
	}
}

func TestParseOS_Linux(t *testing.T) {
	ua := "Mozilla/5.0 (X11; Linux x86_64)"
	os := parseOS(ua)
	if os != "Linux" {
		t.Errorf("os = %q, want %q", os, "Linux")
	}
}

func TestParseDevice_Mobile(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)"
	device := parseDevice(ua)
	if device != "Mobile" {
		t.Errorf("device = %q, want %q", device, "Mobile")
	}
}

func TestParseDevice_Tablet(t *testing.T) {
	ua := "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X)"
	device := parseDevice(ua)
	if device != "Tablet" {
		t.Errorf("device = %q, want %q", device, "Tablet")
	}
}

func TestParseDevice_Desktop(t *testing.T) {
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"
	device := parseDevice(ua)
	if device != "Desktop" {
		t.Errorf("device = %q, want %q", device, "Desktop")
	}
}

func TestParseDevice_Bot(t *testing.T) {
	ua := "Googlebot/2.1 (+http://www.google.com/bot.html)"
	device := parseDevice(ua)
	if device != "Bot" {
		t.Errorf("device = %q, want %q", device, "Bot")
	}
}

func TestParseDevice_CLI(t *testing.T) {
	ua := "curl/8.1.2"
	device := parseDevice(ua)
	if device != "CLI/Tool" {
		t.Errorf("device = %q, want %q", device, "CLI/Tool")
	}
}

func TestWindowsVersion(t *testing.T) {
	tests := []struct {
		nt, want string
	}{
		{"10.0", "10/11"},
		{"6.3", "8.1"},
		{"6.2", "8"},
		{"6.1", "7"},
		{"6.0", "Vista"},
		{"5.1", "XP"},
		{"5.0", "2000"},
		{"99.0", "NT 99.0"},
	}
	for _, tt := range tests {
		got := windowsVersion(tt.nt)
		if got != tt.want {
			t.Errorf("windowsVersion(%q) = %q, want %q", tt.nt, got, tt.want)
		}
	}
}

// ── JWT helpers ────────────────────────────────────────────────────────────

func TestBase64URLDecode_Valid(t *testing.T) {
	// "eyJ0ZXN0IjoxfQ" is base64url of {"test":1}
	decoded, err := base64URLDecode("eyJ0ZXN0IjoxfQ")
	if err != nil {
		t.Fatalf("base64URLDecode error: %v", err)
	}
	if decoded != `{"test":1}` {
		t.Errorf("decoded = %q, want %q", decoded, `{"test":1}`)
	}
}

func TestBase64URLDecode_WithPaddingNeeded(t *testing.T) {
	// Test that padding is properly added
	// "YQ" is base64url of "a" (needs == padding)
	decoded, err := base64URLDecode("YQ")
	if err != nil {
		t.Fatalf("base64URLDecode error: %v", err)
	}
	if decoded != "a" {
		t.Errorf("decoded = %q, want %q", decoded, "a")
	}
}

func TestPrettyJSON_Valid(t *testing.T) {
	input := `{"name":"test","value":42}`
	result := prettyJSON(input)
	if !strings.Contains(result, "\"name\"") {
		t.Errorf("prettyJSON should contain key 'name', got %q", result)
	}
	if !strings.Contains(result, "  ") {
		t.Error("prettyJSON should be indented")
	}
}

func TestPrettyJSON_Invalid(t *testing.T) {
	input := "not json"
	result := prettyJSON(input)
	if result != input {
		t.Errorf("prettyJSON of invalid JSON should return raw input, got %q", result)
	}
}

func TestJsonRawOrString_Valid(t *testing.T) {
	input := `{"key":"value"}`
	result := jsonRawOrString(input)
	if result == nil {
		t.Error("expected non-nil result")
	}
	// Should be parsed as a map
	if _, ok := result.(map[string]interface{}); !ok {
		t.Errorf("expected map, got %T", result)
	}
}

func TestJsonRawOrString_Invalid(t *testing.T) {
	input := "not json"
	result := jsonRawOrString(input)
	if result != input {
		t.Errorf("expected raw string for invalid JSON, got %v", result)
	}
}

func TestExtractTimestampHints_WithTimestamps(t *testing.T) {
	payload := `{"iat":1516239022,"exp":1616239022,"sub":"test"}`
	hints := extractTimestampHints(payload)
	if len(hints) == 0 {
		t.Error("expected timestamp hints from payload")
	}
	// Should contain iat and exp hints
	foundIat := false
	foundExp := false
	for _, h := range hints {
		if strings.Contains(h, "Issued At") {
			foundIat = true
		}
		if strings.Contains(h, "Expires At") {
			foundExp = true
		}
	}
	if !foundIat {
		t.Error("expected 'Issued At' hint")
	}
	if !foundExp {
		t.Error("expected 'Expires At' hint")
	}
}

func TestExtractTimestampHints_NoTimestamps(t *testing.T) {
	payload := `{"sub":"test","name":"John"}`
	hints := extractTimestampHints(payload)
	if len(hints) != 0 {
		t.Errorf("expected no hints, got %d", len(hints))
	}
}

func TestExtractTimestampHints_InvalidJSON(t *testing.T) {
	hints := extractTimestampHints("not json")
	if hints != nil {
		t.Errorf("expected nil hints for invalid JSON, got %v", hints)
	}
}

// ── whois helpers ──────────────────────────────────────────────────────────

func TestExtractWhoisReferral(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"refer: whois.verisign-grs.com\n", "whois.verisign-grs.com"},
		{"whois: whois.example.com\n", "whois.example.com"},
		{"no referral here\n", ""},
		{"refer:\n", ""},
	}
	for _, tt := range tests {
		got := extractWhoisReferral(tt.input)
		if got != tt.want {
			t.Errorf("extractWhoisReferral(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ── skip network tests in short mode ───────────────────────────────────────

func TestDNS_SkipInShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping DNS test in short mode")
	}
	// This would only run with `go test -run TestDNS_SkipInShort` (not -short)
	t.Log("DNS tests would run actual network calls here")
}

func TestWhois_SkipInShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping WHOIS test in short mode")
	}
	t.Log("WHOIS tests would run actual network calls here")
}

func TestIP_SkipInShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping IP test in short mode")
	}
	t.Log("IP tests would run actual network calls here")
}

// ── httpStatusCodes map completeness ───────────────────────────────────────

func TestHTTPStatusCodes_HasStandardCodes(t *testing.T) {
	standardCodes := []int{100, 200, 201, 204, 301, 302, 304, 400, 401, 403, 404, 405, 500, 501, 502, 503, 504}
	for _, code := range standardCodes {
		if _, ok := httpStatusCodes[code]; !ok {
			t.Errorf("httpStatusCodes missing standard code %d", code)
		}
	}
}

func TestHTTPStatusCodes_AllHaveNames(t *testing.T) {
	for code, info := range httpStatusCodes {
		if info.Name == "" {
			t.Errorf("httpStatusCodes[%d] has empty Name", code)
		}
		if info.Description == "" {
			t.Errorf("httpStatusCodes[%d] has empty Description", code)
		}
	}
}
