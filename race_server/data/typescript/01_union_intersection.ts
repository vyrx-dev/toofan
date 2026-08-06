// Topic: Union Types and Type Guards

type Success = { status: "success"; data: string[] };
type Failure = { status: "error"; message: string };
type ApiResponse = Success | Failure;

function handleResponse(res: ApiResponse): string {
    if (res.status === "success") {
        return res.data.join(", ");
    }
    return `Error: ${res.message}`;
}

type Shape = { kind: "circle"; radius: number }
    | { kind: "rect"; width: number; height: number };

function area(shape: Shape): number {
    switch (shape.kind) {
        case "circle": return Math.PI * shape.radius ** 2;
        case "rect": return shape.width * shape.height;
    }
}