# ezMsg Project TODO List

## Phase 1: Core Messaging MVP

**Goal:** Get a basic, real-time chat working between two clients. For this phase, authentication is mocked/hardcoded and there is no E2EE.

#### Server-Side
- [ ] **Establish Basic Server:** Use your existing QUIC implementation to listen for and accept client connections.
- [ ] **Manage Connections:** Implement a simple in-memory map to keep track of active client connections to know who to send messages to.
- [ ] **Define Simple Message Format:** Decide on a simple text-based format for messages, like JSON (`{"sender": "user1", "text": "Hello!"}`), ignoring FlatBuffers and E2EE for now.
- [ ] **Implement Message Relay Logic:** In your `handleStream` function, when a message is received from one client, loop through the connection map and broadcast it to all other connected clients.

#### Client-Side
- [ ] **Build Core Chat UI:** Use your existing Gio UI code to create a message display area, an input widget, and a send button.
- [ ] **Connect to Server:** Write the client-side logic to establish a QUIC connection to the server when the app starts.
- [ ] **Send Messages:** When the user clicks the "Send" button, format the text into the simple message format and send it over the QUIC stream.
- [ ] **Receive and Display Messages:** Create a `goroutine` that constantly listens for incoming data from the server. When a new message is received, add it to the chat UI's list and trigger a refresh.

---

## Phase 2: Implementing Security - Authentication & E2EE

**Goal:** Replace the MVP components with the full, secure architecture.

#### Server-Side
- [ ] **Implement Decoupled DB Schema:** Create the `Users`, `Devices`, `Messages`, and `MessageKeys` tables.
- [ ] **Implement Challenge-Response:** Build the endpoints for issuing and verifying challenges to grant short-lived `accessTokens`.
- [ ] **Implement Message Ingestion:** Create the endpoint that accepts the single, atomic E2EE package from the client and writes it to the decoupled database tables in a single transaction.
- [ ] **Implement Efficient Delivery:** When a message is ingested, create and push a tailored, secure payload to each recipient device.

#### Client-Side
- [ ] **Implement Key Generation & Storage:** Add the logic to generate key pairs and to encrypt/decrypt the private key file using a local password and **Argon2id**.
- [ ] **Implement E2EE Message Sending:** Replace the simple text sending with the full E2EE flow: encrypt the message, encrypt the `MessageKey` for all recipients, and assemble the single, atomic package.
- [ ] **Implement E2EE Message Receiving:** Replace the simple text receiving with the logic to decrypt the tailored payload from the server.
- [ ] **Implement Local Encrypted Database:** Integrate a library like SQLCipher to store the decrypted message history securely at rest.

---

## Phase 3: Implementing Account Management & Recovery

**Goal:** Build all features related to managing a user's account, devices, and recovery.

#### Server-Side
- [ ] **Implement Event Broadcast System:** Create the mechanism for the server to broadcast system events (`device_added`, `key_changed`, etc.).
- [ ] **Implement Device Linking Endpoints:** Build the logic to verify a signed approval message to add a new public key to an account.
- [ ] **Implement Device Removal Endpoint:** Build the logic to handle a signed command to remove a device, including deleting its `MessageKeys`.
- [ ] **Implement Account Recovery Endpoints:** Build the full, multi-step recovery flow (TOTP check, key rotation, storing the history checkpoint, deleting old keys).

#### Client-Side
- [ ] **Build Registration & Recovery UI:** Create the multi-step UI for the mandatory recovery phrase backup and verification process.
- [ ] **Build Device Linking UI:** Create the UI for both displaying a QR code/text phrase and for scanning/approving a new device.
- [ ] **Build Account Recovery UI:** Create the UI for the entire recovery flow.
- [ ] **Implement System Event Handler:** Write the logic to handle incoming server events (updating local keys, displaying UI security notifications).
- [ ] **Implement History Migration/Sync:** Write the logic for both the one-time history migration after recovery and the two-phase history sync for new devices.

---

## Phase 4: Production Hardening

**Goal:** Add the final layers of security and reliability.

- [ ] Implement server-side rate limiting on all public endpoints.
- [ ] Implement expiration and single-use logic for all authentication challenges.
- [ ] Define and implement a permission model for who can add/remove users from chats.
- [ ] `(Future)` Plan and implement a Forward Secrecy protocol (like Double Ratchet).
- [ ] `(Future)` Research and implement metadata protection techniques.