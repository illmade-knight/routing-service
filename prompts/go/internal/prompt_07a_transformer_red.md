# **Prompt 07a: EnvelopeTransformer Unit Tests (Red)**

**Objective**: Generate a complete **unit test** file (transformer\_routing\_test.go) for the EnvelopeTransformer function. This test will serve as the executable specification for the pipeline's transformation and validation stage.

**TDD Phase**: Red (Generate failing tests before the code exists).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/pipeline/transformer\_routing\_test.go file.

This file must contain a **unit test** for the EnvelopeTransformer function, which does not yet exist. The test **must be table-driven** and validate all logic paths, including successful transformation, JSON unmarshaling errors, and validation failures for missing sender or recipient IDs.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer with a strong preference for writing clean, comprehensive, table-driven unit tests.

**Overall Goal**: We are implementing the message processing pipeline. The first stage is a "transformer" responsible for safely parsing and validating raw incoming message payloads. We will define its behavior with a precise unit test before writing the implementation.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/pipeline/transformer\_routing\_test.go**.

**Key Requirements to Test**:

* **Table-Driven Structure**: The test **must** use a table-driven approach to define and execute all test cases.
* **Test Cases**: You **must** create test cases to validate the following scenarios for the EnvelopeTransformer(ctx, msg) function:
    1. **Success Case**: The input messagepipeline.Message contains a valid JSON payload for a transport.SecureEnvelope. The test **must assert** that the function returns the correctly deserialized envelope, a skip value of false, and a nil error.
    2. **Invalid JSON Failure**: The input payload is malformed and not valid JSON. The test **must assert** that the function returns a nil envelope, a skip value of true, and a non-nil error.
    3. **Validation Failure (Missing Recipient)**: The payload is valid JSON but is missing the recipientId field. The test **must assert** that the function returns a nil envelope, a skip value of true, and a non-nil error.
    4. **Validation Failure (Missing Sender)**: The payload is valid JSON but is missing the senderId field. The test **must assert** that the function returns a nil envelope, a skip value of true, and a non-nil error.

**Dependencies**:

* You should assume the pkg/transport, pkg/urn, and a generic messagepipeline package are available for import.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Focuses on a table-driven unit test to specify the behavior of this pure function.

---
