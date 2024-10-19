package api

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "rule-engine/engine"
)

func setupRouter() *gin.Engine {
    r := gin.Default()
    ruleEngine := engine.NewRuleEngine() // Assuming you have a constructor for RuleEngine
    SetupRoutes(r, ruleEngine)
    return r
}

func TestCreateRuleAndVerifyAST(t *testing.T) {
    router := setupRouter()

    // Create an individual rule
    ruleJSON := `{"name": "Temperature Check", "rule_string": "temperature > 30"}`
    req, _ := http.NewRequest(http.MethodPost, "/api/create_rule", strings.NewReader(ruleJSON))
    req.Header.Set("Content-Type", "application/json")
    response := httptest.NewRecorder()

    router.ServeHTTP(response, req)

    assert.Equal(t, http.StatusCreated, response.Code)

    // Verify the AST representation
    // Add your code to verify the AST representation here
    // Example: 
    // ast := ruleEngine.GetAST("Temperature Check") // Hypothetical function
    // assert.Equal(t, expectedAST, ast)
}

func TestCombineRulesAndVerifyAST(t *testing.T) {
    router := setupRouter()

    // Create two individual rules
    ruleJSON1 := `{"name": "Temperature Check", "rule_string": "temperature > 30"}`
    ruleJSON2 := `{"name": "Humidity Check", "rule_string": "humidity < 80"}`

    // Create the first rule
    req1, _ := http.NewRequest(http.MethodPost, "/api/create_rule", strings.NewReader(ruleJSON1))
    req1.Header.Set("Content-Type", "application/json")
    response1 := httptest.NewRecorder()
    router.ServeHTTP(response1, req1)

    // Create the second rule
    req2, _ := http.NewRequest(http.MethodPost, "/api/create_rule", strings.NewReader(ruleJSON2))
    req2.Header.Set("Content-Type", "application/json")
    response2 := httptest.NewRecorder()
    router.ServeHTTP(response2, req2)

    // Combine rules
    combineJSON := `{"rule_ids": ["1", "2"]}` // Adjust rule IDs as necessary
    combineReq, _ := http.NewRequest(http.MethodPost, "/api/combine_rules", strings.NewReader(combineJSON))
    combineReq.Header.Set("Content-Type", "application/json")
    combineResponse := httptest.NewRecorder()
    router.ServeHTTP(combineResponse, combineReq)

    assert.Equal(t, http.StatusOK, combineResponse.Code)

    // Verify the combined AST representation
    // Add your code to verify the combined AST representation here
    // Example: 
    // combinedAST := ruleEngine.GetCombinedAST() // Hypothetical function
    // assert.Equal(t, expectedCombinedAST, combinedAST)
}

func TestEvaluateRuleWithSampleData(t *testing.T) {
    router := setupRouter()

    // Create a rule
    ruleJSON := `{"name": "Temperature Check", "rule_string": "temperature > 30"}`
    req, _ := http.NewRequest(http.MethodPost, "/api/create_rule", strings.NewReader(ruleJSON))
    req.Header.Set("Content-Type", "application/json")
    response := httptest.NewRecorder()
    router.ServeHTTP(response, req)

    // Evaluate rule with sample data
    evaluateJSON := `{"rule_ids": ["1"], "user_data": {"temperature": 35}}` // Adjust rule ID as necessary
    evalReq, _ := http.NewRequest(http.MethodPost, "/api/evaluate_rule", strings.NewReader(evaluateJSON))
    evalReq.Header.Set("Content-Type", "application/json")
    evalResponse := httptest.NewRecorder()
    router.ServeHTTP(evalResponse, evalReq)

    assert.Equal(t, http.StatusOK, evalResponse.Code)
    assert.JSONEq(t, `{"1": true}`, evalResponse.Body.String())
}

func TestCombineAdditionalRules(t *testing.T) {
    router := setupRouter()

    // Create additional rules and combine
    // Repeat the process similar to above for additional rules
    // Assert the expected outcome
}
