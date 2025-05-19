# Ludo Backend Setup on aaPanel

Here's a guide to set up this Ludo game backend on aaPanel:

## Prerequisites

- aaPanel installed on your server
- MongoDB installed (can be done through aaPanel)
- Go installed (version 1.16+)

## Installation Steps

### 1. Install MongoDB

If MongoDB is not already installed:
1. In aaPanel, go to the App Store
2. Search for MongoDB and install it
3. Create a database named "ludo"

### 2. Upload the Ludo Backend

1. Create a directory for the project in your web server root
2. Upload all project files maintaining the directory structure:
   ```
   go.mod
   go.sum
   main.go
   app/
   database/
   models/
   tmp/
   utils/
   ```

### 3. Configure MongoDB Connection

Check and update the MongoDB connection in `database/mongo.go` if needed:

```go
// Update connection string if necessary
```

### 4. Build the Application

Connect to your server via SSH and run:

```bash
cd /path/to/ludo_backend
go mod tidy
go build -o ludo_server main.go
```

### 5. Set up as a Service

Create a systemd service file to run the application:

```bash
sudo nano /etc/systemd/system/ludo-backend.service
```

Add the following content:

```
[Unit]
Description=Ludo Game Backend
After=network.target mongodb.service

[Service]
Type=simple
User=www
WorkingDirectory=/path/to/ludo_backend
ExecStart=/path/to/ludo_backend/ludo_server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable ludo-backend
sudo systemctl start ludo-backend
```

### 6. Configure Reverse Proxy in aaPanel

1. In aaPanel, go to Website management
2. Select your website or create a new one
3. Set up a reverse proxy to forward requests to the Go application (default port is likely 8080, check `main.go` for the actual port)

## API Endpoints

- WebSocket endpoint: `/ws?userId=<user_id>`
- Game creation happens automatically when enough players join a room

## Testing the Setup

You can test the WebSocket connection using tools like Postman or a simple web client.

## Troubleshooting

- Check logs with: `sudo journalctl -u ludo-backend.service`
- If changes are needed, update the code, rebuild, and restart the service
- Ensure MongoDB is running and accessible