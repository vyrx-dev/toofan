# Topic: Conditionals and Arguments
if [ -z "$1" ]; then
    echo "usage: $0 <name>"
    exit 1
fi

name="$1"
count="${2:-5}"

if [ "$count" -gt 10 ]; then
    echo "too many, max 10"
    exit 1
fi

for i in $(seq 1 "$count"); do
    echo "$i: hello $name"
done
