//go:build ignore

// Topic: Maps and Iteration
scores := map[string]int{
    "alice": 92,
    "bob":   85,
    "eve":   97,
}

for name, score := range scores {
    if score >= 90 {
        fmt.Printf("%s passed with %d\n", name, score)
    }
}

delete(scores, "bob")
scores["charlie"] = 88
fmt.Println("total:", len(scores))
