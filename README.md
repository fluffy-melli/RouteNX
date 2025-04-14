
<div align="center">
  <h3>
    <a href="https://github.com/fluffy-melli/RouteNX/blob/main/README.KR.md">KR</a> /
    <a href="https://github.com/fluffy-melli/RouteNX/">EN</a>
  </h3>
</div>

```sh
docker build -t routenx .
```

> [!NOTE]  
> Builds a Docker image tagged **`routenx`**, using the current directory as the build context.

---

```sh
docker run -d \
    --restart unless-stopped \
    -p 80:80 -p 443:443 -p 3000:3000 \
    routenx
```

> [!NOTE]  
> Runs the **`routenx`** container in detached mode,  
> restarts automatically unless stopped,  
> and maps ports **80**, **443**, and **3000** to the host.

---

```json
"port": 80,
"ssl-port": 443,
"web-port": 3000
```

> [!NOTE]  
> Port **80** is for HTTP, **443** for HTTPS (SSL), and **3000** is used by the Web Console interface.

---

```json
"ssl": {
  "enabled": true,
  "testing": true,
  "email": "you@example.com",
  "domains": [
    "example.com",
    "sub.example.com"
  ]
}
```

> [!WARNING]  
> Be sure to replace `"you@example.com"` and `"example.com"` with your actual email and domain.  
> SSL certificate generation will fail if these values are invalid.

> [!TIP]  
> When `"testing"` is set to `true`, SSL certificates will be requested from **Let's Encrypt's staging environment**.  
> This is useful for development and testing because it avoids **rate limits**.  
> Set `"testing": false` for production to obtain real certificates.

> [!NOTE]  
> Enables SSL and listens on port **443** for HTTPS traffic.

---

```json
"firewalls": [
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

> [!TIP]  
> Routes using this firewall rule **only accept packets coming through Cloudflare.**

---

```json
"routes": [
  {
    "host": [
      "*.example.com"
    ],
    "firewall": [
      "cloudflare"
    ],
    "endpoint": "http://localhost:2222"
  }
]
```

> [!NOTE]  
> Routes traffic from `*.example.com` to `localhost:2222`,  
> and only accepts requests from **Cloudflare IPs**.
