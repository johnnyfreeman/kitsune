<div align="center">
  <h1><code>kitsune</code></h1>

  <p><b>A concurrent log multiplexer that interleaves real-time log streams with minimal I/O overhead.</b></p>
</div>

> [!WARNING]
> **This project is in its early stages and may undergo many breaking changes in the future.**

## Features
- **Multi-Stream Support** – Monitor multiple log files in real time.  
- **Concurrency-Optimized** – Efficiently reads logs using Goroutines.  
- **Interleaved Output** – Keeps logs readable and in sync across sources.  
- **Handles Rotating Logs** – Detects file changes and continues streaming.  
- **Optional Color Coding** – Distinguish logs visually at a glance.  
- **Regex Filtering** – Tail only what matters. *(Coming soon!)*  

## Installation
```sh
git clone https://github.com/yourusername/kitsune.git
cd kitsune
go build -o kitsune
```
Or install via `go install`:
```sh
go install github.com/yourusername/kitsune@latest
```

## Usage
### Basic Usage
```sh
kitsune app.log db.log server.log
```
This continuously tails and interleaves logs from all specified files.

### Follow Log Files (`tail -f` style)
```sh
kitsune -f /var/log/syslog /var/log/nginx/access.log
```

### Color-Coded Output (Planned)
```sh
kitsune --color app.log db.log
```

### Regex Filtering (Planned)
```sh
kitsune --filter "ERROR|WARN" app.log
```

## Performance
Kitsune is built for speed, using **Goroutines** and **low-latency file polling** to ensure minimal overhead while reading multiple log streams.

## Contributing
Pull requests are welcome! If you have ideas for features, feel free to open an issue.

## License
MIT License. See [LICENSE](LICENSE) for details.
