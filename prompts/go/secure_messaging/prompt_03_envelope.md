### **prompt\_03\_envelope\_model.md**

# **Prompt 03: Secure Envelope Data Model**

**Objective**: Generate the transport/envelope.go file. This defines the SecureEnvelope struct, the primary data transfer object that the routing service handles.

**TDD Phase**: N/A (Foundational Model)

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the transport/envelope.go file.

This file must define the SecureEnvelope struct. The struct must include fields for SenderID, RecipientID, MessageID, EncryptedData, EncryptedSymmetricKey, and Signature. All fields must have the correct json:"..." tags for serialization.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer creating the core data transfer objects for a secure messaging system.

**Overall Goal**: We are continuing to build our foundational data models. After defining the URN for identification, we now need to define the "container" that will carry the actual message data through the system.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **transport/envelope.go**.

**Key Requirements to Fulfill**:

* **L3-1.2.2 (Secure Envelope)**: The generated struct must perfectly match the specified schema.
* **Content Agnosticism (L1-1.1)**: The fields related to the payload (EncryptedData, etc.) are defined as opaque byte slices (\[\]byte), reinforcing the requirement that the routing service does not need to understand the content it is routing.
* **Type Safety**: The SenderID and RecipientID fields **must** use the urn.URN type.

**Dependencies**:

* You should assume that the urn package and the urn.URN type from prompt\_02\_urn\_model.md are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created from the existing routing-service source code.

---

