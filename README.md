![Let's Encrypt](https://img.shields.io/badge/SSL-Let's%20Encrypt-orange?style=flat-square&logo=letsencrypt&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=Docker&logoColor=white) ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) ![Vue](https://img.shields.io/badge/-Vue.js-42b883?style=flat-square&logo=vue.js&logoColor=white) ![SCSS](https://img.shields.io/badge/SCSS-CC6699?style=flat-square&logo=sass&logoColor=white)

---

### ⚙️ Build Image

```sh
docker build -t routenx .
```

### 🚀 Run Container

```sh
docker run -d \
    --name routenx \
    --restart unless-stopped \
    -p 80:80 -p 443:443 -p 3000:3000 \
    routenx
```

---

## ⚙️ Configuration

### 🔌 Ports

```json
"port": 80,
"ssl-port": 443,
"web-port": 3000
```

- `80`: HTTP Proxy  
- `443`: HTTPS Proxy (SSL)  
- `3000`: Web Console

---

### 🌐 Routes

```json
"routes": [
  {
    "host": ["*.example.com"],
    "firewall": ["cloudflare"],
    "endpoint": "http://localhost:2222"
  }
]
```

- Routes `*.example.com` to `localhost:2222`  
- Allows only Cloudflare IPs

---

### 🔐 Firewall

```json
"firewall": [
  {
    "name": "cloudflare",
    "cidr": [
      "173.245.48.0/20",
      "103.21.244.0/22",
      "... (etc)",
      "2c0f:f248::/32"
    ],
    "block": false
  }
]
```

- Allows requests only from Cloudflare’s IP ranges

---

### 📄 SSL Configuration

```json
"ssl": {
  "enabled": true,
  "email": "you@example.com",
  "domains": [
    "example.com",
    "sub.example.com"
  ]
}
```

> [!WARNING]
> Don't forget to update `"you@example.com"` and `"example.com"` with your real email and domain.
> Without valid information, SSL certificate generation will fail.

- Automatic SSL (Let's Encrypt)  
- Requires valid email and domain  
- Listens on port `443` for `https`

---

### ✅ Summary

- 🌐 Reverse proxy server with Docker support  
- 🔐 SSL certificate management (Let's Encrypt)  
- 🔥 Domain-based routing & Cloudflare firewall  
- 🧰 **Web admin console available at port `3000`**

---