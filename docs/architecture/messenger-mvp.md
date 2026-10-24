# Messenger MVP Flow

This historical sequence describes the minimum contact-request flow between two
messenger clients. It intentionally omits repeated account-log operations so the
application-visible behavior remains readable.

## Simplifying assumption

In the actual protocol, the application subscribes to its account log and the
protocol appends events to that log. The diagram draws some protocol-to-app
arrows directly instead of showing every append and subscription delivery.

## Contact request

```mermaid
sequenceDiagram
	participant AO as Alice store
	participant AP as Alice protocol
	participant AC as Alice app
	actor Alice
	actor Bob
	participant BC as Bob app
	participant BP as Bob protocol
	participant BO as Bob store

	Alice->>AC: Start app
	AC->>AP: Enable contact requests
	AP-->>AO: Subscribe to rendezvous log
	AP->>AC: Return shareable reference
	AC->>Alice: Display reference or QR code
	Alice->>Bob: Share reference and metadata
	Bob->>BC: Import reference
	BC->>BP: Subscribe to one-to-one metadata
	BP-->>BO: Subscribe to metadata log
	BC->>BP: Send contact request
	BP-->>BC: Report request queued
	BP->>BO: Publish request on Alice rendezvous log
	BO->>AO: Replicate rendezvous log
	AO->>AP: Deliver contact request
	AP->>AC: Report incoming request
	AC->>Alice: Display request
	Alice->>AC: Accept request
	AC->>AP: Accept contact request
	AP->>AO: Append member-device event
	AO->>BO: Replicate group metadata
	BO->>BP: Deliver member-device event
	BP->>BC: Report accepted request
	BC->>Bob: Open conversation
```

## Security properties

- The shareable reference bootstraps discovery; it is not a substitute for
  authenticated group membership.
- Contact and membership events must be validated before changing local state.
- Applications should make acceptance explicit and identify the requesting
  contact clearly.
- Shared references should be revocable or replaceable when exposure is
  suspected.

For current message and service types, consult the
[generated protocol reference](../apis/protocoltypes.md).
