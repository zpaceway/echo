# Echo Service

A lightweight HTTP service that returns the client's IP address. This service is designed to help identify the real client IP address behind proxies, load balancers, and other network infrastructure.

## Features

- **Smart IP Detection**: Automatically detects client IP from multiple headers in order of priority:
  - `X-Original-Forwarded-For`
  - `X-Real-IP`
  - `X-Forwarded-For`
  - Falls back to `RemoteAddr` if no headers are present
- **Proxy-Aware**: Handles comma-separated IP lists from proxy chains
- **Lightweight**: Minimal Go HTTP server with no external dependencies
- **Containerized**: Ready-to-deploy Docker image
- **Multi-Architecture**: Supports both AMD64 and ARM64 architectures

## Quick Start

### Using Docker

Pull and run the pre-built image:

```bash
docker run -p 6868:6868 zpaceway/echo:latest
```

### Building from Source

1. Clone the repository:

```bash
git clone <repository-url>
cd echo
```

2. Build and run:

```bash
go build -o echo main.go
./echo
```

3. Test the service:

```bash
curl http://localhost:6868
```

## API

### GET /

Returns the detected client IP address as plain text.

**Response:**

- **Content-Type**: `text/plain`
- **Status Code**: `200 OK`
- **Body**: Client IP address

**Example:**

```bash
$ curl http://localhost:6868
192.168.1.100
```

## Configuration

The service runs on `0.0.0.0:6868` by default. This is currently hardcoded but can be modified in the source code.

## IP Detection Logic

The service uses the following priority order to determine the client IP:

1. **X-Original-Forwarded-For**: Custom header often used by specific proxy configurations
2. **X-Real-IP**: Common header set by reverse proxies like Nginx
3. **X-Forwarded-For**: Standard header for forwarded requests (takes the last IP in comma-separated list)
4. **RemoteAddr**: Direct connection IP (fallback when no proxy headers present)

For comma-separated IP lists (common in `X-Forwarded-For`), the service returns the rightmost (last) IP address after trimming whitespace.

## Docker

### Building

Build for multiple architectures:

```bash
make build-push
```

This builds and pushes images for both `linux/amd64` and `linux/arm64` platforms.

### Manual Build

```bash
docker build -t echo-service .
docker run -p 6868:6868 echo-service
```

## Use Cases

- **Load Balancer Health Checks**: Verify client IP detection through infrastructure
- **Debugging Network Configuration**: Test proxy and load balancer IP forwarding
- **Security Auditing**: Validate that real client IPs are properly preserved
- **Development Testing**: Simple service for testing network setups

## Security Considerations

- The service trusts proxy headers (`X-Forwarded-For`, etc.) which can be spoofed by malicious clients
- Deploy behind trusted proxies/load balancers in production environments
- Consider additional validation if used for security-critical IP-based decisions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]
