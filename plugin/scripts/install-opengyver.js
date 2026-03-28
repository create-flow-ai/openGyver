#!/usr/bin/env node

/**
 * Auto-install openGyver binary from GitHub releases.
 * Detects OS/arch and downloads the correct binary.
 */

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");

const REPO = "create-flow-ai/openGyver";
const BIN_DIR = path.join(__dirname, "..", "bin");

function getPlatform() {
  const platform = os.platform(); // darwin, linux, win32
  const arch = os.arch(); // x64, arm64

  const osMap = { darwin: "darwin", linux: "linux", win32: "windows" };
  const archMap = { x64: "amd64", arm64: "arm64" };

  const goos = osMap[platform];
  const goarch = archMap[arch];

  if (!goos || !goarch) {
    throw new Error(`Unsupported platform: ${platform}/${arch}`);
  }

  const ext = platform === "win32" ? ".zip" : ".tar.gz";
  return { goos, goarch, ext };
}

function fetchJSON(url) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "User-Agent": "opengyver-installer" } }, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        return fetchJSON(res.headers.location).then(resolve).catch(reject);
      }
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try {
          resolve(JSON.parse(data));
        } catch (e) {
          reject(new Error(`Failed to parse response: ${data.slice(0, 200)}`));
        }
      });
      res.on("error", reject);
    });
  });
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, { headers: { "User-Agent": "opengyver-installer" } }, (res) => {
      if (res.statusCode === 302 || res.statusCode === 301) {
        return downloadFile(res.headers.location, dest).then(resolve).catch(reject);
      }
      const file = fs.createWriteStream(dest);
      res.pipe(file);
      file.on("finish", () => {
        file.close();
        resolve();
      });
      res.on("error", reject);
    });
  });
}

async function install() {
  // Check if already installed
  const binaryName = os.platform() === "win32" ? "openGyver.exe" : "openGyver";
  const binaryPath = path.join(BIN_DIR, binaryName);

  if (fs.existsSync(binaryPath)) {
    console.log(`openGyver already installed at ${binaryPath}`);
    return;
  }

  // Check if available via homebrew or PATH
  try {
    execSync("openGyver --version", { stdio: "pipe" });
    console.log("openGyver found in PATH, skipping download");
    return;
  } catch {}

  console.log("Installing openGyver...");

  const { goos, goarch, ext } = getPlatform();
  const assetName = `openGyver_${goos}_${goarch}${ext}`;

  // Get latest release
  let release;
  try {
    release = await fetchJSON(`https://api.github.com/repos/${REPO}/releases/latest`);
  } catch (e) {
    console.warn(`Could not fetch latest release: ${e.message}`);
    console.warn("Please install openGyver manually:");
    console.warn("  brew tap create-flow-ai/tap && brew install opengyver");
    console.warn("  OR: go install github.com/mj/opengyver@latest");
    return;
  }

  const asset = release.assets?.find((a) => a.name === assetName);
  if (!asset) {
    console.warn(`No binary found for ${goos}/${goarch} in release ${release.tag_name}`);
    console.warn(`Available assets: ${release.assets?.map((a) => a.name).join(", ")}`);
    console.warn("Install manually: brew tap create-flow-ai/tap && brew install opengyver");
    return;
  }

  // Download
  fs.mkdirSync(BIN_DIR, { recursive: true });
  const tmpFile = path.join(os.tmpdir(), assetName);

  console.log(`Downloading ${asset.name} (${(asset.size / 1024 / 1024).toFixed(1)} MB)...`);
  await downloadFile(asset.browser_download_url, tmpFile);

  // Extract
  if (ext === ".tar.gz") {
    execSync(`tar xzf "${tmpFile}" -C "${BIN_DIR}" openGyver`, { stdio: "pipe" });
  } else {
    execSync(`unzip -o "${tmpFile}" openGyver.exe -d "${BIN_DIR}"`, { stdio: "pipe" });
  }

  // Make executable
  fs.chmodSync(binaryPath, 0o755);

  // Cleanup
  try {
    fs.unlinkSync(tmpFile);
  } catch {}

  console.log(`openGyver installed to ${binaryPath}`);
}

install().catch((err) => {
  console.error(`Installation failed: ${err.message}`);
  console.error("Install manually: brew tap create-flow-ai/tap && brew install opengyver");
  process.exit(0); // Don't fail the plugin install
});
