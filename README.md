# Go-Live

**Contents**
1. [Overview](#overview)
1. [Usage](#usage)

## Overview
Go Live is an [HLS](https://en.wikipedia.org/wiki/HTTP_Live_Streaming) based live-streaming service. After creating an account, users can stream video from [OBS](https://obsproject.com/) (or another prefered streaming software) to Go Live's HLS server. The server will broadcast the stream to anyone who joins it's room through the Go Live web-app. Other cool features include:

* Text chat for every room.
* Users being able to set thier own profile pictures.
* Rewindable streams up to 5 minutes behind.
* Text searching users.

For more on how this project works, visit my [portfolio](https://sidney-bernardin.github.io/project/?id=go_live).

## Usage
Running Go Live locally on your machine is as simple as clonning this repository.

``` bash
git clone https://github.com/Sidney-Bernardin/Go-Live.git
cd Go-Live
```

Then using docker to spin-up Go Live's various services.

``` bash
# The "--profile web_app" spins-up an isolated/volume-less version of the web_app service's container.
docker compose --profile web_app up --build
```
