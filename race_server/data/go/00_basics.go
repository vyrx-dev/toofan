//go:build ignore

// Topic: Variables and Types
name := "toofan"
version := 1.0
active := true

fmt.Printf("%s v%.1f (active: %t)\n", name, version, active)

nums := []int{3, 1, 4, 1, 5}
sort.Ints(nums)
fmt.Println(nums)
