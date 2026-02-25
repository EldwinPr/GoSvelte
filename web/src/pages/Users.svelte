<script lang="ts">
	import { UserPlus, Trash2 } from 'lucide-svelte';
	import { api } from '../lib/api';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface User {
		id?: string;
		name: string;
		email: string;
		password?: string;
		clearance: number;
	}

	let tableRef: any = $state();

	// Form state for new user
	let newUser = $state<User>({
		name: '',
		email: '',
		password: '',
		clearance: 0
	});

	async function registerUser(e: SubmitEvent) {
		e.preventDefault();
		try {
			await api('/api/register', {
				method: 'POST',
				body: JSON.stringify(newUser)
			});
			// Reset form and refresh list
			newUser = { name: '', email: '', password: '', clearance: 0 };
			tableRef.fetchData();
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

	const columns: Column[] = [
		{ key: 'name', label: 'User Details', sortable: true },
		{ key: 'clearance', label: 'Clearance', sortable: true, align: 'center', width: '140px' },
		{ key: 'actions', label: 'Actions', align: 'right', width: '100px' }
	];
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">User Management</h2>
			<p class="text-slate-500 text-sm font-medium italic opacity-70">Manage system access and security clearance levels.</p>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
		<!-- Registration Form -->
		<section class="card p-8 bg-white border border-slate-200 space-y-6 h-fit shadow-xl">
			<div class="flex items-center gap-2 border-b border-slate-100 pb-4">
				<UserPlus class="text-teal-600" size={24} />
				<h3 class="h3 text-slate-800 font-black uppercase tracking-tight">Register User</h3>
			</div>
			
			<form onsubmit={registerUser} class="space-y-4 text-left">
				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest">Full Name</span>
					<input class="input" type="text" bind:value={newUser.name} placeholder="John Doe" required />
				</label>
				
				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest">Email Address</span>
					<input class="input" type="email" bind:value={newUser.email} placeholder="user@system.com" required />
				</label>
				
				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest">Password</span>
					<input class="input" type="password" bind:value={newUser.password} placeholder="••••••••" required />
				</label>
				
				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest">Clearance Level</span>
					<select class="select" bind:value={newUser.clearance}>
						<option value={0}>0 - Finance</option>
						<option value={10}>10 - Finance Manager</option>
						<option value={20}>20 - Developer</option>
					</select>
				</label>
				
				<button type="submit" class="btn bg-teal-600 text-white w-full mt-4 font-black uppercase tracking-widest py-4 shadow-lg shadow-teal-600/20">Create Account</button>
			</form>
		</section>

		<!-- Users Table -->
		<div class="lg:col-span-2">
			<DataTable 
				bind:this={tableRef}
				endpoint="/api/users" 
				{columns} 
				rowKey="id" 
				searchPlaceholder="Cari nama atau email..."
			>
				{#snippet cell(rowData, key)}
					{@const row = rowData as User}
					{#if key === 'name'}
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 rounded-full bg-teal-100 text-teal-700 flex items-center justify-center font-black shrink-0 border border-teal-200">
								{row.name?.charAt(0) || '?'}
							</div>
							<div class="min-w-0 text-left">
								<p class="font-black text-slate-800 text-sm">{row.name || 'Unknown'}</p>
								<p class="text-xs text-slate-500 font-medium">{row.email || 'No Email'}</p>
							</div>
						</div>
					{:else if key === 'clearance'}
						{@const info = getClearanceLabel(row.clearance)}
						<span class={info.color}>
							{info.label}
						</span>
					{:else if key === 'actions'}
						<button class="text-slate-300 hover:text-red-600 transition-colors p-2" title="Delete User">
							<Trash2 size={16} />
						</button>
					{/if}
				{/snippet}
			</DataTable>
		</div>
	</div>
</div>
