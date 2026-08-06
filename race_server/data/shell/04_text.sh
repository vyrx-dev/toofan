# Topic: Text Processing
grep -rn "TODO" --include="*.go" . | while read -r match; do
    file=$(echo "$match" | cut -d: -f1)
    line=$(echo "$match" | cut -d: -f2)
    echo "[$file:$line] $(echo "$match" | cut -d: -f3-)"
done

wc -l *.go | sort -rn | head -5
