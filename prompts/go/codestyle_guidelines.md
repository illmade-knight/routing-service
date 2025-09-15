## **Go Language Coding & Refactoring Guidelines**

### **1\. General Principles**

* **Grounding**: Work **only** from the provided source files. Do not add, assume, or hallucinate any logic not explicitly present.
* **Completeness**: When providing a refactored file, you **must** provide the **entire file content**, from the package declaration to the final line. Never use snippets or placeholders.
* **Precision**: A refactor must **never** change unrelated code. Apply only the requested change.

---

### **2\. Code Style & Formatting**

* **Error Variable Declaration**: Use the explicit a, err := style on its own line, not the if ...; err \!= nil compound statement.
* **Struct Literals**: Always use the expanded, multi-line format for struct literals for readability.
* **Constructor Parameter Order**: Use the preferred order: context, configuration, clients/dependencies, then other parameters.

---

### **3\. Error Handling**

* **Handle All Critical Errors**: All functions that return an error must have the error value checked.
* **Provide Context**: Wrap downstream errors with fmt.Errorf and the %w verb to create a meaningful error stack.
* **Log Non-Critical Errors**: Errors from non-critical operations like io.Closer.Close() on a reader can be logged at a WARN or INFO level instead of being returned.

---

### **4\. Testing**

* **Test Package Naming**: Test files must belong to a \_test package (e.g., package routing\_test).
* **Test Cleanup**: Prefer t.Cleanup() for scheduling cleanup tasks over defer.
* **Top-Level Test Context**: Every complex test should establish a top-level context.Context with a reasonable timeout to prevent hangs.
* **Avoid Sleeps**: Never use time.Sleep() to wait for an operation. Use require.Eventually or a custom polling loop with a timeout.
* **Table-Driven Tests**: Use table-driven tests for testing multiple scenarios of the same function.

---

### **5\. Documentation**

* **User-Facing Comments**: All public types, functions, and methods should have clear, user-focused godoc comments.
* **Refactoring Comments**: Comments related to a refactoring task should be prefixed with // REFACTOR:.

---

### **Highlighted Additions for Modern Go (1.22+)**

Here are additional guidelines to ensure the library follows modern Go best practices.

* **Use Modern range Semantics**: In Go 1.22 and later, the behavior of range loops was changed to prevent common bugs with closure captures in concurrent code. All for...range loops should be written to take advantage of this new, safer behavior.
* **Prefer Structured Logging with slog**: For all new components, use the standard library's log/slog package for structured, leveled logging. This is the new standard in the Go ecosystem. When logging complex objects, create custom slog.LogValuer types to control their representation.
* **Use Enhanced http.ServeMux for Routing**: For HTTP services, use the enhanced http.ServeMux introduced in Go 1.22. It provides built-in support for HTTP methods and path wildcards (e.g., mux.HandleFunc("GET /items/{id}", ...)), which is cleaner than manual method checking.
* **Aggregate Errors with errors.Join**: When a function needs to return multiple errors (e.g., from a group of concurrent goroutines), use the errors.Join function (introduced in Go 1.20) to wrap them into a single error value.