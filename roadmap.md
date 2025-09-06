# **Routing Service: Development Roadmap**

## **1\. Vision & Core Principles**

The **Routing Service** is a high-performance, secure, and scalable microservice responsible for the real-time delivery of SecureEnvelope messages between clients across multiple transport protocols (Web, Mobile, MQTT).

* **Security First:** The service is a "dumb pipe." It operates on the unencrypted header of an envelope but has zero-knowledge of the encrypted payload.
* **High Scalability:** The architecture must support a large number of concurrent users and high message throughput.
* **Efficiency:** Message delivery should have minimal latency for online users.
* **Resilience:** The service must handle offline users gracefully and be resilient to failure.

## **2\. Core Architecture**

The service will be built around a **Unified Ingestion Pipeline**, leveraging the go-dataflow library. Multiple ingestion points (HTTP, MQTT) will produce messages to a central message bus. A single, scalable processing backend will consume from this bus to perform the core routing logic.

* **Data Storage Strategy:** A two-tiered approach for user presence and delivery information.
    * **Redis (Fast Cache/Session Store):** Primary store for real-time connection status of online users (WebSocket, MQTT, etc.).
        * *Key:* presence:\<userID\>
        * *Value:* connectionInfo (struct containing server instance ID and protocol type).
    * **Firestore (Persistent Fallback):** Long-term storage for offline push notification tokens.
        * *Key:* devices:\<userID\>
        * *Value:* A list of device tokens (FCM, APNS).

## **3\. Development Phases**

### **Phase 1: Foundation & Unified Ingestion Pipeline**

*Objective: Build the core message processing pipeline and its multi-protocol ingestion frontends.*

1. **Ingestion Points (Producers):**
    * **HTTP Ingestion (for Web/Mobile):**
        * Create an HTTP handler (POST /send) that accepts a SecureEnvelope.
        * This handler will validate the request, create a messagepipeline.Message, and publish it to a central Google Pub/Sub topic (e.g., envelope-ingress).
        * **Reuse:** googleproducer.go.
    * **MQTT Ingestion (for IoT/Mobile):**
        * Implement a standalone ingestion component using the provided mqttconverter package.
        * This component will use the mqttconverter.MqttConsumer to subscribe to a broker topic (e.g., envelopes/send).
        * The MqttConsumer will automatically transform incoming MQTT messages into messagepipeline.Message objects.
        * These messages will then be published to the **same** central Pub/Sub topic (envelope-ingress) using a googleproducer.
        * **Reuse:** mqttconsumer.go, googleproducer.go.
2. **Message Processing Pipeline (Consumer):**
    * Implement a messagepipeline.StreamingService as the core of the routing logic.
    * **Consumer:** Use a messagepipeline.GooglePubsubConsumer to subscribe to the envelope-ingress topic.
    * **Transformer:** A MessageTransformer will parse the messagepipeline.Message payload back into a transport.SecureEnvelope struct. This transformer can be enhanced with WithPayloadValidation to reject malformed envelopes early.
    * **Processor:** The StreamProcessor will contain the core routing logic.
3. **Routing Logic (The Processor):**
    * For each envelope, perform a **presence lookup** in Redis using the recipientID.
    * **If found (user is online):** Forward the envelope to the specific server instance/protocol handler indicated in the connectionInfo from Redis. This is done by publishing to a targeted Pub/Sub topic (e.g., delivery-instance-123).
    * **If not found (user is offline):** Query Firestore for the user's device tokens and trigger a push notification via FCM/APNS.
4. **Data Adapters:**
    * Implement the caching/storage layers.
    * Create a RedisCache that implements cache.Cache\[string, ConnectionInfo\].
    * Create a FirestoreFetcher that implements cache.Fetcher\[string, \[\]DeviceToken\].
    * **Reuse:** cache.go, firestore.go.

### **Phase 2: Multi-Transport Delivery & API**

*Objective: Implement the real-time connection handling for delivering messages to online clients.*

1. **Real-time Connection Manager:**
    * Implement a manager for persistent client connections. When a client connects, it authenticates (JWT) and registers its userID and connectionInfo (protocol type, server instance ID) in the Redis session store.
    * **WebSocket Handler:** Manages connections from web and mobile clients.
    * **MQTT Publisher:** While ingestion uses the MqttConsumer, delivery requires a way to *publish* back to a specific topic for an MQTT client. The connection manager will handle this, publishing messages to topics like clients/\<userID\>/receive.
    * Each server instance will subscribe to its own targeted Pub/Sub topic to receive routed messages from the pipeline processor.
2. **API Endpoints:**
    * **POST /send**: HTTP ingestion endpoint.
    * **GET /connect**: WebSocket upgrade endpoint.
    * **POST /register-device**: For mobile clients to register push notification tokens.

### **Phase 3 & 4: Production Hardening & Deployment**

*(These phases remain largely the same but are strengthened by the robust pipeline architecture)*

* **Security:** JWT authentication will be required for HTTP, WebSocket, and MQTT connections.
* **Observability:** Add metrics for each ingestion path (HTTP/MQTT) and for the depth of the central Pub/Sub queue.
* **Configuration:** Add configuration for MQTT broker details and topic mappings.
* **Deployment:** The MQTT ingestion component can be deployed as a separate "sidecar" service alongside the main routing service, allowing them to be scaled independently.