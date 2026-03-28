#!/usr/bin/env node

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  Tool,
} from "@modelcontextprotocol/sdk/types.js";
import { execFile } from "child_process";
import { promisify } from "util";
import { existsSync } from "fs";
import { join } from "path";

const exec = promisify(execFile);

// Find openGyver binary
function findBinary(): string {
  const candidates = [
    join(__dirname, "..", "..", "bin", "openGyver"),
    "openGyver", // PATH
  ];
  for (const c of candidates) {
    if (c === "openGyver" || existsSync(c)) return c;
  }
  return "openGyver";
}

const BINARY = findBinary();

async function runOpenGyver(args: string[]): Promise<string> {
  try {
    const { stdout, stderr } = await exec(BINARY, args, { timeout: 30000 });
    return stdout || stderr;
  } catch (err: any) {
    if (err.stdout) return err.stdout;
    throw new Error(`openGyver error: ${err.message}`);
  }
}

// Tool definitions
const tools: Tool[] = [
  // --- General ---
  {
    name: "opengyver_run",
    description:
      "Run any openGyver command. Use this for commands not covered by specific tools. Pass the full command arguments as an array.",
    inputSchema: {
      type: "object",
      properties: {
        args: {
          type: "array",
          items: { type: "string" },
          description: 'Command arguments, e.g. ["convert", "100", "cm", "in"]',
        },
      },
      required: ["args"],
    },
  },

  // --- Encoding ---
  {
    name: "opengyver_encode",
    description:
      "Encode or decode text: base64, base32, hex, url, html, binary, rot13, morse, punycode, jwt",
    inputSchema: {
      type: "object",
      properties: {
        format: {
          type: "string",
          enum: ["base64", "base32", "base58", "hex", "url", "html", "binary", "rot13", "morse", "punycode", "jwt"],
        },
        input: { type: "string", description: "Text to encode/decode" },
        decode: { type: "boolean", description: "Decode instead of encode", default: false },
      },
      required: ["format", "input"],
    },
  },

  // --- Hashing ---
  {
    name: "opengyver_hash",
    description: "Hash text with MD5, SHA-1, SHA-256, SHA-512, bcrypt, CRC32, or HMAC",
    inputSchema: {
      type: "object",
      properties: {
        algorithm: {
          type: "string",
          enum: ["md5", "sha1", "sha256", "sha384", "sha512", "bcrypt", "crc32", "adler32", "hmac"],
        },
        input: { type: "string" },
        key: { type: "string", description: "HMAC key (for hmac algorithm)" },
      },
      required: ["algorithm", "input"],
    },
  },

  // --- Unit Conversion ---
  {
    name: "opengyver_convert",
    description:
      "Convert between units: temperature (c/f/k), length (cm/in/ft/m/km/mi), weight (kg/lb/oz), volume, area, speed, data, time, and 38 currencies (live rates)",
    inputSchema: {
      type: "object",
      properties: {
        value: { type: "number" },
        from: { type: "string", description: "Source unit (e.g. cm, kg, usd)" },
        to: { type: "string", description: "Target unit (e.g. in, lb, eur)" },
      },
      required: ["value", "from", "to"],
    },
  },

  // --- JSON Tools ---
  {
    name: "opengyver_json",
    description: "JSON tools: format, minify, validate, path query, escape/unescape",
    inputSchema: {
      type: "object",
      properties: {
        action: { type: "string", enum: ["format", "minify", "validate", "path", "escape", "unescape"] },
        input: { type: "string" },
        path: { type: "string", description: "JSON path for path action (e.g. data.users[0].name)" },
      },
      required: ["action", "input"],
    },
  },

  // --- Generate ---
  {
    name: "opengyver_generate",
    description:
      "Generate passwords, passphrases, API keys, secrets, UUIDs, nanoids, OTP secrets, and more",
    inputSchema: {
      type: "object",
      properties: {
        type: {
          type: "string",
          enum: ["password", "passphrase", "apikey", "secret", "uuid", "nanoid", "otp", "snowflake", "shortid"],
        },
        length: { type: "number", description: "Length (for password, string, apikey)" },
        count: { type: "number", description: "How many to generate", default: 1 },
        prefix: { type: "string", description: "Prefix for API keys" },
        version: { type: "number", description: "UUID version (4 or 6)", default: 4 },
      },
      required: ["type"],
    },
  },

  // --- Color ---
  {
    name: "opengyver_color",
    description: "Color tools: convert between hex/rgb/hsl/cmyk, check WCAG contrast, generate palettes",
    inputSchema: {
      type: "object",
      properties: {
        action: { type: "string", enum: ["convert", "contrast", "palette", "name", "random"] },
        color: { type: "string", description: "Color value (hex, rgb, hsl)" },
        color2: { type: "string", description: "Second color (for contrast)" },
        to: { type: "string", description: "Target format for convert (hex, rgb, hsl, cmyk)" },
      },
      required: ["action"],
    },
  },

  // --- Time ---
  {
    name: "opengyver_time",
    description:
      "Time tools: current time in any timezone, convert between timezones, Unix epoch, date math, format dates",
    inputSchema: {
      type: "object",
      properties: {
        action: {
          type: "string",
          enum: ["now", "to-utc", "to-tz", "to-unix", "from-unix", "format", "diff", "add", "info", "epoch"],
        },
        input: { type: "string", description: "Time string or epoch" },
        timezone: { type: "string", description: "Target timezone (e.g. America/New_York)" },
        format: { type: "string", description: "Output format name (iso8601, rfc2822, etc.)" },
        duration: { type: "string", description: "Duration for add (e.g. 2h30m, 30d)" },
      },
      required: ["action"],
    },
  },

  // --- Network ---
  {
    name: "opengyver_network",
    description: "Network tools: DNS lookup, WHOIS, public IP, CIDR calculator, URL parser, HTTP status codes",
    inputSchema: {
      type: "object",
      properties: {
        action: { type: "string", enum: ["dns", "ip", "whois", "cidr", "urlparse", "httpstatus", "useragent"] },
        input: { type: "string" },
        type: { type: "string", description: "DNS record type (A, MX, TXT, etc.)" },
      },
      required: ["action"],
    },
  },

  // --- Data Format ---
  {
    name: "opengyver_dataformat",
    description: "Convert between data formats: YAML, TOML, XML, CSV to/from JSON",
    inputSchema: {
      type: "object",
      properties: {
        conversion: {
          type: "string",
          enum: ["yaml2json", "json2yaml", "toml2json", "json2toml", "csv2json", "json2csv", "xml2json", "json2xml"],
        },
        input: { type: "string" },
      },
      required: ["conversion", "input"],
    },
  },

  // --- Validate ---
  {
    name: "opengyver_validate",
    description: "Validate data formats: HTML, CSV, XML, YAML, TOML, JSON",
    inputSchema: {
      type: "object",
      properties: {
        format: { type: "string", enum: ["html", "csv", "xml", "yaml", "toml"] },
        input: { type: "string" },
      },
      required: ["format", "input"],
    },
  },

  // --- Format / Beautify ---
  {
    name: "opengyver_format",
    description: "Format/beautify/minify: HTML, XML, CSS, SQL",
    inputSchema: {
      type: "object",
      properties: {
        language: { type: "string", enum: ["html", "xml", "css", "sql"] },
        input: { type: "string" },
        minify: { type: "boolean", description: "Minify instead of beautify", default: false },
      },
      required: ["language", "input"],
    },
  },

  // --- Crypto ---
  {
    name: "opengyver_crypto",
    description: "Crypto tools: AES encrypt/decrypt, generate RSA/SSH keys, self-signed certificates, CSRs",
    inputSchema: {
      type: "object",
      properties: {
        action: { type: "string", enum: ["aes", "rsa", "sshkey", "cert", "csr"] },
        input: { type: "string", description: "Plaintext for AES, CN for cert/csr" },
        key: { type: "string", description: "Key for AES" },
        decrypt: { type: "boolean", description: "Decrypt mode for AES" },
        bits: { type: "number", description: "Key bits for RSA (default 2048)" },
      },
      required: ["action"],
    },
  },

  // --- Weather ---
  {
    name: "opengyver_weather",
    description: "Weather lookup: current conditions, forecast, or historical data for any city worldwide",
    inputSchema: {
      type: "object",
      properties: {
        city: { type: "string" },
        date: { type: "string", description: "Specific date (YYYY-MM-DD)" },
        from: { type: "string", description: "Start date for range" },
        to: { type: "string", description: "End date for range" },
        units: { type: "string", enum: ["celsius", "fahrenheit"], default: "celsius" },
      },
      required: ["city"],
    },
  },

  // --- Stock ---
  {
    name: "opengyver_stock",
    description: "Stock price lookup from 35+ global markets. No API key needed.",
    inputSchema: {
      type: "object",
      properties: {
        ticker: { type: "string" },
        market: { type: "string", description: "Market (kospi, kosdaq, tokyo, london, etc.)" },
        date: { type: "string", description: "Specific date (YYYY-MM-DD)" },
        from: { type: "string", description: "Start date for range" },
        to: { type: "string", description: "End date for range" },
      },
      required: ["ticker"],
    },
  },

  // --- Math ---
  {
    name: "opengyver_math",
    description: "Math tools: evaluate expressions, percentages, GCD, LCM, factorial, Fibonacci",
    inputSchema: {
      type: "object",
      properties: {
        action: { type: "string", enum: ["eval", "percent", "gcd", "lcm", "factorial", "fibonacci"] },
        expression: { type: "string", description: "Math expression for eval" },
        args: { type: "array", items: { type: "string" }, description: "Arguments for other actions" },
      },
      required: ["action"],
    },
  },
];

