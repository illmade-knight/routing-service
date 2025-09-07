# **Routing Service**

The Routing Service is a high-performance, secure, and scalable microservice responsible for the real-time delivery of SecureEnvelope messages between clients. It acts as a "dumb pipe," routing messages based on their unencrypted headers with zero knowledge of the encrypted payload, ensuring end-to-end security.

The architecture is built around a **Unified Ingestion Pipeline** using the go-dataflow library. This allows multiple transport protocols (e.g., HTTP, MQTT) to produce messages to a central message bus, which are then consumed by a single, scalable backend that performs the core routing logic.

---

## **Directory Structure**

The repository is organized using standard Go project layout conventions to maintain a clear separation between public APIs, internal logic, and the final executable.

Plaintext

.  
├── cmd/  
│   └── routing-service/  
│       └── main.go              \# Assembles and runs the service  
├── e2e/  
│   └── routing\_e2e\_test.go      \# End-to-end integration tests  
├── internal/  
│   ├── api/                     \# Private HTTP API handlers  
│   ├── pipeline/                \# Core message processing stages (transform, process)  
│   └── platform/                \# Concrete adapters for external services (e.g., Pub/Sub)  
├── pkg/  
│   └── routing/                 \# Public "contract" (domain models, interfaces, config)  
├── routingservice/  
│   └── service.go               \# The public, embeddable service library/wrapper  
└── test/  
└── e2e\_helpers.go           \# Public test helpers for E2E tests

* **cmd/**: Contains the executable entry point for the service. Its main.go file is responsible for reading configuration, creating concrete dependencies (like a real Google Pub/Sub client), and injecting them into the service wrapper to run the application.
* **e2e/**: Holds the end-to-end integration tests. These tests treat the service as a black box, interacting with its public APIs (like the HTTP endpoint) to validate the entire workflow.
* **internal/**: Contains all the private application logic that cannot be imported by other projects. This includes the HTTP handlers (api), the core pipeline logic (pipeline), and the concrete adapters for external services (platform).
* **pkg/**: Holds public, shareable libraries. pkg/routing is the service's public contract, defining the data structures, interfaces, and configuration that other services or the application's own entry point can use.
* **routingservice/**: The primary, public-facing service library. It provides the Wrapper and its constructor, New(), which assembles the internal components into a runnable service. It is designed to be imported by executables.
* **test/**: Provides public helper functions specifically for testing. This allows our e2e tests to create concrete service dependencies without violating the internal package boundary.

---

## **Current State vs. Roadmap**

The refactoring has successfully established the core architecture, which allows for direct integration with the pre-built components from the go-dataflow library. A significant portion of **Phase 1** is either complete or ready for immediate integration.

### **What's Complete (Phase 1 Foundation)**

The current codebase provides a robust and fully decoupled foundation for the Unified Ingestion Pipeline.

* ✅ **HTTP Ingestion:** The POST /send endpoint is implemented and publishes messages to the central message bus.
* ✅ **Message Processing Pipeline:** A fully functional StreamingService is in place. It is designed to use a GooglePubsubConsumer to subscribe to the ingress topic, a Transformer to parse the SecureEnvelope, and a Processor to handle the routing logic.
* ✅ **Core Routing Logic:** The Processor contains the complete logic for checking user presence and deciding between online and offline delivery paths. It depends on generic Fetcher interfaces for its data lookups.

### **Ready for Integration (Phase 1 Data Adapters)**

The "missing" data adapters from the roadmap do not need to be re-implemented. The next step is to instantiate and inject the existing, generic implementations from the go-dataflow library within cmd/routing-service/main.go.

* ➡️ **Redis for Presence:** The PresenceCache dependency can be satisfied by instantiating cache.NewRedisCache\[string, routing.ConnectionInfo\] from go-dataflow. This cache should be configured with no fallback, as a miss in Redis definitively means the user is offline.
* ➡️ **Firestore for Device Tokens:** The DeviceTokenFetcher dependency can be satisfied by instantiating cache.NewFirestore\[string, \[\]routing.DeviceToken\] from go-dataflow. For improved performance, this Firestore fetcher can be used as the fallback for another caching layer, such as RedisCache or InMemoryLRUCache.

### **What's Next (Future Roadmap Items)**

The following major components from the roadmap have not yet been started and represent the next phases of development:

* ❌ **MQTT Ingestion:** The MQTT ingestion point needs to be built and configured to produce to the same ingress topic.
* ❌ **Phase 2 \- Real-time Delivery:** The entire real-time connection manager, including the **WebSocket handler** and **MQTT publisher** for delivering messages to online clients, needs to be implemented.
* ❌ **Phase 2 \- Additional APIs:** The GET /connect (WebSocket) and POST /register-device endpoints are not yet created.
* ❌ **Phase 3/4 \- Production Hardening:** Features like JWT authentication, observability, and deployment configurations are still pending.