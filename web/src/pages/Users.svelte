<script lang="ts">
	import { onMount } from 'svelte';
	import { UserPlus, Trash2 } from 'lucide-svelte';
	import { api } from '../lib/api';

	interface User {
		id?: string;
		name: string;
		email: string;
		password?: string;
		clearance: number;
	}

	let users = $state<User[]>([]);
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	// Form state for new user
	let newUser = $state<User>({
		name: '',
		email: '',
		password: '',
		clearance: 0
	});

	async function fetchUsers() {
		try {
			users = await api<User[]>('/api/users');
		} catch (err: any) {
			error = err.message || "Failed to load users. Ensure you are logged in as a Developer.";
		} finally {
			isLoading = false;
		}
	}

	async function registerUser(e: SubmitEvent) {
		e.preventDefault();
		try {
			await api('/api/register', {
				method: 'POST',
				body: JSON.stringify(newUser)
			});
			// Reset form and refresh list
			newUser = { name: '', email: '', password: '', clearance: 0 };
			fetchUsers();
		} catch (err: any) {
			alert("Error registering user: " + err.message);
		}
	}

	function getClearanceLabel(level: number) {
		const lvl = level || 0;
		if (lvl >= 20) return { label: 'Developer', color: 'badge preset-filled-error-500' };
		if (lvl >= 10) return { label: 'Manager', color: 'badge preset-filled-warning-500' };
		return { label: 'Finance', color: 'badge preset-filled-success-500' };
	}

	onMount(fetchUsers);
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 text-slate-900">User Management</h2>
			<p class="text-slate-500 text-sm">Manage system access and security clearance levels.</p>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
		<!-- Registration Form -->
		<section class="card p-6 bg-white border border-slate-200 space-y-4 h-fit shadow-sm">
			<div class="flex items-center gap-2 mb-4">
				<UserPlus class="text-primary-500" size={24} />
				<h3 class="h3 text-slate-800 font-bold">Register User</h3>
			</div>
			
			<form onsubmit={registerUser} class="space-y-4">
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Full Name</span>
					<input class="input" type="text" bind:value={newUser.name} placeholder="John Doe" required />
				</label>
				
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Email Address</span>
					<input class="input" type="email" bind:value={newUser.email} placeholder="john@erp.com" required />
				</label>
				
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Password</span>
					<input class="input" type="password" bind:value={newUser.password} placeholder="••••••••" required />
				</label>
				
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Clearance Level</span>
					<select class="select" bind:value={newUser.clearance}>
						<option value={0}>0 - Finance</option>
						<option value={10}>10 - Finance Manager</option>
						<option value={20}>20 - Developer</option>
					</select>
				</label>
				
				<button type="submit" class="btn preset-filled-primary-500 w-full mt-2 font-bold shadow-md shadow-primary-500/20">Create Account</button>
			</form>
		</section>

		<!-- Users Table -->
		<section class="lg:col-span-2 card bg-white border border-slate-200 overflow-hidden shadow-sm">
			{#if isLoading}
				<div class="p-12 text-center text-slate-400">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent mb-4"></div>
					<p>Loading users...</p>
				</div>
			{:else if error}
				<div class="p-12 text-center text-error-600 bg-error-50">
					<p class="font-bold mb-2">Access Denied</p>
					<p class="text-sm">{error}</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full text-left">
						<thead class="bg-slate-50 text-slate-600 font-bold border-b border-slate-200">
							<tr>
								<th class="p-4">User Details</th>
								<th class="p-4">Clearance</th>
								<th class="p-4 text-right">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each users as user (user.id)}
								{@const info = getClearanceLabel(user.clearance)}
								<tr class="hover:bg-slate-50 transition-colors">
									<td class="p-4">
										<div class="flex items-center gap-3">
											<div class="w-10 h-10 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center font-bold shrink-0">
												{user.name?.charAt(0) || '?'}
											</div>
											<div class="min-w-0">
												<p class="font-bold text-slate-800 truncate text-sm">{user.name || 'Unknown'}</p>
												<p class="text-xs text-slate-500 truncate">{user.email || 'No Email'}</p>
											</div>
										</div>
									</td>
									<td class="p-4">
										<span class={info.color}>
											{info.label}
										</span>
									</td>
									<td class="p-4 text-right">
										<button class="btn-icon btn-icon-sm text-slate-400 hover:text-error-500 hover:bg-error-50 transition-all" title="Delete User">
											<Trash2 size={16} />
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</section>
	</div>
</div>
