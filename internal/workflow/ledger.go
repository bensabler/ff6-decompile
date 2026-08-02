package workflow

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// RunIdentitySchemaVersion is the immutable run-identity format.
	RunIdentitySchemaVersion = "1.0"
	// RunEventSchemaVersion is the normalized event format.
	RunEventSchemaVersion = "1.0"
	genesisEventHash      = "0000000000000000000000000000000000000000000000000000000000000000"
)

var trustBases = map[TrustBasis]bool{
	TrustCollectorObserved: true, TrustBackendExitStatus: true,
	TrustOperatorAttested: true, TrustSelfReported: true, TrustUnsupported: true,
}

var (
	runIDPattern = regexp.MustCompile(`^run-[0-9a-f]{32}$`)
	hashPattern  = regexp.MustCompile(`^[0-9a-f]{64}$`)
	headPattern  = regexp.MustCompile(`^[0-9a-f]{40}(?:[0-9a-f]{24})?$`)
)

// RunIdentity is the immutable binding between a human workflow, one generated
// run, its frozen contract, and the Git repository state where it began.
type RunIdentity struct {
	SchemaVersion      string `json:"schema_version"`
	WorkflowID         string `json:"workflow_id"`
	RunID              string `json:"run_id"`
	ContractHash       string `json:"contract_hash"`
	CreatedAt          string `json:"created_at"`
	RepositoryRoot     string `json:"repository_root"`
	RepositoryIdentity string `json:"repository_identity"`
	StartingBranch     string `json:"starting_branch"`
	StartingHead       string `json:"starting_head"`
}

// SourceKind is an event's claimed provenance. Provenance is evaluated
// independently from structural validity and hash-chain integrity.
type SourceKind string

const (
	SourceProviderHook         SourceKind = "provider_hook"
	SourceDeterministicBackend SourceKind = "deterministic_backend"
	SourceOperatorRecord       SourceKind = "operator_record"
	SourceManualImport         SourceKind = "manual_import"
	SourceLegacyTranscript     SourceKind = "legacy_transcript"
	SourceUnknown              SourceKind = "unknown"
)

var sourceKinds = map[SourceKind]bool{
	SourceProviderHook: true, SourceDeterministicBackend: true,
	SourceOperatorRecord: true, SourceManualImport: true,
	SourceLegacyTranscript: true, SourceUnknown: true,
}

// TrustBasis states why a collector believes an event. These values are not
// signatures and do not make an event cryptographically authentic.
type TrustBasis string

const (
	TrustCollectorObserved TrustBasis = "collector_observed"
	TrustBackendExitStatus TrustBasis = "backend_exit_status"
	TrustOperatorAttested  TrustBasis = "operator_attested"
	TrustSelfReported      TrustBasis = "self_reported"
	TrustUnsupported       TrustBasis = "unsupported"
)

// EventKind is provider-neutral. The schema recognizes more event kinds than
// the current eligibility boundary trusts or converts into observations.
type EventKind string

const (
	EventAgentStarted     EventKind = "agent_started"
	EventAgentFinished    EventKind = "agent_finished"
	EventSkillSelected    EventKind = "skill_selected"
	EventToolStarted      EventKind = "tool_started"
	EventToolFinished     EventKind = "tool_finished"
	EventBackendFinished  EventKind = "backend_finished"
	EventOutputObserved   EventKind = "output_observed"
	EventValidationFinish EventKind = "validation_finished"
	EventOperatorApproval EventKind = "operator_approval"
)

var eventKinds = map[EventKind]bool{
	EventAgentStarted: true, EventAgentFinished: true, EventSkillSelected: true,
	EventToolStarted: true, EventToolFinished: true, EventBackendFinished: true,
	EventOutputObserved: true, EventValidationFinish: true,
	EventOperatorApproval: true,
}

