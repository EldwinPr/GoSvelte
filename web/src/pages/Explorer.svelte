<script lang="ts">
	import { Database, Plus, Edit2, Trash2, X, Save, RefreshCcw, AlertTriangle } from 'lucide-svelte';
	import { api } from '../lib/api';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	const tables = [
		{ id: 'users', name: 'Users' },
		{ id: 'customer_credits', name: 'Customer Credits' },
		{ id: 'invoices', name: 'Invoices' },
		{ id: 'invoice_details', name: 'Invoice Details' },
		{ id: 'invoice_payments', name: 'Invoice Payments' },
		{ id: 'company_transactions', name: 'Company Transactions' },
		{ id: 'requisitions', name: 'Requisitions' },
		{ id: 'company_balances', name: 'Company Balances' }
	];

	// State
	let selectedTable = $state(tables[0].id);
	let tableRef: any = $state();
	
	// Schema-aware columns (EXCLUDING ID)
	const tableSchemas: Record<string, string[]> = {
		'users': ['name', 'email', 'clearance', 'deleted_at'],
		'customer_credits': ['customer_name', 'credit_limit', 'used_credit', 'deleted_at'],
		'invoices': ['number', 'customer_id', 'total_amount', 'status', 'deleted_at'],
		'invoice_details': ['invoice_id', 'description', 'quantity', 'unit_price', 'deleted_at'],
		'invoice_payments': ['invoice_id', 'amount', 'payment_date', 'method', 'deleted_at'],
		'company_transactions': ['date', 'description', 'amount', 'type', 'category', 'deleted_at'],
		'requisitions': ['name', 'amount', 'status', 'type', 'category', 'deleted_at'],
		'company_balances': ['account_name', 'balance', 'deleted_at']
	};

	let columns = $derived.by(() => {
		const baseCols = tableSchemas[selectedTable] || ['deleted_at'];
		const cols: Column[] = [
			{ key: 'actions', label: 'Actions', width: '80px', align: 'center' },
			...baseCols.map(key => ({
				key,
				label: key.replace(/_/g, ' '),
				sortable: true
			}))
		];
		return cols;
	});

	// CRUD States
	let showModal = $state(false);
	let isEditing = $state(false);
	let currentRecord = $state<any>({});
	let isSaving = $state(false);

	let editableColumns = $derived(columns.filter(c => !['actions', 'id', 'created_at', 'updated_at', 'deleted_at'].includes(c.key)).map(c => c.key));

	function openCreate() {
		isEditing = false;
		currentRecord = {};
		editableColumns.forEach(c => currentRecord[c] = '');
		showModal = true;
	}

	function openEdit(row: any) {
		isEditing = true;
		currentRecord = { ...row };
		showModal = true;
	}

	async function deleteRecord(id: any) {
		if (!confirm('Are you sure you want to delete this record?')) return;
		try {
			await api(`/api/explorer?table=${selectedTable}&id=${id}`, { method: 'DELETE' });
			tableRef.fetchData();
		} catch (err: any) {
			alert('Delete failed: ' + err.message);
		}
	}

	async function handleSave(e: SubmitEvent) {
		e.preventDefault();
		isSaving = true;
		const method = isEditing ? 'PUT' : 'POST';
		try {
			const payload = { ...currentRecord };
			for (let key in payload) {
				if (typeof payload[key] === 'string' && !isNaN(Number(payload[key])) && payload[key] !== '') {
					if (payload[key].length < 15) payload[key] = Number(payload[key]);
				}
			}

			await api(`/api/explorer?table=${selectedTable}`, {
				method,
				body: JSON.stringify(payload)
			});
			showModal = false;
			tableRef.fetchData();
		} catch (err: any) {
			alert('Save failed: ' + err.message);
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<header class="flex flex-col md:flex-row md:items-center justify-between gap-6">
		<div class="flex items-center gap-4 text-left">
			<div class="p-3 bg-teal-600 text-white rounded shadow shadow-teal-600/20">
				<Database size={24} />
			</div>
			<div>
				<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Database Explorer</h2>
				<p class="text-slate-500 text-sm font-medium">Monitoring system tables including soft-deleted records.</p>
			</div>
		</div>
		
		<div class="flex items-center gap-2">
			<select class="select w-56 bg-white border-slate-200 font-bold text-slate-700" bind:value={selectedTable}>
				{#each tables as table (table.id)}
					<option value={table.id}>{table.name}</option>
				{/each}
			</select>

			<button class="btn bg-teal-600 text-white font-black uppercase tracking-widest text-xs flex items-center gap-2 py-3 px-6 shadow-teal-600/20 hover:opacity-90" onclick={openCreate}>
				<Plus size={18} />
				New Record
			</button>
		</div>
	</header>

	<DataTable 
		bind:this={tableRef}
		endpoint="/api/explorer" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Search system records..."
		extraParams={{ table: selectedTable }}
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as any}
			{#if key === 'actions'}
				<div class="flex items-center justify-center gap-3">
					<button class="text-slate-300 hover:text-teal-600 transition-colors" onclick={() => openEdit(row)} title="Edit">
						<Edit2 size={14} />
					</button>
					{#if !row.deleted_at}
						<button class="text-slate-300 hover:text-red-600 transition-colors" onclick={() => deleteRecord(row.id)} title="Delete">
							<Trash2 size={14} />
						</button>
					{/if}
				</div>
			{:else if key === 'deleted_at'}
				{#if row[key]}
					<span class="badge preset-filled-error-500 border-none px-3 flex items-center gap-1 mx-auto">
						<AlertTriangle size={10} />
						Deleted
					</span>
				{:else}
					<span class="text-[10px] font-black text-slate-200 uppercase tracking-widest">Active</span>
				{/if}
			{:else}
				<span class="text-slate-700 font-medium text-xs truncate block max-w-[250px] {row.deleted_at ? 'opacity-30 grayscale' : ''}" title={String(row[key])}>
					{row[key] === null ? '-' : String(row[key])}
				</span>
			{/if}
		{/snippet}
	</DataTable>
</div>

<!-- Modal -->
{#if showModal}
	<div class="fixed inset-0 bg-slate-900/60 backdrop-blur-md z-50 flex items-center justify-center p-4">
		<div class="card p-8 bg-white border border-slate-200 max-w-2xl w-full space-y-8 max-h-[90vh] overflow-y-auto">
			<div class="flex justify-between items-center border-b border-slate-100 pb-6 text-left">
				<h3 class="h3 font-black uppercase tracking-tight text-slate-900">
					{isEditing ? 'Update Record' : 'Create Record'}
				</h3>
				<button class="text-slate-400 hover:text-slate-600 p-2" onclick={() => showModal = false}>
					<X size={24} />
				</button>
			</div>

			<form onsubmit={handleSave} class="space-y-6 text-left" autocomplete="off">
				<div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
					{#each editableColumns as col (col)}
						<label class="label">
							<span class="text-[10px] font-black uppercase tracking-widest text-slate-400 mb-2 block">{col.replace(/_/g, ' ')}</span>
							{#if typeof currentRecord[col] === 'number'}
								<input class="input font-bold" type="number" bind:value={currentRecord[col]} step="any" required />
							{:else if col.includes('date')}
								<input class="input font-bold" type="date" bind:value={currentRecord[col]} required />
							{:else}
								<input class="input font-bold" type="text" bind:value={currentRecord[col]} required />
							{/if}
						</label>
					{/each}
				</div>

				<div class="flex gap-4 pt-10 sticky bottom-0 bg-white border-t border-slate-50 mt-8">
					<button type="button" class="btn bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200 flex-1 font-black uppercase tracking-widest text-xs py-4" onclick={() => showModal = false}>Cancel</button>
					<button type="submit" class="btn bg-teal-600 text-white flex-1 font-black uppercase tracking-widest text-xs py-4 shadow-teal-600/30" disabled={isSaving}>
						{#if isSaving}
							<RefreshCcw size={18} class="animate-spin mr-2" />
							Processing...
						{:else}
							<Save size={18} class="mr-2" />
							{isEditing ? 'Update' : 'Create'}
						{/if}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
