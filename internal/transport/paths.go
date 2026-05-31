// Package transport provides low-level HTTP helpers shared across all services.
package transport

import "fmt"

const apiV1 = "/api/open/v1"

// Conversation paths
const PathConversations = apiV1 + "/conversations"

func PathConversation(id string) string    { return PathConversations + "/" + id }
func PathMessages(id string) string        { return PathConversation(id) + "/messages" }
func PathChat(id string) string            { return PathConversation(id) + "/chat" }
func PathChatStop(id string) string        { return PathChat(id) + "/stop" }
func PathChatHitlDecide(id string) string  { return PathChat(id) + "/hitl/decide" }

// Cron paths
const PathCronTasks = apiV1 + "/cron/tasks"

func PathCronTask(taskID string) string       { return PathCronTasks + "/" + taskID }
func PathCronTaskPause(taskID string) string  { return PathCronTask(taskID) + "/pause" }
func PathCronTaskResume(taskID string) string { return PathCronTask(taskID) + "/resume" }

// Agent paths (cloud mode only)
const PathAgents = apiV1 + "/agent"

func PathAgent(agentID int) string      { return fmt.Sprintf("%s/%d", PathAgents, agentID) }
func PathAgentMates(agentID int) string { return PathAgent(agentID) + "/mates" }
