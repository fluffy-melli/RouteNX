
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
> 현재 디렉토리를 빌드 컨텍스트로 사용하여 **`routenx`**라는 태그를 가진 Docker 이미지를 생성합니다.

---

```sh
docker run -d \
    --restart unless-stopped \
    -p 80:80 -p 443:443 -p 3000:3000 \
    routenx
```

> [!NOTE]  
> **`routenx`** 컨테이너를 **백그라운드(detached mode)**로 실행하며,  
> 수동으로 중지하지 않는 이상 자동으로 재시작됩니다.  
> 포트 **80**, **443**, **3000**을 호스트와 연결합니다.

---

```json
"port": 80,
"ssl-port": 443,
"web-port": 3000
```

> [!NOTE]  
> 포트 **80**은 HTTP, **443**은 HTTPS(SSL), **3000**은 웹 콘솔 인터페이스 용도로 사용됩니다.

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
> `"you@example.com"`과 `"example.com"`을 실제 이메일과 도메인으로 반드시 교체하세요.  
> 유효하지 않은 값일 경우 SSL 인증서 발급이 실패합니다.

> [!TIP]  
> `"testing"`이 `true`로 설정되어 있으면 **Let's Encrypt의 스테이징 환경**에서 테스트용 인증서를 발급합니다.  
> 이는 개발 및 테스트에 유용하며 **요청 제한(rate limit)**을 피할 수 있습니다.  
> 실제 운영 환경에서는 `"testing": false`로 설정하여 실제 인증서를 받으세요.

> [!NOTE]  
> SSL이 활성화되며, HTTPS 트래픽을 위해 **443번 포트**에서 수신합니다.

---

```json
"firewalls": [
  {
    "name": "cloudflare",
    "cidr": [
      "173.245.48.0/20",
      "103.21.244.0/22",
      "... (기타)",
      "2c0f:f248::/32"
    ],
    "block": false
  }
]
```

> [!TIP]  
> 이 방화벽 규칙을 사용하는 라우트는 **Cloudflare를 통해 들어오는 패킷만 허용**합니다.

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
> `*.example.com`에서 들어오는 트래픽을 `localhost:2222`로 전달하며,  
> 오직 **Cloudflare IP로부터의 요청만 허용**합니다.
