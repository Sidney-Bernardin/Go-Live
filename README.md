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
## How to Use
1. Go to ```localhost:5173``` in your browser.

2. Sign-up for a new account.

   ![alt text](https://github.com/Sidney-Bernardin/Go-Live/blob/main/images/image1.png)

3. Click the Go Live button in the dropdown menu.

   ![alt text](https://github.com/Sidney-Bernardin/Go-Live/blob/main/images/image2.png)

4. Use OBS (or a similar software) to stream with the URI and key.

   ![alt text](https://github.com/Sidney-Bernardin/Go-Live/blob/main/images/image3.png)
   
   You can also stream a video file using FFMPEG: ```ffmpeg -re -i "<your video file>" -c:v copy -c:a aac -ar 44100 -ac 1 -f flv "<URI and key here>"```

5. Get you and your friends visit your profile page at (/your_username) to watch and chat about your live stream.

   ![alt text](https://github.com/Sidney-Bernardin/Go-Live/blob/main/images/image4.png)
