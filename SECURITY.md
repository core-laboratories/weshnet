# Security Policy

## Supported versions

Security fixes are applied to the latest maintained release and the current
default branch. Older releases may require upgrading before a fix can be used.

## Reporting a vulnerability

Do not open a public issue for a suspected vulnerability. Use GitHub's private
security advisory flow for the
[Weshnet repository](https://github.com/core-laboratories/weshnet/security/advisories/new).

Include the following information when available:

- affected version or commit;
- affected platform and configuration;
- impact and attack prerequisites;
- reproducible steps or a minimal proof of concept;
- suggested remediation;
- whether the report may be shared with additional maintainers.

Do not include real user secrets, private keys, recovery material, message
content, or production credentials in a report.

## Response process

Maintainers will acknowledge a complete report, investigate its impact, and
coordinate remediation and disclosure. Response and release timing depends on
severity, reproducibility, and the affected dependency chain.

Please allow maintainers a reasonable opportunity to address the issue before
publishing details. Good-faith research that avoids privacy violations, data
destruction, service disruption, and unauthorized access is appreciated.

## Scope

Reports are especially useful when they concern:

- cryptographic key handling;
- authorization or group membership;
- message confidentiality or authenticity;
- rendezvous and discovery privacy;
- malformed protocol input;
- local datastore protection;
- dependency vulnerabilities with a demonstrated Weshnet impact.

General support requests and non-security bugs belong in the public
[issue tracker](https://github.com/core-laboratories/weshnet/issues).
