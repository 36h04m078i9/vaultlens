// Package snapshot provides facilities for capturing, storing, comparing,
// and managing point-in-time snapshots of Vault secret paths.
//
// # Retention Policies
//
// A RetentionPolicy encapsulates rules that determine how many snapshots to
// keep and for how long. Policies are stored by name in a PolicyStore and
// can be applied individually via ApplyPolicy or in bulk via ApplyAllPolicies.
//
// Example:
//
//	ps := snapshot.NewPolicyStore()
//	_ = ps.Add("daily", snapshot.RetentionPolicy{
//		KeepLast:  30,
//		OlderThan: 30 * 24 * time.Hour,
//	})
//
//	ss, _ := snapshot.NewStore("/var/vaultlens/snapshots")
//	_ = snapshot.ApplyPolicy("daily", ps, ss)
//
// Policies may also be scoped to a path prefix so that different parts of
// the Vault hierarchy can have independent retention rules.
package snapshot
