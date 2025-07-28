# Wallet Tracker

A comprehensive cryptocurrency wallet tracking application built with Go, designed to help users monitor their digital assets across multiple blockchain networks.

## 🚀 Features

- **Multi-Chain Support**: Track wallets on Ethereum, BSC, and Polygon networks
- **Token Management**: Add and monitor custom tokens for each wallet
- **Real-time Balance Tracking**: Get up-to-date token balances with USD valuations
- **User Authentication**: Secure JWT-based authentication system
- **Caching System**: Redis-powered caching for improved performance
- **RESTful API**: Clean and well-documented API endpoints

## 🏗️ Architecture

The application follows a clean architecture pattern with clear separation of concerns:

```
wallet-tracker/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/         # Configuration management
│   ├── handler/        # HTTP request handlers
│   ├── middleware/     # HTTP middleware (auth, etc.)
│   ├── model/          # Data models and entities
│   ├── repository/     # Data access layer
│   └── service/        # Business logic layer
├── pkg/                # Public packages
│   ├── blockchain/     # Blockchain integration
│   ├── cache/          # Caching utilities
│   └── database/       # Database connection management
├── configs/            # Configuration files
└── docs/              # Documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.24.5
- **Framework**: Gin (HTTP web framework)
- **Database**: MySQL with GORM ORM
- **Cache**: Redis
- **Authentication**: JWT
- **Blockchain**: Ethereum, BSC, Polygon integration
- **Configuration**: Viper + YAML

## 📋 Prerequisites

- Go 1.24.5 or higher
- MySQL 8.0 or higher
- Redis 6.0 or higher
- Infura API key (for Ethereum mainnet access)

## 🚀 Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd wallet-tracker
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Configure the application**

   ```bash
   # Edit configs/config.yaml with your settings
   # Update RPC URLs and API keys
   ```

5. **Set up the database**

   ```bash
   # Create MySQL database
   mysql -u root -p
   CREATE DATABASE wallet_tracker;
   ```

6. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## ⚙️ Configuration

The application uses a YAML configuration file located at `configs/config.yaml`:

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  name: wallet_tracker

redis:
  host: localhost
  port: 6379
  db: 0
  password: ""

blockchain:
  ethereum:
    rpc_url: "https://mainnet.infura.io/v3/YOUR_INFURA_KEY"
  bsc:
    rpc_url: "https://bsc-dataseed1.binance.org/"
  polygon:
    rpc_url: "https://polygon-rpc.com/"

cache:
  token_balance_ttl: 24h
```

## 📚 API Endpoints

### Authentication

- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login

### Wallet Management (Protected)

- `POST /api/v1/wallets` - Add a new wallet
- `GET /api/v1/wallets` - Get user's wallets
- `POST /api/v1/wallets/:wallet_id/tokens` - Add token to wallet
- `GET /api/v1/balances` - Get wallet balances
- `POST /api/v1/refresh-cache` - Refresh cached data

## 🔐 Authentication

The application uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 💾 Database Schema

### Users

- `id` - Primary key
- `username` - Unique username
- `email` - Unique email address
- `password` - Hashed password
- `created_at`, `updated_at` - Timestamps

### Wallets

- `id` - Primary key
- `user_id` - Foreign key to users
- `address` - Wallet address
- `chain_id` - Blockchain network ID
- `chain_name` - Blockchain network name
- `name` - User-defined wallet name
- `created_at`, `updated_at` - Timestamps

### Wallet Tokens

- `id` - Primary key
- `wallet_id` - Foreign key to wallets
- `token_address` - Token contract address
- `token_symbol` - Token symbol
- `token_name` - Token name
- `token_decimals` - Token decimal places
- `is_active` - Token tracking status
- `created_at`, `updated_at` - Timestamps

## 🔄 Caching Strategy

The application uses Redis for caching:

- **Token balances**: Cached for 24 hours by default
- **User sessions**: JWT token validation
- **Blockchain data**: RPC call results

## 🧪 Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## 🚀 Deployment

### Docker (Recommended)

```bash
# Build the image
docker build -t wallet-tracker .

# Run the container
docker run -p 8080:8080 wallet-tracker
```

### Manual Deployment

1. Build the binary: `go build -o wallet-tracker cmd/server/main.go`
2. Set up environment variables
3. Run the binary: `./wallet-tracker`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/wallet-tracker/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

## 🔮 Roadmap

- [ ] Support for more blockchain networks (Solana, Avalanche, etc.)
- [ ] Real-time price alerts
- [ ] Portfolio analytics and charts
- [ ] Mobile application
- [ ] Webhook notifications
- [ ] Multi-language support
