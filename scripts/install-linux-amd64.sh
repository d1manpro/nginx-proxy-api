#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="npapi"
BIN_URL="https://github.com/d1manpro/nginx-proxy-api/releases/download/v1.0.0/npa-linux-amd64"
BIN_PATH="/usr/local/bin/${SERVICE_NAME}"
CONFIG_DIR="/etc/${SERVICE_NAME}"
SYSTEMD_UNIT="/etc/systemd/system/${SERVICE_NAME}.service"
LOG_FILE="/var/log/${SERVICE_NAME}.log"

echo "[*] Installing ${SERVICE_NAME}..."

if [[ $EUID -ne 0 ]]; then
  echo "[-] This script must be run as root" >&2
  exit 1
fi

echo "[*] Downloading binary from GitHub..."
curl -fsSL -o "${BIN_PATH}" "${BIN_URL}"
chmod +x "${BIN_PATH}"
chown root:root "${BIN_PATH}"

echo "[*] Setting up configs..."
mkdir -p "${CONFIG_DIR}"

TOKEN=$(cat /proc/sys/kernel/random/uuid)
echo "-> Generated API token"

if [[ ! -f "${CONFIG_DIR}/config.yml" ]]; then
  echo "  -> Creating default config.yml at ${CONFIG_DIR}/"
  cat > "${CONFIG_DIR}/config.yml" <<EOF
http_server:
  host: "0.0.0.0"
  port: 8080
  origins: ["*"]

access:
  token: "${TOKEN}"
  allowed_ips: ["127.0.0.1", "::1"]

cloudflare:
  token: "your_cloudflare_api_token"
  node_ip: "0.0.0.0"
  domains:
    "example.com": "cloudflare_zone_id"

email: "admin@example.com"

debug_mode: false
EOF
else
  echo "  -> Config already exists..."
fi
  chmod 600 "${CONFIG_DIR}/config.yml"

if [[ ! -f "${CONFIG_DIR}/template.conf" ]]; then
  echo "  -> Creating default template.conf at ${CONFIG_DIR}/"
  cat > "${CONFIG_DIR}/template.conf" <<'EOF'
server {
    listen 80;
    server_name {{.Domain}};
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name {{.Domain}};

    ssl_certificate /etc/letsencrypt/live/{{.Cert}}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/{{.Cert}}/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;

    location / {
        proxy_pass http://{{.Target}}/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF
else
  echo "  -> Template already exists..."
fi
  chmod 644 "${CONFIG_DIR}/template.conf"

touch "${LOG_FILE}"
chmod 644 "${LOG_FILE}"

echo "[*] Creating systemd service..."
cat > "${SYSTEMD_UNIT}" <<EOF
[Unit]
Description=Nginx Proxy API (npapi) service
After=network.target nginx.service
Wants=nginx.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=${CONFIG_DIR}
ExecStart=${BIN_PATH}
Restart=on-failure
RestartSec=5s

StandardOutput=append:${LOG_FILE}
StandardError=append:${LOG_FILE}

NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

echo "[*] Enabling and starting service..."
systemctl daemon-reload
systemctl enable --now "${SERVICE_NAME}"

echo "[+] Installation complete!"
echo "  Binary: ${BIN_PATH}"
echo "  Configs: ${CONFIG_DIR}/"
echo "  Log: ${LOG_FILE}"
echo
echo "[!] Your API token (save it safely):"
echo "    ${TOKEN}"
