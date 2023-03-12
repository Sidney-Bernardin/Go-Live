# Go-Live

## Overview
Go Live is a fullstack live-streaming web-app. Users can create accounts and live-stream to other Go Live users, by using OBS (or similar software) to stream live video to Go Live's RTMP server. Each room has a live chat so viewers to talk to each other.

## Installation
1. Download and enter this repository.

   ``` bash
   git clone https://github.com/Sidney-Bernardin/Go-Live.git
   cd Go-Live
   ```
2. Run services in development mode with Docker Compose.

   ``` bash
   docker compose --profile web_app up --build
   ```
