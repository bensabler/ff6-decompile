# Error Policy

Return errors rather than panic for invalid external data. Include operation, domain, offset, expected length, and actual length. Use `errors.Is` and `errors.As` for stable classification.
