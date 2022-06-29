# aspenjames.dev -- Infrastructure

Infrastructure supporting aspenjames.dev:

## Google Compute Instance

This VM hosts the Go web service & the website content. The Go service runs as a
Docker container, managed through systemd. The website content is stored in the
home directory and mounted to the web service container by the systemd daemon.
The web server & web content may be managed independently.

New versions of the Go web service may be deployed by:

- Building & pushing a new container image
- Triggering a systemd service restart on the instance

New versions of the website content may be deployed by:

- Building CSS & any static files
- Syncing content files to the instance
- Triggering a systemd service restart

## Fastly Service

(TBD) Sets up a CDN distribution for the content hosted on the compute instance
& handles TLS termination.
