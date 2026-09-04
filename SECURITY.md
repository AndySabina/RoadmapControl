# Security Policy

## Supported versions

RoadmapControl is in a pre-1.0 bootstrap stage and has no released or operational version. No current commit should be treated as a supported security product. Once releases begin, this section will identify supported release lines; until then, security fixes may be applied only to the default branch and are not subject to backward-compatibility guarantees.

## Report a vulnerability privately

Do not open a public issue for a suspected vulnerability or include exploit details, credentials, tokens, private repository data, or personal information in public discussions.

Use this repository's **Security** tab to submit a private vulnerability report or open a private draft GitHub Security Advisory. If private reporting is unavailable, contact GitHub Support for help reaching the repository's private security-reporting facilities rather than sending details to a personal email address.

Include only the evidence needed to investigate:

- affected revision, planned component, or document;
- impact and realistic attack conditions;
- minimal reproduction steps or proof of concept;
- whether the issue is already being exploited or publicly known; and
- suggested mitigation, if known.

Sanitize logs and attachments. Replace tokens, secrets, usernames, repository names, hostnames, and private URLs with safe placeholders.

## What to expect

Maintainers will acknowledge reports when they are able, assess scope and severity, and coordinate next steps through the private report. Response and remediation times depend on maintainer availability, complexity, and release readiness; this project does not promise a fixed SLA. Please keep the report private until maintainers agree that coordinated disclosure is safe.

Submitting a report does not authorize testing against systems or repositories you do not own or have explicit permission to test.

## Secrets and credentials

- Never commit credentials, tokens, private keys, `.env` contents, or private repository data.
- Revoke and rotate an exposed secret immediately; removing it from the latest commit is not sufficient.
- Use GitHub Actions secrets or an approved credential manager for required credentials.
- Share sensitive evidence only through GitHub's private security-reporting facilities.

## Non-security support

Use a bug report for reproducible defects and a feature request for proposed behavior. Questions about setup, roadmap approval, contribution authorization, or general support belong in the repository's public issue/discussion facilities when available. Security reporting is not a route to priority, roadmap approval, or private general support.