// RunEvent is one normalized JSONL ledger record. Optional values are encoded
// explicitly as empty strings or null, never omitted.
//
// The event hash makes mutation and broken ordering detectable. It does not
// prove who emitted the event. Eligibility is a separate decision.
type RunEvent struct {
	SchemaVersion      string     `json:"schema_version"`
	Sequence           uint64     `json:"sequence"`
	EventID            string     `json:"event_id"`
	PreviousHash       string     `json:"previous_hash"`
	EventHash          string     `json:"event_hash"`
	WorkflowID         string     `json:"workflow_id"`
	RunID              string     `json:"run_id"`
	ContractHash       string     `json:"contract_hash"`
	ObservedAt         string     `json:"observed_at"`
	Provider           string     `json:"provider"`
	SourceKind         SourceKind `json:"source_kind"`
	CollectorID        string     `json:"collector_id"`
	TrustBasis         TrustBasis `json:"trust_basis"`
	SessionID          string     `json:"session_id"`
	TurnID             string     `json:"turn_id"`
	RepositoryIdentity string     `json:"repository_identity"`
	CWD                string     `json:"cwd"`
	Branch             string     `json:"branch"`
	Head               string     `json:"head"`
	EventKind          EventKind  `json:"event_kind"`
	Selector           string     `json:"selector"`
	ToolUseID          string     `json:"tool_use_id"`
	ExitStatus         *int       `json:"exit_status"`
	EvidenceRef        string     `json:"evidence_ref"`
}

// LedgerVerification is the result of reading the complete ledger. Events are
// released to reconciliation only when Valid is true.
type LedgerVerification struct {
	Valid        bool       `json:"valid"`
	Events       []RunEvent `json:"events,omitempty"`
	Problems     []string   `json:"problems,omitempty"`
	TailSequence uint64     `json:"tail_sequence"`
	TailHash     string     `json:"tail_hash"`
}

// EventEligibility is the separate provenance/trust decision for a valid
// event. A selector match is only one input to this decision.
type EventEligibility struct {
	Eligible        bool
	ObservationKind ObservationKind
	Reason          string
}

func newRunID() (string, error)   { return randomID("run-") }
func newEventID() (string, error) { return randomID("evt-") }

func randomID(prefix string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate %s id: %w", strings.TrimSuffix(prefix, "-"), err)
	}
	return prefix + hex.EncodeToString(b), nil
}

func (id RunIdentity) validate() []string {
	var problems []string
	if id.SchemaVersion != RunIdentitySchemaVersion {
		problems = append(problems, fmt.Sprintf("unsupported run identity schema_version %q", id.SchemaVersion))
	}
	if !workflowIDPattern.MatchString(id.WorkflowID) {
		problems = append(problems, fmt.Sprintf("invalid workflow_id %q", id.WorkflowID))
	}
	if !runIDPattern.MatchString(id.RunID) {
		problems = append(problems, fmt.Sprintf("invalid run_id %q", id.RunID))
	}
	if !hashPattern.MatchString(id.ContractHash) {
		problems = append(problems, "invalid contract_hash")
	}
	if _, err := time.Parse(time.RFC3339Nano, id.CreatedAt); err != nil {
		problems = append(problems, "invalid created_at: "+err.Error())
	}
	if !filepath.IsAbs(id.RepositoryRoot) {
		problems = append(problems, "repository_root is not absolute")
	}
	if strings.TrimSpace(id.RepositoryIdentity) == "" {
		problems = append(problems, "repository_identity is empty")
	}
	if strings.TrimSpace(id.StartingBranch) == "" {
		problems = append(problems, "starting_branch is empty")
	}
	if !headPattern.MatchString(id.StartingHead) {
		problems = append(problems, "starting_head is not a Git object id")
	}
	return problems
}

