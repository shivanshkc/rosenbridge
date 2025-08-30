# Rosenbridge Documentation

## Overview

Rosenbridge is a distributed websocket hub written in Go that enables real-time communication between clients across multiple nodes in a cluster. It provides a scalable solution for routing messages between connected users, even when they are connected to different nodes in the cluster.

## Architecture

### Core Components

#### 1. Distributed Cluster Architecture
- **Node Discovery**: Nodes discover each other through a shared MongoDB connection
- **Service Registration**: Each node registers itself in the MongoDB `bridges` collection upon startup
- **Peer Communication**: Nodes communicate with each other via HTTP/HTTPS for message routing

#### 2. Client Connection Management
- **WebSocket Bridges**: Each client connection is represented as a "bridge" with a unique ID
- **Client Identification**: Clients are identified by a `client_id` parameter during connection
- **Connection Tracking**: All active connections are stored in MongoDB for cluster-wide visibility

#### 3. Message Routing System
- **Cross-Node Routing**: Messages can be sent between clients connected to different nodes
- **Database Lookup**: MongoDB is queried to find which node hosts the target client
- **Direct Delivery**: Messages are routed directly to the appropriate node and delivered to the client

## Key Features

### 1. Scalable Architecture
- **Horizontal Scaling**: Add new nodes to the cluster without downtime
- **Load Distribution**: Clients can connect to any available node
- **Fault Tolerance**: Node failures don't affect clients on other nodes

### 2. Real-time Communication
- **WebSocket Protocol**: Low-latency bidirectional communication
- **Message Delivery**: Guaranteed delivery to online clients
- **Offline Handling**: Graceful handling of offline/disconnected clients

### 3. Cloud-Native Support
- **Google Cloud Run**: Built-in support for Cloud Run deployment
- **Local Development**: Solo mode for single-node development
- **Container Ready**: Docker support for containerized deployments

## System Flow

### Node Startup Process
1. **Configuration Loading**: Node loads configuration from `configs.yaml`
2. **MongoDB Connection**: Establishes connection to shared MongoDB instance
3. **Index Creation**: Creates necessary database indexes for performance
4. **Service Discovery**: Registers node address in the cluster
5. **HTTP Server**: Starts HTTP server with WebSocket upgrade capability

### Client Connection Flow
1. **WebSocket Handshake**: Client initiates WebSocket connection with `client_id` parameter
2. **Bridge Creation**: Node creates a unique bridge (connection) with UUID
3. **Database Registration**: Bridge information stored in MongoDB `bridges` collection
4. **Connection Upgrade**: HTTP connection upgraded to WebSocket protocol
5. **Message Handler**: Bridge configured with message handling capabilities

### Message Routing Flow
1. **Message Reception**: Node receives message from connected client
2. **Target Resolution**: MongoDB queried to find target client's node location
3. **Route Determination**: System determines if target is local or remote
4. **Message Delivery**: 
   - **Local**: Direct delivery to local bridge
   - **Remote**: HTTP request sent to target node's internal API
5. **Response Handling**: Delivery status reported back to sender

## Database Schema

### Bridges Collection
```javascript
{
  "client_id": "user123",           // Client identifier
  "bridge_id": "uuid-string",       // Unique bridge identifier
  "node_addr": "http://node1:8080", // Node hosting this connection
  "connected_at": "2024-01-01T00:00:00Z" // Connection timestamp
}
```

### Database Indexes
- `client_id`: B-tree index for client lookups
- `bridge_id`: B-tree index for bridge identification
- `node_addr`: B-tree index for node-based queries

## API Endpoints

### External API (`/api`)

#### GET `/api/bridge?client_id={id}`
- **Purpose**: Establishes WebSocket connection for a client
- **Parameters**: `client_id` (required) - Unique identifier for the client
- **Response**: WebSocket upgrade with custom headers
- **Headers**: 
  - `x-bridge-id`: Unique bridge identifier
  - `x-node-addr`: Node address hosting the connection

#### POST `/api/message`
- **Purpose**: Send message to target clients
- **Authentication**: None (external API)
- **Body**: OutgoingMessageReq format
- **Response**: Delivery status for each target

### Internal API (`/api/internal`)

#### POST `/api/internal/message`
- **Purpose**: Inter-node message delivery
- **Authentication**: Basic Auth (internal_username/internal_password)
- **Body**: OutgoingMessageInternalReq format
- **Response**: Delivery status and errors

## Configuration

### Application Settings
```yaml
application:
  name: rosenbridge
  version: 2.0.0
  solo_mode: false  # Set to true for single-node operation

auth:
  internal_username: dev  # Username for inter-node communication
  internal_password: dev  # Password for inter-node communication

bridges:
  max_bridge_limit: 10000          # Maximum total bridges per node
  max_bridge_limit_per_client: 10  # Maximum bridges per client

http_server:
  addr: 0.0.0.0:8080                    # Server listen address
  discovery_addr: http://0.0.0.0:8080   # External discovery address

logger:
  level: info  # Logging level (debug, info, warn, error)

mongo:
  addr: mongodb://dev:dev@localhost:27017/?retryWrites=true&w=majority
  database_name: rosenbridge
  operation_timeout_sec: 60
```

