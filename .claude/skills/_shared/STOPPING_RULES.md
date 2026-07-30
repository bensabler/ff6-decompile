# Stopping Rules

Stop the active unit when any condition occurs:

- the question is answered to its planned confidence threshold;
- the next step requires a new experiment;
- a contradiction blocks interpretation;
- required state cannot be reproduced;
- the task exceeds its declared scope;
- the remaining work is implementation and should move to `/implement-discovery`;
- context or time requires a checkpoint.

Do not expand scope silently.
