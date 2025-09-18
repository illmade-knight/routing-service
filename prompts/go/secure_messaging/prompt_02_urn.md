# **Prompt 02: URN Data Model**

**Objective**: Generate the urn/urn.go file. This file will define the strongly-typed URN identifier, which is a foundational data type for the entire secure messaging ecosystem.

**TDD Phase**: N/A (Foundational Model)

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the urn/urn.go file.

The file must define a URN struct that represents a validated urn:sm:{type}:{id} identifier. The implementation must include Parse, String, and IsZero methods, along with custom JSON marshaling/unmarshaling logic. The JSON unmarshaler **must** be backward-compatible, correctly parsing both full URN strings and legacy string user IDs into a user-typed URN.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer focused on creating robust and type-safe data models.

**Overall Goal**: We are building the foundational data models for our routing service. The first and most critical model is a validated Uniform Resource Name (URN) that will be used to identify all entities in the system securely and unambiguously.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **urn/urn.go**.

**Key Requirements to Fulfill**:

* **L3-1.2.1 (Uniform Resource Name)**: The implementation must strictly conform to the urn:sm:{entityType}:{entityID} format.
* **Encapsulation**: The fields of the URN struct **must** be unexported to ensure that all instances are created via a validating constructor, preventing the creation of invalid URNs.
* **Validation**:
    * The New constructor **must** validate that its input components (namespace, entityType, entityID) are not empty.
    * The Parse function **must** validate that the input string has the correct number of parts and the correct urn scheme.
* **Backward Compatibility**: The UnmarshalJSON method **must** correctly handle two cases:
    1. A valid, modern URN string (e.g., "urn:sm:user:user-bob").
    2. A legacy user ID string (e.g., "user-bob"), which it should parse into a URN with the default namespace (sm) and entity type (user).

**Dependencies**: This is the first file. It has no other code dependencies.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created from the existing routing-service source code.

---
