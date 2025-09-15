# **L1: Business & Functional Requirements**

## **1\. Core Mission**

| ID | Requirement | Description |
| :---- | :---- | :---- |
| **1.1** | **Content-Agnostic Message Routing** | The service **shall** act as a content-agnostic router for SecureEnvelope messages. It **must** route messages based solely on the SenderID and RecipientID URNs in the envelope's metadata, without inspecting the encrypted payload. |
| **1.2** | **Guaranteed Offline Message Storage** | The service **shall** guarantee the persistence of messages for users who are not online. It **must** store these messages until the recipient successfully retrieves them. |
| **1.3** | **Asynchronous Ingestion** | The service **shall** provide an asynchronous ingestion mechanism. It **must** immediately accept valid messages for processing and return a success status to the client before the message has been fully routed and delivered. |

## **2\. Key Features**

| ID | Feature | Description |
| :---- | :---- | :---- |
| **2.1** | **Message Ingestion API** | The service **shall** expose an HTTP endpoint (POST /send) to accept a SecureEnvelope for routing. The endpoint **must** validate that the envelope is well-formed JSON and contains a valid RecipientID. |
| **2.2** | **Offline Message Retrieval API** | The service **shall** expose an HTTP endpoint (GET /messages) for users to retrieve all their stored offline messages. The user **must** be identified via an X-User-ID header containing their URN. If no messages are stored, the service **shall** return a 204 No Content status. |
| **2.3** | **Automatic Message Cleanup** | Upon successful retrieval of messages by a user from the /messages endpoint, the service **shall** automatically and asynchronously delete those specific messages from the offline store to prevent re-delivery. |
| **2.4** | **Push Notification Trigger** | When a message for an offline user is stored, the service **shall** attempt to trigger a push notification to the recipient's registered devices to inform them of the new message. The success or failure of the push notification **must not** affect the success of the message storage operation. |

---


