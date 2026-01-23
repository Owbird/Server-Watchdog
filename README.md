# 🛡️ Server-Watchdog

**Server-Watchdog** is a real-time SSH activity monitor for Linux servers.  
It gives you a live, visual overview of who is connecting to your server, where they are coming from, and how long they stay connected — all inside a rich, interactive terminal UI.

---

## What It Does

Server-Watchdog continuously watches port **22 (SSH)** and builds a timeline of activity by:

- Detecting **live SSH connections**
- Tracking **session duration and history**
- Resolving **IP geolocation** in real time
- Persisting historical attempts across restarts
- Visually highlighting suspicious behavior using **color-coded “heat”**

---

## Installation

```bash
git clone https://github.com/owbird/server-watchdog.git
cd server-watchdog
go build -o server-watchdog
```

---

## Usage

```bash
./server-watchdog
```

You’ll see a live dashboard showing:

| Column      | Description               |
| ----------- | ------------------------- |
| `#`         | Entry index               |
| `IP`        | Remote IP address         |
| `Status`    | LIVE or NIL               |
| `Last seen` | NOW or relative time      |
| `Country`   | Geo-resolved country      |
| `Sessions`  | Attempts + total duration |

The display refreshes automatically every second.

---

## Configuration

### Whitelisting IPs

Create or edit `whitelist.json`:

```json
[
  "192.168.1.1",
  "10.0.0.5"
]
```

Whitelisted IPs:

* Are ignored by the watchdog
* Still visible in the UI panel

---

### SSH History Storage

`ssh-attempts.json` is automatically managed.

* Only **completed** sessions are saved
* Live sessions are excluded until they end
* Safe to delete if you want a clean slate

---

## How It Works

1. Reads active SSH connections using:

   ```bash
   ss -tnp | grep :22
   ```
2. Extracts remote IPs
3. Compares against whitelist
4. Updates session state:

   * Starts new sessions
   * Ends inactive ones
5. Resolves country (once per IP)
6. Renders everything in a live terminal UI


## ⚠️ Limitations

* Linux only
* Relies on `ss` output format
* IP geolocation depends on a third-party service
* Not a replacement for fail2ban or firewall rules

---


## Contributing

PRs welcome

---
