### **L3: Technical & Configuration Requirements (L3\_Technical\_Requirements.md)**

This document defines the *specific* technologies, schemas, and technical details required for this implementation.

# **L3: Technical & Configuration Requirements**

## **1\. Configuration Schema**

### **1.1. Application Configuration**

The system **shall** be configured via environment variables that map to the following fields:

| Field | Env Variable | Required | Description |
| :---- | :---- | :---- | :---- |
| ProjectID | GCP\_PROJECT\_ID | Yes | The target Google Cloud project ID. |
| HTTPListenAddr | N/A (hardcoded) | N/A | The host and port for the HTTP API server (e.g., :8082). |
| IngressSubscriptionID | N/A (hardcoded) | N/A | The name of the Pub/Sub subscription for incoming messages. |
| IngressTopicID | N/A (hardcoded) | N/A | The name of the Pub/Sub topic for message ingestion. |
| NumPipelineWorkers | N/A (hardcoded) | N/A | The number of concurrent workers for the message pipeline. |

### **1.2. Data Schemas**

The system **shall** use the following data structures for its core operations.

#### **1.2.1. Uniform Resource Name (URN)**

A validated identifier for all entities, conforming to the urn:sm:{entityType}:{entityID} format. It **must** support JSON serialization to and from a string. It **must** also support backward-compatible deserialization from a plain user ID string.

#### **1.2.2. Secure Envelope**

The primary data transfer object, which the service routes but does not decrypt.

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| senderId | urn.URN | Yes | The URN of the message sender. |
| recipientId | urn.URN | Yes | The URN of the message recipient. |
| messageId | String | No | A unique identifier for the message. |
| encryptedData | \[\]byte | Yes | The opaque, encrypted message payload. |
| encryptedSymmetricKey | \[\]byte | Yes | The encrypted symmetric key for the payload. |
| signature | \[\]byte | Yes | The signature of the encrypted data. |

## **2\. Supported Cloud Services (GCP)**

| Category | Service | Implementation Detail |
| :---- | :---- | :---- |
| **Messaging** | Google Cloud Pub/Sub | Used for the primary message ingestion bus. |
| **Database** | Google Cloud Firestore | Used as the persistence layer (MessageStore) for offline messages. |
| **Storage Schema** | Firestore Collections | Offline messages **shall** be stored in a sub-collection structure: user-messages/{recipient-urn}/messages/{message-id}. |

