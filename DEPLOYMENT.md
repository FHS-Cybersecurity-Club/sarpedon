# Fairview tryouts deployment

The prepared deployment lives at `/opt/sarpedon` on `ubuntu@ampere-free`.
The stack builds natively on ARM64 and publishes no host ports. MongoDB uses
an internal Docker network and persists its data in the `mongo-data` volume.

## Finish setup

1. Edit `/opt/sarpedon/.env` and set `TUNNEL_TOKEN` to your dedicated,
   remotely managed Cloudflare Tunnel token.
2. Configure the tunnel's published application as
   `scoring.knightscyber.org` → `http://sarpedon:4013`.
3. Replace `Windows-Tryout-PLACEHOLDER` and `Linux-Tryout-PLACEHOLDER` in
   `sarpedon.conf` with the exact `name` values from the Aeacus images.
4. Start the stack:

   ```bash
   cd /opt/sarpedon
   chmod 600 .env sarpedon.conf
   docker compose build
   docker compose up -d
   docker compose ps
   docker compose logs --tail=100 sarpedon cloudflared
   ```

Compose rejects an empty tunnel token. To validate or build before adding it,
use `TUNNEL_TOKEN=validation-only docker compose config --quiet` or
`TUNNEL_TOKEN=validation-only docker compose build sarpedon`; this value is
only for validation/builds, never for starting the tunnel.

The shared scoring password and a separate generated `admin` login password
are in `sarpedon.conf`. This file and `.env` are excluded from Git and the
Docker build context. Treat the generated admin password as private.

## Aeacus settings

For both images, set `remote = "https://scoring.knightscyber.org"` and use
the top-level `password` from `sarpedon.conf`. Each image's `name` must
exactly match its Sarpedon image name, including capitalization.

Each participant's `TeamID.txt` must contain their full name exactly as
listed in `sarpedon.conf`, with spaces. The testing team uses `testing_id`
and appears as `Scoring Test`.

The configured playtime is `2h`, with `enforce = false`, as in the setup
discussion. This records playtime but does not send the client a kill signal
after two hours; manage the end of tryouts separately. Timestamps use
`America/Denver`.

Leave Aeacus's `/update` and `/status/*` routes accessible without an
interactive Cloudflare Access login. Otherwise scoring clients cannot connect.

## Operations

After editing `sarpedon.conf`, run `docker compose restart sarpedon`.
`docker compose down` retains the MongoDB volume; `docker compose down -v`
deletes it and all stored scores. Back up MongoDB before removing its volume.

For a non-Docker installation, `MONGODB_URI` defaults to
`mongodb://localhost:27017`; the Compose stack sets it to `mongodb://mongo:27017`.
