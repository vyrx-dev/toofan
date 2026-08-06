// Topic: Template Literal Types

type HttpMethod = "get" | "post" | "put" | "delete";
type ApiRoute = `/api/${string}`;

type EventMap = {
    click: MouseEvent;
    keydown: KeyboardEvent;
    scroll: Event;
};

type EventHandler<T extends keyof EventMap> = (
    event: EventMap[T]
) => void;

function on<T extends keyof EventMap>(
    event: T,
    handler: EventHandler<T>
): void {
    document.addEventListener(event, handler as EventListener);
}