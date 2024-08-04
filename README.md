# Netrunner: Building HTTP from the Ground Up

![MEME](./images/meme.png)

## Overview

Netrunner is an educational project that demonstrates how to build an HTTP server from scratch using Go. This series covers everything from basic TCP connections to a fully functional HTTP/HTTPS server with advanced features.

If you liked it please make sure to star the repository and follow me for more such walkthroughs!!!!

## Series Contents

1. **Part 0: Foundations and Prerequisites**
   - Setting up the Go environment
   - Introduction to TCP/IP and HTTP basics

2. **Part 1: TCP Foundations**
   - Implementing a basic TCP server
   - Handling connections

3. **Part 2: HTTP Basics - Parsing Requests**
   - Understanding HTTP request structure
   - Implementing request parsing

4. **Part 3: HTTP Responses and Status Codes**
   - Creating HTTP responses
   - Implementing status codes

5. **Part 4: HTTP Methods and Routing**
   - Implementing GET and POST methods
   - Creating a basic routing system

6. **Part 5: Middleware and Static File Serving**
   - Implementing middleware functionality
   - Adding support for serving static files

7. **Part 6: Performance Optimization and Error Handling**
   - Implementing connection pooling
   - Enhancing error handling

8. **Part 7: Implementing HTTPS**
   - Understanding HTTPS and TLS
   - Adding HTTPS support to the server

## Getting Started

To run this project:

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/netrunner.git
   ```

2. Navigate to the project directory:
   ```
   cd netrunner
   ```

3. Run the server:
   ```
   go build -o netrunner cmd/server/main.go
   ./netrunner
   ```

## Features

- HTTP and HTTPS support
- Custom routing with path parameters
- Middleware support
- Static file serving
- JSON request/response handling
- Connection pooling
- Rate limiting
- CORS support
- Graceful shutdown

## Contributing

Contributions to Netrunner are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Minami Sato

- Twitter: [@minamisatokun](https://x.com/minamisatokun)
- Blog: [Minami's Blog](https://minami.bearblog.dev/blog/)
- Substack: [Minami Sato's Substack](https://minamisato.substack.com/)

## Acknowledgments

- Thanks to all the readers and contributors who followed along with this series.
- Special thanks to the Go community for their excellent documentation and resources.