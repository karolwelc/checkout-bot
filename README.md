# Zalando & ASOS Checkout Bot

## Overview
This is a high-performance checkout bot for Zalando and ASOS, written in Golang. The bot automates the checkout process, allowing users to quickly purchase items before they sell out.

## Features
- Fast and efficient checkout process
- Proxy support for avoiding bans
- Multi-threading for handling multiple tasks simultaneously
- Automatic session handling
- Error handling and retry mechanism
- Support for account login and guest checkout

## Installation

### Prerequisites
- Golang (version 1.18+ recommended)
- MongoDB or PostgreSQL (for storing session and order data)
- A set of proxies (optional but recommended)

### Clone the Repository
```sh
git clone https://github.com/karolwelc/checkout-bot.git
cd checkout-bot
```

### Install Dependencies
```sh
go mod tidy
```

## Configuration
Create a `.env` file in the root directory and add your configuration details:
```env
ZALANDO_BASE_URL=https://www.zalando.com
ASOS_BASE_URL=https://www.asos.com
DB_CONNECTION_STRING=mongodb://localhost:27017/botdb
PROXY_LIST=proxy1:port,proxy2:port,proxy3:port
USER_AGENT=your-custom-user-agent
```

## Usage
### Running the Bot
To start the bot, run:
```sh
go run main.go
```

## Roadmap
- [ ] Support for more retailers

## Contributing
Feel free to submit issues or pull requests to improve the bot!

## Disclaimer
This bot is for educational purposes only. Use it at your own risk.
