# Development rules

- Read README.md and docs/MIGRATION.md; record exact upstream commit, source files and changes.
- Preserve upstream packet/wire semantics and existing security policy. Do not roll custom crypto or introduce plaintext fallback.
- Do not modify either product's master during core extraction. Pin versioned dependencies; migrate on product integration branches with rollback.
- Product authorization is separate; core must not grant membership by trusting client-provided room or device IDs.
- No secrets, private keys or user addresses. No anonymous relay or unrestricted broadcast.
- Run relevant tests only; report actual results, NOT RUN gates, and failures. Do not claim real network or game compatibility based on unit tests.
- License and dependency/security gates must be resolved before public distribution; do not select project licensing unilaterally.
