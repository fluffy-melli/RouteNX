# Running with Docker

### 1. **Build the Image:**

To build the Docker image, use the following command:

```sh
docker build -t routenx .
```

This command will create an image named `routenx` using the Dockerfile in the current directory.

### 2. **Run the Container:**

After building the image, you can run the container with this command:

```sh
docker run -d -p 8080:8080 routenx
```

---

# Configuration Example

This configuration defines a proxy server with specific routing and firewall rules. Below is an explanation of how each part works.

### Configuration Breakdown

#### 1. **Port**
- The proxy server listens on port `8080`.

```json
"port": 8080
```

#### 2. **Routes**
This section defines the routing rules for incoming requests.

- **Route 1:**
  - **Host:** Accepts requests to `*.example.com` and any host on port `8080`.
  - **Firewall:** Only requests from Cloudflare IP ranges are allowed.
  - **Endpoint:** Routes to `http://localhost:2222`.

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

#### 3. **Firewall Rules**
- **Cloudflare IP ranges:** This firewall rule ensures only requests from Cloudflare’s IP addresses are allowed to access the proxy server.

```json
"firewall": [
  {
    "name": "cloudflare",
    "cidr": [
      "173.245.48.0/20",
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "141.101.64.0/18",
      "108.162.192.0/18",
      "190.93.240.0/20",
      "188.114.96.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "162.158.0.0/15",
      "104.16.0.0/13",
      "104.24.0.0/14",
      "172.64.0.0/13",
      "131.0.72.0/22",
      "2400:cb00::/32",
      "2606:4700::/32",
      "2803:f800::/32",
      "2405:b500::/32",
      "2405:8100::/32",
      "2a06:98c0::/29",
      "2c0f:f248::/32"
    ],
    "block": false
  }
]
```

---

### Summary

This configuration sets up a reverse proxy server that:

- Listens on port `8080`.
- Routes traffic to `http://localhost:2222` for specific domains (`*.example.com`).
- Enforces firewall rules that only allow traffic from Cloudflare's IP ranges.

Feel free to adjust the domain names, IP ranges, and endpoints to match your infrastructure. This setup ensures secure and controlled routing through the reverse proxy.
