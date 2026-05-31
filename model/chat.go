package model

// ChatRequest is the full request body for a chat turn.
type ChatRequest struct {
	Message       string                 `json:"message"`
	AgentID       string                 `json:"agent_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	HitlDecisions []HitlDecision         `json:"hitl_decisions,omitempty"`
}

// HitlDecision is a Human-In-The-Loop decision submitted to resume a paused run.
type HitlDecision struct {
	Decision     string                 `json:"decision"` // "approve" | "reject" | "modify"
	Reason       string                 `json:"reason,omitempty"`
	ModifiedArgs map[string]interface{} `json:"modified_args,omitempty"`
}

// Approve returns an approve decision.
func Approve() HitlDecision { return HitlDecision{Decision: "approve"} }

// Reject returns a reject decision with the given reason.
func Reject(reason string) HitlDecision { return HitlDecision{Decision: "reject", Reason: reason} }

// Modify returns a modify decision with replacement arguments.
func Modify(args map[string]interface{}) HitlDecision {
	return HitlDecision{Decision: "modify", ModifiedArgs: args}
}

// HitlDecideRequest is the wire body for the HITL decide endpoint.
type HitlDecideRequest struct {
	Decisions []HitlDecision `json:"decisions"`
}
