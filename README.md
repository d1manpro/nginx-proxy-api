# nginx-proxy-api

**nginx-proxy-api** is a REST API service for automating Nginx reverse proxy configuration, managing SSL certificates via `certbot`, and handling DNS records through Cloudflare.  
It allows automatic provisioning of HTTPS proxying for subdomains or custom domains pointing to specified backend targets.

---

## 🚀 Features

- **Add new proxy**
  - Creates `A` DNS record in Cloudflare pointing to your node IP  
  - Verifies or issues SSL certificates via `certbot`
  - Generates and enables Nginx site configuration
  - Reloads Nginx automatically

- **Remove proxy**
  - Deletes Nginx config and symlink
  - Removes DNS record from Cloudflare (if applicable)
  - Deletes SSL certificate via `certbot`

- **Secure API**
  - Access control by allowed IPs  
  - Token-based authentication  
  - CORS configuration for specific origins  
  - Graceful shutdown and structured logging with `zap`

---

## ▶️ Installation & Setup

### 1. Install NPA using the Installer Script

We provide a convenient installer script that automates binary download, configuration setup, and systemd service creation.

> ⚠️ Warning! Read the configuration and template setup below before running the installer.

```bash
curl -fsSL https://raw.githubusercontent.com/d1manpro/nginx-proxy-api/main/install.sh | sudo bash
````

The installer will:

* Download the latest binary from GitHub Releases
* Create `/usr/local/bin/npapi` and make it executable
* Set up configuration directory `/etc/npapi` with:

  * `config.yml`
  * `template.conf` for Nginx reverse proxy
* Create log file `/var/log/npapi.log`
* Install and enable systemd service `npapi`

After installation, the service will start automatically.

---

### 2. Initial Configuration

After installation, edit `/etc/npapi/config.yml` to set:

* `cloudflare.token` — your Cloudflare API token
* `cloudflare.node_ip` — the IP address for new DNS records
* `cloudflare.domains` — map of domain names to Cloudflare zone IDs
* `email` — your Lets Encrypt email address for CertBot

You can also customize the Nginx template in `/etc/npapi/template.conf`.

---

### 3. Test access

Use the generated API token (printed at the end of installation) for authentication:

```bash
curl -X GET http://localhost:8080/test \
  -H "Authorization: Bearer <your_token>"
```

---

## 🪵 Logging

Uses `zap` with a human-readable console encoder and timestamps in `YYYY.MM.DD HH:MM:SS.mmm` format.

Saves logs into `/var/log/npapi.log`

---

## 🧰 Requirements
* Nginx installed and configured with:

  ```
  /etc/nginx/sites-available/
  /etc/nginx/sites-enabled/
  ```
* `certbot` (with nginx and dns-cloudflare plugins)
* Valid Cloudflare API token with `Zone.DNS` permissions

---

## 🧩 API Endpoints

| Method | Path            | Description         | Auth Required |
| ------ | --------------- | ------------------- | ------------- |
| POST   | `/add-proxy`    | Add proxy config    | ✅            |
| POST   | `/remove-proxy` | Remove proxy config | ✅            |
| GET    | `/test`         | Health check        | ✅            |

---

## Examples

**Add proxy**

```bash
curl -X POST https://api.example.com/add-proxy \
  -H "Authorization: Bearer your_api_token" \
  -d '{"domain": "sub.example.com", "target": "node.example.com:8800"}'
```

**Remove proxy**

```bash
curl -X POST https://api.example.com/remove-proxy \
  -H "Authorization: Bearer your_api_token" \
  -d '{"domain": "sub.example.com"}'
```


## 🛑 Graceful Shutdown

When receiving `SIGINT` or `SIGTERM`, the service:

1. Stops the HTTP server gracefully
2. Closes all open connections
3. Writes shutdown messages to the log

---

## 🧑‍💻 Author

Developed by **[@d1manpro](https://github.com/d1manpro)**. Licensed under MIT License.
