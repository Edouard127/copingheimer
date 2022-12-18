# Copingheimer

Copingheimer is a simple tool for coping with the fact that you can't use copenheimer :3

## Features

- [x] Cope
- [x] Blacklist. Avoid getting enterprises getting mad and sending you letters
- [x] Concurrency. Cope with multiple IPs at once

## Usage

```bash
$ copingheimer -h
Usage of copingheimer:
  -bf, -blacklist-file string
        Path to the blacklist file (default "blacklist.txt")
  -c, -config string
        Path to the config file (default "config.env")
  -cs, -cpu-saver
        Whether to enable the CPU saver or not (default: true) (default true)
  -d, -database string
        Database to use (default: mongodb) (mongodb, bolt) (default "mongodb")
  -du, -database-url string
        URL to the database (default "mongodb://localhost:27017")
  -h, -help Show this help message
  -i, -instances int
        Number of instances to run (default: 1) (default 1)
  -ip string
        IP address to start from, only used in order mode (default "0.0.0.0")
  -m, -mode string
        Mode to run in (default: "random") (random, order) (default "random")
  -t, -timeout int
        Timeout for each ping (default: 1000) (default 1000)
```

## Risks
- [ ] You might get sued or get a letter by enterprises if you keep scanning their IPs
- [ ] You might get banned from your ISP

I am not responsible for any of the above or any other damages caused by this tool.

Although, you can reduce these risks by:
- [x] Using a VPN
- [x] Using a VPS
- [x] Using a Tor exit node
- [x] Blacklisting IP ranges that are susceptible to getting you in trouble
- [x] Using the random mode instead of the order mode
