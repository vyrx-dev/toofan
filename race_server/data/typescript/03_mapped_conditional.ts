// Topic: Mapped and Conditional Types

type FormErrors<T> = {
    [K in keyof T]?: string;
};

type User = { name: string; email: string; age: number };

const errors: FormErrors<User> = {
    name: "required",
    email: "invalid format",
};

type NonNullableFields<T> = {
    [K in keyof T]: NonNullable<T[K]>;
};