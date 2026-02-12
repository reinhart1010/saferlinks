# Enforcement Policies

SaferLinks uses a hierarchical set of Enforcement Policies:

## Which tag is appropriate?

First, look at how a certain policy (e.g., requires to view on a mobile app) matches your observations in the table:

| Policy Tag | Effect Level* | User can opt-out from enforcement | Content intended to be public | Content intended to be private | Site has anti-bot protections | Site has anti-NSFW content filter |
|---|---|---|---|---|---|---|
| `not_enforced` | 0 | ✅ | ❌ | ❌ | ❌ | ❌ |
| `optionally_enforced` | 1 | ✅ | Any | Any | Maybe | Maybe |
| `enforced_if_mass_traffic` | 2 | ❌ | ❌ | Maybe | ✅ | Maybe |
| `enforced_if_private_intended` | 3 | ❌ | ❌ | ✅ | ✅ (Assumed) | Maybe |
| `enforced_if_auto_nsfw` | 4 | ❌ | ✅ (if assumed NSFW by platform) | ✅ | ✅ (Assumed) | ✅ |
| `enforced_if_interacting` | 5 | ❌ | ✅ (e.g., to like or comment) | ✅ | ✅ (Assumed) | Maybe |
| `enforced_if_consuming` | 6 | ❌ | ✅ (e.g., to read more the article, or view video for longer than 30 seconds) | ✅ | ✅ (Assumed) | Maybe |
| `always_enforced` | 7 | ✅ | ✅ | ✅ | ✅ (Assumed) | Maybe |

*DO NOT use the raw Effect Level values as the policy rankings may change in the future. If you are writing YAML policies, use the tag names in `snake_case` convention. If you are working with Go code, use the official enum name in `MixedCaps` convention.