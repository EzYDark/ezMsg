# Security Implementation Checklist

### **Server-Side Quests (The Gatekeeper)**
*Your server is the bouncer. It checks IDs and permissions before letting anyone in.*

-   **Authentication & Session:**
    -   [ ] **Validate Session Token:** Check the token on every single request.
    -   [ ] **Hash Passwords:** Use **Argon2id** or **scrypt**. Never store plaintext passwords.
    -   [ ] **Secure Session Tokens:** Generate cryptographically random tokens on login.
-   **Authorization (Permissions):**
    -   [ ] **Check Permissions for Every Action:** Before doing anything, check if the user is allowed.
        -   *Is user X in chat Y?*
        -   *Can user A access file B?*
-   **Anti-Replay / DoS:**
    -   [ ] **Time-Window Check:** Reject requests older than your window (e.g., 5 mins).
    -   [ ] **Nonce Check:** If the request is fresh, check the nonce against a persistent store like Redis.
    -   [ ] **Use Redis for Nonce Store:** This is key. It survives restarts and has automatic expiration (TTL).
    -   [ ] **Rate Limiting:** Stop users from spamming requests.
-   **Data Integrity:**
    -   [ ] **Measure File Size:** Ignore the client's claimed size. Measure the encrypted blob yourself and save *that* size.
    -   [ ] **Generate Server Timestamps:** Create the `server_timestamp` when a message is received. This is the truth for message ordering.

---

### **Client-Side Quests (The Final Checkpoint)**
*Your client receives a sealed box. It must check the contents carefully before opening.*

-   **Payload Validation (Most Critical):**
    -   [ ] **Use Flatbuffers Verifier:** **Always** run the `Verify...` function on data from the server *before* parsing it. This prevents crashes from bad data.
-   **Content Sanitization:**
    -   [ ] **Limit String Lengths:** Don't let a huge message crash the app.
    -   [ ] **Validate Logic:** Does a `reply_to_message_uid` actually exist? If not, don't render it as a reply.
    -   [ ] **Prevent Duplicates:** Check the `message_uid` on pushed messages to avoid showing the same message twice.
-   **Trust Management:**
    -   [ ] **Trust Server Metadata:** Always use `server_timestamp` for sorting and `authoritative_size` for files.
-   **Local Storage Security:**
    -   [ ] **Encrypt Local Data:** If you save keys or messages on the device, encrypt them using the OS keychain/keystore.

---

### **E2EE & Crypto Quests**
*The core of your privacy promise.*

-   [ ] **Use Standard Libraries:** Don't invent your own crypto.
-   [ ] **Unique IV/Nonce for Every Encryption:** This is non-negotiable. One key, one IV, one time.
-   [ ] **(Advanced) Key Verification:** For top-tier security, add a way for users to confirm they're talking to the right person (e.g., QR code scan).