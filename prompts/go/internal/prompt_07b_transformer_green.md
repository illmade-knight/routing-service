# **Prompt 07b: EnvelopeTransformer Implementation (Green)**

**Objective**: Generate the complete implementation for the EnvelopeTransformer function. The generated code must satisfy the previously defined unit tests.

**TDD Phase**: Green (Write the code to make the failing tests pass).

## **Primary Prompt (Concise)**

Your task is to generate the Go code for the internal/pipeline/transformer\_routing.go file.

I have provided a complete unit test file (transformer\_routing\_test.go) that defines the function's behavior. Your implementation of the EnvelopeTransformer function **must** make all the tests in the provided test file pass.

## **Reinforcement Prompt (Detailed Context)**

**Persona**: You are an expert Go developer writing a data processing pipeline and are focused on security and robustness, ensuring no invalid data can proceed to later stages.

**Overall Goal**: We have an executable specification in the form of a failing unit test (transformer\_routing\_test.go). Your task is now to write the Go code for the EnvelopeTransformer function that satisfies that specification.

**Your Specific Task**: Your task is to generate the complete code for a single Go file: **internal/pipeline/transformer\_routing.go**.

**Key Requirements to Fulfill**:

* **Pass All Tests**: This is the most critical requirement.
* **Unmarshaling (L2-2.2)**: The function **must** first attempt to json.Unmarshal the msg.Payload into a transport.SecureEnvelope. If this fails, it must immediately return with skip=true and an error.
* **Strict Validation (L2-2.2)**: After a successful unmarshal, the function **must** immediately perform a validation check to ensure the envelope is not malformed. It **must** check if envelope.SenderID.IsZero() or envelope.RecipientID.IsZero(). If either field is missing, the function must return with skip=true and an error. This is a critical guardrail for the pipeline.
* **Correct Return Values**: The function must return the exact tuple of (\*transport.SecureEnvelope, bool, error) that satisfies each case in the test specification.

**Dependencies**:

1. **The Specification**: internal/pipeline/transformer\_routing\_test.go
2. **The Data Models**: pkg/transport/envelope.go and a generic messagepipeline package.

---

## **Prompt Changelog**

### **RV\_20250918**

* Initial prompt version created. Mandates strict post-unmarshal validation to fulfill the input validation requirement (L2-2.2).