func eventHash(event RunEvent) (string, error) {
	event.EventHash = ""
	b, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("canonicalize event: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func runIdentityHash(identity RunIdentity) (string, error) {
	b, err := json.Marshal(identity)
	if err != nil {
		return "", fmt.Errorf("canonicalize run identity: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func verifyLedger(identity RunIdentity, path, closedAt string, expectedSequence uint64, expectedHash string) LedgerVerification {
	result := LedgerVerification{TailHash: genesisEventHash}
	if problems := identity.validate(); len(problems) > 0 {
		result.Problems = append(result.Problems, problems...)
		return result
	}
	f, err := os.Open(path)
	if err != nil {
		result.Problems = append(result.Problems, fmt.Sprintf("open ledger: %v", err))
		return result
	}
	defer f.Close()

	created, _ := time.Parse(time.RFC3339Nano, identity.CreatedAt)
	var closed time.Time
	if closedAt != "" {
		var parseErr error
		closed, parseErr = time.Parse(time.RFC3339Nano, closedAt)
		if parseErr != nil {
			result.Problems = append(result.Problems, "invalid run closure time: "+parseErr.Error())
		}
	}

	seenIDs := map[string]bool{}
	expectedSeq := uint64(1)
	expectedPrevious := genesisEventHash
	reader := bufio.NewReader(f)
	for lineNumber := 1; ; lineNumber++ {
		line, readErr := reader.ReadBytes('\n')
		if errors.Is(readErr, io.EOF) {
			if len(line) != 0 {
				result.Problems = append(result.Problems,
					fmt.Sprintf("line %d: truncated final record (missing newline)", lineNumber))
			}
			break
		}
		if readErr != nil {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: read ledger: %v", lineNumber, readErr))
			break
		}
		line = bytes.TrimSuffix(line, []byte{'\n'})
		var event RunEvent
		if err := decodeStrictJSON(line, &event); err != nil {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: malformed JSON: %v", lineNumber, err))
			continue
		}
		problems := validateEvent(event, identity, created, closed)
		for _, problem := range problems {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: %s", lineNumber, problem))
		}
		if event.Sequence != expectedSeq {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: sequence %d, want %d", lineNumber, event.Sequence, expectedSeq))
		}
		if seenIDs[event.EventID] {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: duplicate event_id %q", lineNumber, event.EventID))
		}
		seenIDs[event.EventID] = true
		if event.PreviousHash != expectedPrevious {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: previous_hash %s, want %s", lineNumber, event.PreviousHash, expectedPrevious))
		}
		computed, err := eventHash(event)
		if err != nil {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: compute event hash: %v", lineNumber, err))
		} else if event.EventHash != computed {
			result.Problems = append(result.Problems,
				fmt.Sprintf("line %d: event_hash %s, canonical hash is %s", lineNumber, event.EventHash, computed))
		}

		result.Events = append(result.Events, event)
		result.TailSequence = event.Sequence
		result.TailHash = event.EventHash
		expectedSeq++
		expectedPrevious = event.EventHash
	}

	if expectedSequence != 0 || expectedHash != "" {
		if result.TailSequence != expectedSequence || result.TailHash != expectedHash {
			result.Problems = append(result.Problems,
				fmt.Sprintf("ledger tail is sequence %d hash %s; closed state records sequence %d hash %s",
					result.TailSequence, result.TailHash, expectedSequence, expectedHash))
		}
	}
	result.Valid = len(result.Problems) == 0
	if !result.Valid {
		result.Events = nil
	}
	return result
}

func decodeStrictJSON(line []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(line))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func validateEvent(event RunEvent, identity RunIdentity, created, closed time.Time) []string {
	var problems []string
	if event.SchemaVersion != RunEventSchemaVersion {
		problems = append(problems, fmt.Sprintf("unsupported event schema_version %q", event.SchemaVersion))
	}
	if event.Sequence == 0 {
		problems = append(problems, "sequence must be positive")
	}
	if strings.TrimSpace(event.EventID) == "" {
		problems = append(problems, "event_id is empty")
	}
	if !hashPattern.MatchString(event.PreviousHash) {
		problems = append(problems, "previous_hash is not SHA-256")
	}
	if !hashPattern.MatchString(event.EventHash) {
		problems = append(problems, "event_hash is not SHA-256")
	}
	if event.WorkflowID != identity.WorkflowID {
		problems = append(problems, "workflow_id does not match run identity")
	}
	if event.RunID != identity.RunID {
		problems = append(problems, "run_id does not match run identity")
	}
	if event.ContractHash != identity.ContractHash {
		problems = append(problems, "contract_hash does not match run identity")
	}
	observed, err := time.Parse(time.RFC3339Nano, event.ObservedAt)
	if err != nil {
		problems = append(problems, "invalid observed_at: "+err.Error())
	} else {
		if observed.Before(created) {
			problems = append(problems, "event predates run creation")
		}
		if !closed.IsZero() && observed.After(closed) {
			problems = append(problems, "event is after terminal closure")
		}
	}
	if strings.TrimSpace(event.Provider) == "" {
		problems = append(problems, "provider is empty")
	}
	if !sourceKinds[event.SourceKind] {
		problems = append(problems, fmt.Sprintf("unsupported source_kind %q", event.SourceKind))
	}
	if strings.TrimSpace(event.CollectorID) == "" {
		problems = append(problems, "collector_id is empty")
	}
	if !trustBases[event.TrustBasis] {
		problems = append(problems, fmt.Sprintf("unsupported trust_basis %q", event.TrustBasis))
	}
	if event.RepositoryIdentity != identity.RepositoryIdentity {
		problems = append(problems, "repository_identity does not match run identity")
	}
	if !pathWithin(identity.RepositoryRoot, event.CWD) {
		problems = append(problems, "cwd is not inside the run repository")
	}
	if strings.TrimSpace(event.Branch) == "" {
		problems = append(problems, "branch is empty")
	}
	if !headPattern.MatchString(event.Head) {
		problems = append(problems, "head is not a Git object id")
	}
	if !eventKinds[event.EventKind] {
		problems = append(problems, fmt.Sprintf("unsupported event_kind %q", event.EventKind))
	}
	return problems
}

func pathWithin(root, cwd string) bool {
	if !filepath.IsAbs(cwd) {
		return false
	}
	canonicalCWD, err := canonicalPath(cwd)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(root, canonicalCWD)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// evaluateEventEligibility decides whether a structurally valid event can
// satisfy a reconciliation requirement. Hash-chain validity never changes an
// event's provenance or trust basis.
func evaluateEventEligibility(event RunEvent) EventEligibility {
	if event.EventKind == EventOutputObserved {
		return EventEligibility{Eligible: true, ObservationKind: ObsArtifactPresent,
			Reason: "output observation is eligible only as artifact evidence"}
	}
	if event.SourceKind == SourceManualImport || event.SourceKind == SourceLegacyTranscript ||
		event.SourceKind == SourceUnknown {
		return EventEligibility{Reason: fmt.Sprintf("%s provenance cannot satisfy an invocation", event.SourceKind)}
	}
	if event.SourceKind == SourceOperatorRecord {
		return EventEligibility{Reason: "operator-record eligibility is outside this atomic change"}
	}
	if event.SourceKind == SourceProviderHook {
		if event.TrustBasis != TrustCollectorObserved {
			return EventEligibility{Reason: "provider-hook invocation is not collector-observed"}
		}
		if event.SessionID == "" || event.TurnID == "" {
			return EventEligibility{Reason: "provider-hook invocation lacks session or turn identity"}
		}
		switch event.EventKind {
		case EventAgentStarted, EventAgentFinished:
			if event.Selector == "" {
				return EventEligibility{Reason: "agent event has no selector"}
			}
			return EventEligibility{Eligible: true, ObservationKind: ObsAgentCall,
				Reason: "collector-observed provider agent event"}
		case EventSkillSelected:
			if event.Selector == "" {
				return EventEligibility{Reason: "skill event has no selector"}
			}
			return EventEligibility{Eligible: true, ObservationKind: ObsSkillCall,
				Reason: "collector-observed provider skill event"}
		case EventToolFinished:
			if event.Selector == "" {
				return EventEligibility{Reason: "tool event has no selector"}
			}
			return EventEligibility{Eligible: true, ObservationKind: ObsBackendRun,
				Reason: "collector-observed provider tool completion"}
		default:
			return EventEligibility{Reason: fmt.Sprintf("event_kind %s is not invocation evidence", event.EventKind)}
		}
	}
	if event.SourceKind == SourceDeterministicBackend {
		if event.TrustBasis != TrustBackendExitStatus {
			return EventEligibility{Reason: "backend event lacks a captured-exit-status trust basis"}
		}
		if event.EventKind != EventBackendFinished && event.EventKind != EventValidationFinish {
			return EventEligibility{Reason: fmt.Sprintf("event_kind %s is not a backend completion", event.EventKind)}
		}
		if event.Selector == "" {
			return EventEligibility{Reason: "backend event has no selector"}
		}
		return EventEligibility{Eligible: true, ObservationKind: ObsBackendRun,
			Reason: "deterministic backend completion with captured exit status field"}
	}
	return EventEligibility{Reason: "source provenance is not eligible"}
}

func evidenceFromLedger(verification LedgerVerification) Evidence {
	if !verification.Valid {
		why := strings.Join(verification.Problems, "; ")
		return Evidence{
			Incomplete: map[ObservationKind]string{
				ObsAgentCall: why, ObsSkillCall: why, ObsBackendRun: why,
				ObsOperatorAction: why, ObsArtifactPresent: why,
			},
			IntegrityErrors: append([]string(nil), verification.Problems...),
		}
	}
	ev := Evidence{}
	for _, event := range verification.Events {
		eligibility := evaluateEventEligibility(event)
		if !eligibility.Eligible {
			ev.Limitations = append(ev.Limitations,
				fmt.Sprintf("event %s is ineligible: %s", event.EventID, eligibility.Reason))
			continue
		}
		if eligibility.ObservationKind == ObsArtifactPresent {
			ev.Limitations = append(ev.Limitations,
				fmt.Sprintf("event %s is artifact evidence only; it cannot prove an invocation", event.EventID))
		}
		ref := fmt.Sprintf("ledger:%s:%d:%s:%s", event.RunID, event.Sequence, event.EventID, event.EventHash)
		ev.Observations = append(ev.Observations, Observation{
			Kind: eligibility.ObservationKind, Selector: event.Selector,
			ExitStatus: event.ExitStatus, Timestamp: event.ObservedAt, EvidenceRef: ref,
		})
	}
	return ev
}

func marshalJSONLine(value any) ([]byte, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func writeJSONExclusive(path string, value any, dirPerm, filePerm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("create directory for %s: %w", path, err)
	}
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	return writeFileExclusive(path, append(b, '\n'), filePerm)
}

func writeFileExclusive(path string, contents []byte, perm os.FileMode) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	ok := false
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close %s: %w", path, closeErr)
		}
		if !ok {
			_ = os.Remove(path)
		}
	}()
	if _, err = f.Write(contents); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err = f.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", path, err)
	}
	ok = true
	return nil
}

func writeJSONAtomic(path string, value any, perm os.FileMode) error {
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	b = append(b, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".state-*")
	if err != nil {
		return fmt.Errorf("create temporary state for %s: %w", path, err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return fmt.Errorf("chmod temporary state: %w", err)
	}
	if _, err := tmp.Write(b); err != nil {
		tmp.Close()
		return fmt.Errorf("write temporary state: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temporary state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary state: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return syncDirectory(filepath.Dir(path))
}

func syncDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open directory %s for sync: %w", path, err)
	}
	defer dir.Close()
	if err := dir.Sync(); err != nil {
		return fmt.Errorf("sync directory %s: %w", path, err)
	}
	return nil
}
