// Topic: Generics with Constraints

function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
    return obj[key];
}

interface Repository<T extends { id: string }> {
    findById(id: string): Promise<T>;
    save(entity: T): Promise<void>;
}