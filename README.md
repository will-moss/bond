<p align="center">
    <h1 align="center">Bond</h1>
    <p align="center">
      Self-hostable headless QR code generator
      <br />
      Generate QR codes with a one-endpoint API
   </p>
   <p align="center">
      <a href="#table-of-contents">Table of Contents</a> -
      <a href="#deployment-and-examples">Install</a> -
      <a href="#configuration">Configure</a>
    </p>
    <p align="center">
        <a href="#free-hosted-service">
            <img src="https://img.shields.io/badge/FREE%20HOSTED%20SERVICE-AVAILABLE-green?style=for-the-badge&labelColor=black" />
        </a>
    </p>
</p>

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Deployment and Examples](#deployment-and-examples)
  * [Deploy with Docker](#deploy-with-docker)
  * [Deploy with Docker Compose](#deploy-with-docker-compose)
  * [Deploy as a standalone application](#deploy-as-a-standalone-application)
- [Configuration](#configuration)
- [API Reference](#api-reference)
  * [Generate a QR code](#generate-a-qr-code)
- [Usage](#usage)
  * [Curl](#curl)
  * [Wget](#wget)
  * [Javascript](#javascript)
- [Free Hosted Service](#free-hosted-service)
- [Troubleshoot](#troubleshoot)
- [Credits](#credits)

## Introduction

Bond is a tiny, simple, and self-hostable service that enables you to generate QR codes by calling an API.
It was born out of a need to generate QR codes for my business when I couldn't find a fully free and secure API without limitations, along with Google shutting down their service.
I also wanted something rudimentary without gimmicks or customization (colors, redirection, logo, shapes, etc.), hence I decided to make it myself.

## Features

Bond has all these features implemented :
- Generate a QR code of any size with any content
- Simple security using a defined secret to deter bots
- Support for HTTP and HTTPS
- Support for health checks
- Support for standalone / proxy deployment

On top of these, one may appreciate the following characteristics :
- Written in Go
- Holds in a single file with few dependencies
- Holds in a ~14 MB compressed Docker image

For more information, read about [Configuration](#configuration) and [API Reference](#api-reference).

## Deployment and Examples

### Deploy with Docker

You can run Bond with Docker on the command line very quickly.

You can use the following commands :

```sh
# Create a .env file
touch .env

# Edit .env file ...

# Option 1 : Run Bond attached to the terminal (useful for debugging)
docker run --env-file .env -p <YOUR-PORT-MAPPING> mosswill/bond

# Option 2 : Run Bond as a daemon
docker run -d --env-file .env -p <YOUR-PORT-MAPPING> mosswill/bond
```

### Deploy with Docker Compose

To help you get started quickly, multiple example `docker-compose` files are located in the ["examples/"](examples) directory.

Here's a description of every example :

- `docker-compose.simple.yml`: Run Bond as a front-facing service on port 80, with environment variables supplied in the `docker-compose` file directly.

- `docker-compose.volume.yml`: Run Bond as a front-facing service on port 80, with environment variables supplied as a `.env` file mounted as a volume.

- `docker-compose.ssl.yml`:  Run Bond as a front-facing service on port 443, listening for HTTPS requests, with certificate and private key provided as mounted volumes.

- `docker-compose.proxy.yml`: A full setup with Bond running on port 80, behind a proxy listening on port 443.

When your `docker-compose` file is on point, you can use the following commands :
```sh
# Run Bond in the current terminal (useful for debugging)
docker-compose up

# Run Bond in a detached terminal (most common)
docker-compose up -d

# Show the logs written by Bond (useful for debugging)
docker logs <NAME-OF-YOUR-CONTAINER>
```

### Deploy as a standalone application

Deploying Bond as a standalone application assumes the following prerequisites :
- You have Go installed on your server
- You have properly filled your `.env` file
- Your DNS and networking configuration is on point

When all the prerequisites are met, you can run the following commands in your terminal :

```sh
# Retrieve the code
git clone https://github.com/will-moss/bond
cd bond

# Create a new .env file
cp sample.env .env

# Edit .env file ...

# Build the code into an executable
go build -o bond main.go

# Option 1 : Run Bond in the current terminal
./bond

# Option 2 : Run Bond as a background process
./bond &

# Option 3 : Run Bond using screen
screen -S bond
./bond
<CTRL+A> <D>
```

## Configuration

To run Bond, you will need to set the following environment variables in a `.env` file located next to your executable :

> **Note :** Regular environment variables provided on the commandline work too

| Parameter               | Type      | Description                | Default |
| :---------------------- | :-------- | :------------------------- | ------- |
| `SSL`            | `boolean` | Whether HTTPS should be used in place of HTTP. When configured, Bond will look for `certificate.pem` and `key.pem` next to the executable for configuring SSL. Note that if Bond is behind a proxy that already handles SSL, this should be set to `false`. | False        |
| `PORT`           | `integer` | The port Bond listens on. | 80        |
| `SECRET`         | `string`  | The secret used to secure your Bond instance against bots / malicious usage. (This parameter can be left empty to disable security) | a-very-long-and-complicated-secret |
| `MAX_SIZE`       | `integer` | The max size for your QR codes, in pixels, such that a QR code can never be greater than MAX_SIZE x MAX_SIZE pixels. | 1024 |
| `RECOVERY_LEVEL` | `string`  | The recovery level used to generate the QR codes. One of : Low, Medium, High, and Highest (case-insensitive). | Medium |
| `ENABLE_LOGS`    | `boolean` | Whether all the HTTP requests should be displayed in the console / logs. | TRUE |

> **Note :** Boolean values are case-insensitive, and can be represented via "ON" / "OFF" / "TRUE" / "FALSE" / 0 / 1.

> **Tip :** You can generate a random secret with the following command :

```sh
head -c 1024 /dev/urandom | base64 | tr -cd "[:lower:][:upper:][:digit:]" | head -c 32
```

## API Reference

Bond exposes the following API, consisting of a single endpoint :

#### Generate a QR code

```
  GET /
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `secret`  | `string` | **Required.** Your server secret (can be empty if your `SECRET` setting is empty). |
| `size`    | `string` | **Required.** The size (in pixels) of the QR code to generate. (The QR code will be size x size pixels.) |
| `content` | `string` | **Required.** The data to encode in the QR code. |

The API will directly return the image representing the QR code generated using your settings.

#### Perform a health check

```
  GET /health
```

The API will return a `HTTP 200 OK` response with an empty body, indicating that the server is up and running.

## Usage

To generate QR codes using Bond, you can copy and adapt the following examples :

### curl

```sh
curl -o qr-code.png "https://bond.your-domain.tld/?content=YOUR-CONTENT&size=512&secret=YOUR-SECRET"
```

### wget

```sh
wget -O qr-code.png "https://bond.your-domain.tld/?content=YOUR-CONTENT&size=512&secret=YOUR-SECRET"
```

### Javascript

```javascript
async function to_qrcode(text) {
  const url = `https://bond.your-domain.tld/?` + new URLSearchParams({
    size: 512,
    content: text,
    secret: 'YOUR-SECRET'
  });
  let response = await fetch(url);

  if (response.status !== 200) {
    console.log('HTTP-Error: ' + response.status);
    return null;
  }

  const blob = await response.blob();
  const objectURL = URL.createObjectURL(blob);

  const image = document.createElement('img');
  image.src = objectURL;

  const container = document.getElementById('YOUR-CONTAINER');
  container.append(image);
}

await to_qrcode("YOUR-CONTENT");
```

## Free Hosted Service

Using our hosted service, you can use Bond to generate QR codes for free and without account registration.

The service is provided freely with the following characteristics :
- **URL :** `https://endpoint.bond/`
- **Secret :** `52e679fae92441942a2ed4390ad9e8639eab9347a74a19ebaa00ef4a5494f7f3`
- **Limitations :**
  * **Max QR code size :** 512x512 pixels
  * **Recovery level :** Low
  * **HTTP verbs :** OPTIONS, GET
  * **Rate limit :** 1 request per second
- **Miscellaneous :**
  * Requests are not logged
  * HTTPS is required

For more advanced needs, please open an issue, or send an email to the address displayed on my Github Profile.

## Troubleshoot

Should you encounter any issue running Bond, please refer to the following common problems that may occur.

> If none of these matches your case, feel free to open an issue.

#### Bond is unreachable over HTTP / HTTPS

Please make sure that the following requirements are met :

- If Bond runs as a standalone application without proxy :
    - Make sure your server / firewall accepts incoming connections on Bond's port.
    - Make sure your DNS configuration is correct. (Usually, such record should suffice : `A bond XXX.XXX.XXX.XXX` for `https://bond.your-server-tld`)
    - Make sure your `.env` file is well configured according to the [Configuration](#configuration) section.

- If Bond runs behind Docker / a proxy :
    - Perform the previous (standalone) verifications first.
    - Make sure that `PORT` (Bond's port) is well set in `.env`.
    - Check your proxy forwarding rules.

In any case, the crucial part is [Configuration](#configuration).

#### Bond returns an error 4xx instead of a QR code

Please make sure that :
- You're using the `GET` HTTP method.
- You've included the `secret` parameter, and the value of it equals the value of the `SECRET` defined in your `.env`.
- The `size` you requested fits within the range 1 <= `size` <= `MAX_SIZE`.


#### Something else

Please feel free to open an issue, explaining what happens, and describing your environment.

## Credits

Hey hey ! It's always a good idea to say thank you and mention the people and projects that help us move forward.

Big thanks to the individuals / teams behind these projects :
- [go-qrcode](https://github.com/skip2/go-qrcode) : For the QR code generation.
- [go-chi](https://github.com/go-chi/chi) : For the web server.
- The countless others!

And don't forget to mention Bond if you like it or if it helps you in any way!
