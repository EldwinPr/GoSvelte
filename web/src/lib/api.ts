export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
    const response = await fetch(path, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Something went wrong');
    }

    if (response.status === 204) {
        return {} as T;
    }

    return response.json();
}