## Deployment

### Local Development
1. **Prerequisites**:
   - Go 1.17+
   - MongoDB instance
   - Docker (optional)

2. **Setup**:
   ```bash
   # Clone repository
   git clone https://github.com/shivanshkc/rosenbridge
   cd rosenbridge
   
   # Configure MongoDB connection in configs.yaml
   # Start MongoDB (if not already running)
   
   # Run the application
   go run main.go
   ```

3. **Solo Mode**:
   - Set `solo_mode: true` in configuration
   - Uses local in-memory storage instead of MongoDB
   - Suitable for development and testing

### Production Deployment

#### Google Cloud Run
1. **Build and Deploy**:
   ```bash
   # Build Docker image
   docker build -t rosenbridge .
   
   # Deploy to Cloud Run
   gcloud run deploy rosenbridge --image gcr.io/PROJECT-ID/rosenbridge
   ```

2. **Environment Variables**:
   - `K_SERVICE`: Automatically set by Cloud Run
   - MongoDB connection string in configs.yaml

#### Multi-Node Cluster
1. **Shared MongoDB**: All nodes must connect to the same MongoDB instance
2. **Load Balancer**: Use a load balancer to distribute client connections
3. **Service Discovery**: Nodes automatically discover each other via MongoDB
4. **Health Checks**: Implement health check endpoints for monitoring

## Message Formats

### Outgoing Message Request
```json
{
  "sender_id": "user123",
  "receiver_ids": ["user456", "user789"],
  "message": {
    "type": "chat",
    "content": "Hello, world!",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

### Bridge Message (WebSocket)
```json
{
  "type": "outgoing_message_res",
  "request_id": "uuid-string",
  "body": {
    "code": "OK",
    "reason": "Message delivered successfully",
    "report": {
      "user456": [
        {
          "bridge_id": "bridge-uuid",
          "code": "OK",
          "reason": "Delivered"
        }
      ]
    }
  }
}
```

## Error Handling

### Common Error Codes
- `OK`: Operation successful
- `OFFLINE`: Target client is not connected
- `BAD_REQUEST`: Invalid request format or parameters
- `INTERNAL_ERROR`: Server-side error occurred

### Connection Cleanup
- **Graceful Closure**: Proper cleanup when client disconnects
- **Database Cleanup**: Automatic removal of bridge records
- **Stale Connection Handling**: System detects and cleans stale entries

## Monitoring and Logging

### Logging Features
- **Structured Logging**: JSON-formatted logs using Zap logger
- **Request Tracing**: Unique request IDs for request correlation
- **Error Tracking**: Detailed error logging with stack traces
- **Performance Metrics**: Connection and message delivery statistics

### Health Monitoring
- **Connection Counts**: Track active connections per node
- **Message Throughput**: Monitor message delivery rates
- **Database Performance**: MongoDB operation latency tracking
- **Node Status**: Cluster node health and availability

## Security Considerations

### Authentication
- **Internal API**: Basic authentication for inter-node communication
- **WebSocket Origin**: Origin checking for WebSocket connections
- **Client Validation**: Input validation for all client parameters

### Network Security
- **HTTPS/WSS**: Use secure connections in production
- **Firewall Rules**: Restrict internal API access to cluster nodes only
- **Secret Management**: Secure storage of authentication credentials

## Performance Optimization

### Database Optimization
- **Connection Pooling**: Efficient MongoDB connection management
- **Index Usage**: Optimized queries using proper indexes
- **Batch Operations**: Bulk operations for improved performance

### WebSocket Optimization
- **Connection Reuse**: HTTP client connection pooling for inter-node communication
- **Message Batching**: Efficient message delivery mechanisms
- **Memory Management**: Proper cleanup of connections and resources

## Troubleshooting

### Common Issues

#### Connection Problems
- **MongoDB Connection**: Verify MongoDB connectivity and credentials
- **Port Binding**: Ensure HTTP server port is available
- **WebSocket Upgrade**: Check for proxy or firewall interference

#### Message Delivery Issues
- **Node Discovery**: Verify nodes can communicate with each other
- **Authentication**: Check inter-node authentication credentials
- **Database Queries**: Monitor MongoDB query performance

#### Performance Issues
- **Connection Limits**: Check if bridge limits are reached
- **Database Performance**: Monitor MongoDB server resources
- **Network Latency**: Verify network connectivity between nodes

### Debug Mode
- Set `logger.level: debug` for detailed logging
- Monitor HTTP request/response cycles
- Check database query execution times
- Verify WebSocket message flow

## Contributing

### Development Setup
1. Fork the repository
2. Set up local development environment
3. Run tests: `go test ./...`
4. Follow Go coding standards
5. Submit pull request with tests

### Code Structure
- `src/core/`: Core business logic and interfaces
- `src/handlers/`: HTTP request handlers
- `src/impl/`: Implementation of core interfaces
- `src/middlewares/`: HTTP middlewares
- `src/mongodb/`: Database operations
- `src/utils/`: Utility functions and helpers

This documentation provides a comprehensive overview of Rosenbridge's architecture, deployment, and usage. For specific implementation details, refer to the source code and inline documentation.