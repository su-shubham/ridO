package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "rule-engine/engine"
)

// Request structure for creating a new rule
type CreateRuleRequest struct {
    Name       string `json:"name" binding:"required"`
    RuleString string `json:"rule_string" binding:"required"`
}

// Request structure for evaluating rules
type EvaluateRuleRequest struct {
    RuleIDs  []string               `json:"rule_ids" binding:"required"`
    UserData map[string]interface{} `json:"user_data" binding:"required"`
}

// SetupRoutes initializes all API routes
func SetupRoutes(router *gin.Engine, ruleEngine *engine.RuleEngine) {
    router.POST("/api/create_rule", createRuleHandler(ruleEngine))
    router.GET("/api/get_rules", getRulesHandler(ruleEngine))
    router.POST("/api/evaluate_rule", evaluateRuleHandler(ruleEngine))
}

// Handler for creating a new rule
func createRuleHandler(ruleEngine *engine.RuleEngine) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateRuleRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := ruleEngine.CreateRule(req.Name, req.RuleString); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Rule created successfully"})
    }
}

// Handler for fetching all rules
func getRulesHandler(ruleEngine *engine.RuleEngine) gin.HandlerFunc {
    return func(c *gin.Context) {
        rules, err := ruleEngine.GetRules()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"rules": rules})
    }
}

// Handler for evaluating rules
func evaluateRuleHandler(ruleEngine *engine.RuleEngine) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req EvaluateRuleRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        results, err := ruleEngine.EvaluateRule(req.RuleIDs, req.UserData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"results": results})
    }
}