// Build CLI args from tool input
function buildArgs(name: string, args: Record<string, any>): string[] {
  switch (name) {
    case "opengyver_run":
      return args.args;

    case "opengyver_encode": {
      const a = ["encode", args.format, args.input];
      if (args.decode) a.splice(2, 0, "-d");
      a.push("-j");
      return a;
    }

    case "opengyver_hash": {
      const a = ["hash", args.algorithm, args.input, "-j"];
      if (args.key) a.splice(3, 0, "--key", args.key);
      return a;
    }

    case "opengyver_convert":
      return ["convert", String(args.value), args.from, args.to, "-j"];

    case "opengyver_json": {
      const a = ["json", args.action, args.input];
      if (args.path) a[2] = args.path;
      if (args.action === "path") a.push("--file", "/dev/stdin"); // handled differently
      return a;
    }

    case "opengyver_generate": {
      const a = ["generate", args.type === "uuid" ? "uuid" : args.type, "-j"];
      if (args.type === "uuid") return ["uuid", "--version", String(args.version || 4), "--count", String(args.count || 1), "-j"];
      if (args.length) a.push("--length", String(args.length));
      if (args.count) a.push("--count", String(args.count));
      if (args.prefix) a.push("--prefix", args.prefix);
      return a;
    }

    case "opengyver_color": {
      if (args.action === "convert") return ["color", "convert", args.color, "--to", args.to || "rgb", "-j"];
      if (args.action === "contrast") return ["color", "contrast", args.color, args.color2, "-j"];
      if (args.action === "palette") return ["color", "palette", args.color, "-j"];
      if (args.action === "name") return ["color", "name", args.color, "-j"];
      return ["color", "random", "-j"];
    }

    case "opengyver_time": {
      if (args.action === "epoch") return ["epoch", "-j"];
      if (args.action === "now") {
        const a = ["timex", "now", "-j"];
        if (args.timezone) a.push("--tz", args.timezone);
        return a;
      }
      if (args.action === "to-utc") return ["timex", "to-utc", args.input, "-j"];
      if (args.action === "to-tz") return ["timex", "to-tz", args.input, "--tz", args.timezone, "-j"];
      if (args.action === "to-unix") return ["timex", "to-unix", args.input];
      if (args.action === "from-unix") return ["timex", "from-unix", args.input, "-j"];
      if (args.action === "format") return ["timex", "format", args.input, "--to", args.format || "iso8601"];
      if (args.action === "diff") return ["timex", "diff", args.input, args.duration, "-j"];
      if (args.action === "add") return ["timex", "add", args.input, args.duration, "-j"];
      if (args.action === "info") return ["timex", "info", args.input, "-j"];
      return ["timex", args.action, args.input || ""];
    }

    case "opengyver_network": {
      const a = ["network", args.action];
      if (args.input) a.push(args.input);
      if (args.type) a.push("--type", args.type);
      a.push("-j");
      return a;
    }

    case "opengyver_dataformat":
      return ["dataformat", args.conversion, args.input];

    case "opengyver_validate":
      return ["validate", args.format, args.input, "-j"];

    case "opengyver_format": {
      const a = ["format", args.language, args.input];
      if (args.minify) a.push("--minify");
      return a;
    }

    case "opengyver_crypto": {
      if (args.action === "aes") {
        const a = ["crypto", "aes", args.input, "--key", args.key, "-j"];
        if (args.decrypt) a.splice(3, 0, "-d");
        return a;
      }
      if (args.action === "rsa") return ["crypto", "rsa", "--bits", String(args.bits || 2048), "-j"];
      if (args.action === "sshkey") return ["crypto", "sshkey", "-j"];
      if (args.action === "cert") return ["crypto", "cert", "--cn", args.input, "-j"];
      if (args.action === "csr") return ["crypto", "csr", "--cn", args.input, "-j"];
      return ["crypto", args.action, "-j"];
    }

    case "opengyver_weather": {
      const a = ["weather", args.city, "-j"];
      if (args.date) a.push("--date", args.date);
      if (args.from) a.push("--from", args.from);
      if (args.to) a.push("--to", args.to);
      if (args.units) a.push("--units", args.units);
      return a;
    }

    case "opengyver_stock": {
      const a = ["stock", args.ticker, "-j"];
      if (args.market) a.push("--market", args.market);
      if (args.date) a.push("--date", args.date);
      if (args.from) a.push("--from", args.from);
      if (args.to) a.push("--to", args.to);
      return a;
    }

    case "opengyver_math": {
      if (args.action === "eval") return ["math", "eval", args.expression, "-j"];
      const a = ["math", args.action, ...(args.args || []), "-j"];
      return a;
    }

    default:
      return args.args || [];
  }
}

// Server setup
const server = new Server(
  { name: "opengyver", version: "1.0.0" },
  { capabilities: { tools: {} } }
);

server.setRequestHandler(ListToolsRequestSchema, async () => ({ tools }));

server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: toolArgs } = request.params;
  const cliArgs = buildArgs(name, toolArgs as Record<string, any>);

  try {
    const result = await runOpenGyver(cliArgs);
    return { content: [{ type: "text", text: result.trim() }] };
  } catch (err: any) {
    return {
      content: [{ type: "text", text: `Error: ${err.message}` }],
      isError: true,
    };
  }
});

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch(console.error);
