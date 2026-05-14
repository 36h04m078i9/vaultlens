package snapshot_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultlens/internal/snapshot"
)

func makePolicyApplyStore(t *testing.T) *snapshot.Store {
	t.Helper()
	s, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestApplyPolicyNilPolicyStoreReturnsError(t *testing.T) {
	ss := makePolicyApplyStore(t)
	if err := snapshot.ApplyPolicy("x", nil, ss); err == nil {
		t.Fatal("expected error for nil policy store")
	}
}

func TestApplyPolicyNilSnapshotStoreReturnsError(t *testing.T) {
	ps := snapshot.NewPolicyStore()
	if err := snapshot.ApplyPolicy("x", ps, nil); err == nil {
		t.Fatal("expected error for nil snapshot store")
	}
}

func TestApplyPolicyUnknownNameReturnsError(t *testing.T) {
	ps := snapshot.NewPolicyStore()
	ss := makePolicyApplyStore(t)
	if err := snapshot.ApplyPolicy("missing", ps, ss); err == nil {
		t.Fatal("expected error for unknown policy name")
	}
}

func TestApplyPolicyPrunesSnapshots(t *testing.T) {
	ps := snapshot.NewPolicyStore()
	_ = ps.Add("keep1", snapshot.RetentionPolicy{KeepLast: 1})

	ss := makePolicyApplyStore(t)
	now := time.Now()
	for i, id := range []string{"snap-a", "snap-b", "snap-c"} {
		s := &snapshot.Snapshot{ID: id, CapturedAt: now.Add(time.Duration(i) * time.Minute)}
		if err := ss.Save(s); err != nil {
			t.Fatalf("Save %s: %v", id, err)
		}
	}

	if err := snapshot.ApplyPolicy("keep1", ps, ss); err != nil {
		t.Fatalf("ApplyPolicy: %v", err)
	}

	snaps, err := ss.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(snaps) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(snaps))
	}
}

func TestApplyAllPoliciesNilPolicyStoreReturnsError(t *testing.T) {
	ss := makePolicyApplyStore(t)
	if err := snapshot.ApplyAllPolicies(nil, ss); err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyAllPoliciesRunsEachPolicy(t *testing.T) {
	ps := snapshot.NewPolicyStore()
	_ = ps.Add("p1", snapshot.RetentionPolicy{KeepLast: 10})
	_ = ps.Add("p2", snapshot.RetentionPolicy{OlderThan: 24 * time.Hour})
	ss := makePolicyApplyStore(t)
	if err := snapshot.ApplyAllPolicies(ps, ss); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
