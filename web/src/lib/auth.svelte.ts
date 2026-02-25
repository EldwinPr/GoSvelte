import { api } from './api';

export interface User {
	id?: string;
	name: string;
	email: string;
	clearance: number;
}

class AuthStore {
	#user = $state<User | null>(null);

	get user() { return this.#user; }
	set user(value: User | null) { this.#user = value; }

	async init() {
		try {
			this.#user = await api<User>('/api/me');
		} catch (err) {
			this.#user = null;
		}
	}

	async login(email: string, password: string) {
		const user = await api<User>('/api/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		this.#user = user;
		return user;
	}

	async logout() {
		try {
			await api('/api/logout', { method: 'POST' });
		} finally {
			this.#user = null;
			window.location.hash = '#/login';
		}
	}
}

export const auth = new AuthStore();
