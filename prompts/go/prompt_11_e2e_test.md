### **prompt\_11\_e2e\_test.md**

# **Prompt 11: End-to-End Service Test**

**Objective**: Generate a final, comprehensive **end-to-end test** (routing\_e2e\_test.go). This test will validate the entire service, running it as a complete application with live emulators to simulate its cloud environment.

**TDD Phase**: N/A (Final Validation)

## **Primary Prompt (Concise)**

Your task is to generate the Go code for e2e/routing\_e2e\_test.go.

This file must contain a comprehensive **end-to-end test** that validates the entire service flow. The test must start the complete routing service, using **live emulators for Pub/Sub and Firestore**, and interact with it via its public HTTP API to verify the full message store-and-retrieve lifecycle for an offline user.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer specializing in complex, multi-component integration and end-to-end testing for cloud-native microservices.

**Overall Goal**: We have now unit-tested all the individual components of the routing-service in isolation. The final step is to create a definitive test that proves all these components work together correctly as a cohesive system, interacting with a production-like environment.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **e2e/routing\_e2e\_test.go**.

**Key Requirements for the Test Flow**:

* **Emulator Setup**: The test **must** programmatically set up and connect to both a **Google Cloud Pub/Sub emulator** and a **Google Cloud Firestore emulator**.  
* **Service Startup**: It must assemble the full routingservice with all its concrete, production-ready dependencies (wired up to the emulator clients) and start it.  
* **Test Scenario**: The test must execute the following two-phase scenario:  
  1. **Phase 1: Send Message to Offline User & Verify Storage**  
     * Use a standard Go HTTP client to POST a SecureEnvelope to the running service's /send endpoint.  
     * Use a mock PushNotifier to confirm that the message was processed via the offline path.  
     * Use require.Eventually to poll the Firestore emulator's database directly and **assert** that the message has been correctly persisted in the expected collection (user-messages/{urn}/messages/{msgId}).  
  2. **Phase 2: Retrieve Message & Verify Cleanup**  
     * Use the HTTP client to perform a GET request to the /messages endpoint, providing the correct user URN in the X-User-ID header.  
     * **Assert** that the HTTP response is 200 OK and that its JSON body contains the exact message that was originally sent.  
     * Use require.Eventually to poll the Firestore emulator again and **assert** that the message document has been deleted, confirming the async cleanup logic worked.

**Dependencies**:

* This test will depend on the entire routing-service codebase being available.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Defines the final, high-fidelity end-to-end test to validate the full service behavior against live emulators.