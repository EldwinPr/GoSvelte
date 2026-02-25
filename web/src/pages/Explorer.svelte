<script lang="ts">
	import { onMount } from 'svelte';
	import { Database, RefreshCw, AlertCircle } from 'lucide-svelte';
	import { api } from '../lib/api';

	let tables = [
		{ id: 'users', name: 'Users' },
		{ id: 'customer_credits', name: 'Customer Credits' },
		{ id: 'invoices', name: 'Invoices' },
		{ id: 'invoice_details', name: 'Invoice Details' },
		{ id: 'invoice_payments', name: 'Invoice Payments' },
		{ id: 'company_transactions', name: 'Company Transactions' },
		{ id: 'requisitions', name: 'Requisitions' },
		{ id: 'company_balances', name: 'Company Balances' }
	];

	let selectedTable = $state(tables[0].id);
	let data = $state<any[]>([]);
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	let columns = $derived(data && data.length > 0 ? Object.keys(data[0]) : []);

	async function fetchData() {
		if (isLoading) return;
		isLoading = true;
		error = null;
		try {
			const result = await api<any[]>(`/api/explorer?table=${selectedTable}`);
			data = result || [];
		} catch (err: any) {
			error = err.message;
			data = [];
		} finally {
			isLoading = false;
		}
	}

	onMount(fetchData);
</script>

<div class="space-y-6 bg-white p-2 rounded-xl">
	<header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 flex items-center gap-2 text-slate-900">
				<Database class="text-primary-500" />
				Database Explorer
			</h2>
			<p class="text-slate-500 text-sm italic">Direct read access to system tables (Developer Only).</p>
		</div>
		
		<div class="flex items-center gap-2">
			<select class="select w-auto min-w-200 bg-slate-50 border-slate-200" bind:value={selectedTable} onchange={fetchData}>
				{#each tables as table (table.id)}
					<option value={table.id}>{table.name}</option>
				{/each}
			</select>
			<button class="btn preset-filled-primary-500 shadow-md shadow-primary-500/20" onclick={fetchData} disabled={isLoading}>
				<RefreshCw size={18} class={isLoading ? 'animate-spin' : ''} />
				<span>Refresh</span>
			</button>
		</div>
	</header>

	{#if error}
		<aside class="alert preset-tonal-error">
			<AlertCircle size={24} />
			<div class="alert-message">
				<h3 class="font-bold">Error</h3>
				<p>{error}</p>
			</div>
		</aside>
	{/if}

	<section class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
		{#if isLoading}
			<div class="p-20 text-center text-slate-400">
				<div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-primary-500 border-t-transparent mb-4"></div>
				<p class="font-medium">Querying {selectedTable}...</p>
			</div>
		{:else if data.length === 0}
			<div class="p-20 text-center text-slate-400 italic bg-slate-50">
				No records found in {selectedTable}.
			</div>
		{:else}
			<div class="overflow-x-auto max-h-[calc(100vh-250px)] overflow-y-auto">
				<table class="table table-hover table-compact w-full text-xs">
					<thead class="bg-slate-100 sticky top-0 z-10 border-b border-slate-200">
						<tr>
							{#each columns as col (col)}
								<th class="p-3 uppercase font-bold text-slate-600 tracking-wider">
									{col.replace(/_/g, ' ')}
								</th>
							{/each}
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each data as row (row.id || JSON.stringify(row))}
							<tr class="hover:bg-slate-50 transition-colors">
								{#each columns as col (col)}
									<td class="p-3 whitespace-nowrap overflow-hidden text-ellipsis max-w-[200px] text-slate-700" title={String(row[col])}>
										{row[col] === null ? 'NULL' : String(row[col])}
									</td>
								{/each}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</section>
</div>
