package engine

import (
    "encoding/json"
    "fmt"
    "go/ast"
    "go/parser"
    "strconv"
    "rule-engine/storage"
)

type Node struct {
    Type  string `json:"type"`
    Value string `json:"value,omitempty"`
    Left  *Node  `json:"left,omitempty"`
    Right *Node  `json:"right,omitempty"`
}

type RuleEngine struct {
    store storage.RuleStore
}

func NewRuleEngine(store storage.RuleStore) *RuleEngine {
    return &RuleEngine{store: store}
}

func (e *RuleEngine) CreateRule(name, ruleString string) error {
    if name == "" || ruleString == "" {
        return fmt.Errorf("name and rule string cannot be empty")
    }

    node, err := createRule(ruleString)
    if err != nil {
        return fmt.Errorf("failed to create rule: %w", err)
    }

    nodeJSON, err := json.Marshal(node)
    if err != nil {
        return fmt.Errorf("failed to marshal rule: %w", err)
    }

    if _, err := e.store.StoreRule(name, string(nodeJSON)); err != nil {
        return fmt.Errorf("failed to store rule: %w", err)
    }

    return nil
}

func (e *RuleEngine) GetRules() ([]storage.Rule, error) {
    return e.store.GetRules()
}

func (e *RuleEngine) EvaluateRule(ruleIDs []string, userData map[string]interface{}) (bool, error) {
    rules, err := e.store.GetRulesByIDs(ruleIDs)
    if err != nil {
        return false, err
    }

    var nodes []*Node
    for _, rule := range rules {
        var node Node
        if err := json.Unmarshal([]byte(rule.RuleString), &node); err != nil {
            return false, err
        }
        nodes = append(nodes, &node)
    }

    return evaluateRule(combineRules(nodes), userData), nil
}

func createRule(ruleString string) (*Node, error) {
    expr, err := parser.ParseExpr(ruleString)
    if err != nil {
        return nil, fmt.Errorf("invalid rule syntax: %v", err)
    }
    return parseExpression(expr)
}

func parseExpression(expr ast.Expr) (*Node, error) {
    switch e := expr.(type) {
    case *ast.BinaryExpr:
        left, err := parseExpression(e.X)
        if err != nil {
            return nil, err
        }
        right, err := parseExpression(e.Y)
        if err != nil {
            return nil, err
        }
        return &Node{Type: "operator", Value: map[string]string{"&&": "AND", "||": "OR"}[e.Op.String()], Left: left, Right: right}, nil
    case *ast.ParenExpr:
        return parseExpression(e.X)
    case *ast.Ident:
        return &Node{Type: "operand", Value: e.Name}, nil
    case *ast.BasicLit:
        return &Node{Type: "operand", Value: e.Value}, nil
    default:
        return nil, fmt.Errorf("unsupported expression type: %T", expr)
    }
}

func combineRules(nodes []*Node) *Node {
    if len(nodes) == 0 {
        return nil
    }
    combined := nodes[0]
    for _, node := range nodes[1:] {
        combined = &Node{Type: "operator", Value: "AND", Left: combined, Right: node}
    }
    return combined
}

func evaluateRule(node *Node, data map[string]interface{}) bool {
    if node == nil {
        return false
    }

    switch node.Type {
    case "operator":
        left := evaluateRule(node.Left, data)
        right := evaluateRule(node.Right, data)
        if node.Value == "AND" {
            return left && right
        }
        return left || right
    case "operand":
        parts := splitOperand(node.Value)
        if len(parts) != 3 {
            return false
        }
        return compare(data[parts[0]], parts[1], parseValue(parts[2]))
    }
    return false
}

func splitOperand(s string) []string {
    var result []string
    var current string
    inQuotes := false

    for _, char := range s {
        if char == '\'' {
            inQuotes = !inQuotes
        } else if char == ' ' && !inQuotes {
            if current != "" {
                result = append(result, current)
                current = ""
            }
        } else {
            current += string(char)
        }
    }

    if current != "" {
        result = append(result, current)
    }

    return result
}

func parseValue(s string) interface{} {
    if i, err := strconv.Atoi(s); err == nil {
        return i
    }
    if f, err := strconv.ParseFloat(s, 64); err == nil {
        return f
    }
    return s
}

func compare(left interface{}, op string, right interface{}) bool {
    switch op {
    case ">":
        return compareGreater(left, right)
    case "<":
        return compareLess(left, right)
    case "=":
        return compareEqual(left, right)
    default:
        return false
    }
}

func compareGreater(left, right interface{}) bool {
    switch l := left.(type) {
    case int:
        if r, ok := right.(int); ok {
            return l > r
        }
    case float64:
        if r, ok := right.(float64); ok {
            return l > r
        }
    }
    return false
}

func compareLess(left, right interface{}) bool {
    switch l := left.(type) {
    case int:
        if r, ok := right.(int); ok {
            return l < r
        }
    case float64:
        if r, ok := right.(float64); ok {
            return l < r
        }
    }
    return false
}

func compareEqual(left, right interface{}) bool {
    return left == right
}
