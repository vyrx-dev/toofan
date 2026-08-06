# Topic: Loops and File Processing
for file in *.txt; do
    lines=$(wc -l < "$file")
    echo "$file: $lines lines"
done

while read -r line; do
    key=$(echo "$line" | cut -d= -f1)
    val=$(echo "$line" | cut -d= -f2)
    export "$key=$val"
done < config.env

find . -name '*.log' -mtime +7 \
    -exec rm -f {} \;
echo "cleaned old logs"
