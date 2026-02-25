<script lang="ts">
	import { RefreshCcw, Check, User, CreditCard, Wallet } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { auth } from '../lib/auth.svelte';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface Requisition {
		id: string;
		name: string;
		category: string;
		description: string;
		amount: number;
		status: 'pending' | 'approved' | 'given';
		type: 'cash' | 'reimburse';
		user_id: string | null;
		user?: { name: string };
		user_name: string | null;
		approved_by_id: string | null;
		processed_by_id: string | null;
		created_at: string;
	}

	const columns: Column[] = [
		{ key: 'user_id', label: 'Pengaju / Tipe', sortable: true, width: '160px' },
		{ key: 'category', label: 'Kategori & Nama', sortable: true },
		{ key: 'amount', label: 'Jumlah', sortable: true, align: 'right', width: '160px' },
		{ key: 'status', label: 'Status', sortable: true, align: 'center', width: '120px' },
		{ key: 'actions', label: 'Tindakan', align: 'right', width: '120px' }
	];

	let isProcessing = $state<string | null>(null);
	let tableRef: any = $state();

	async function updateStatus(id: string, action: 'approve' | 'reject') {
		if (isProcessing) return;
		isProcessing = id;
		const endpoint = action === 'approve' ? '/api/requisitions/approve' : '/api/requisitions/reject';
		try {
			await api(endpoint, {
				method: 'POST',
				body: JSON.stringify({ id })
			});
			tableRef.fetchData();
		} catch (e: any) {
			alert(`Gagal memproses: ` + e.message);
		} finally {
			isProcessing = null;
		}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	function getStatusBadge(status: string) {
		const s = (status || '').toLowerCase();
		switch (s) {
			case 'given': return 'badge preset-filled-success-500';
			case 'approved': return 'badge preset-filled-primary-500';
			case 'pending': return 'badge preset-filled-warning-500';
			default: return 'badge preset-tonal-surface';
		}
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 text-left">
		<div>
			<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Persetujuan Pengajuan</h2>
			<p class="text-slate-500 text-sm font-medium italic opacity-70">Khusus Manager: Periksa dan berikan persetujuan untuk pengajuan staf.</p>
		</div>
	</div>

	<DataTable 
		bind:this={tableRef}
		endpoint="/api/requisitions" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari pengajuan..."
		extraParams={{ pending: true }}
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as Requisition}
			{#if key === 'user_id'}
				<div class="flex flex-col text-left">
					<span class="text-sm font-bold text-slate-800">{row.user?.name || row.user_name || 'Anonim'}</span>
					<span class="text-xs flex items-center gap-1 text-slate-500 font-medium">
						{#if row.type === 'cash'} <Wallet size={12}/> {:else} <CreditCard size={12}/> {/if}
						{row.type === 'cash' ? 'Tunai' : 'Reimburse'}
					</span>
				</div>
			{:else if key === 'category'}
				<div class="flex flex-col text-left">
					<span class="text-xs font-bold text-teal-600 uppercase tracking-tighter">{row.category}</span>
					<span class="text-sm font-medium text-slate-700 line-clamp-1">{row.name}</span>
				</div>
			{:else if key === 'amount'}
				<span class="font-bold text-slate-900">{formatCurrency(row.amount)}</span>
			{:else if key === 'status'}
				<span class={getStatusBadge(row.status)}>{row.status.toUpperCase()}</span>
			{:else if key === 'actions'}
				{#if (auth.user?.clearance ?? 0) >= 10}
					<div class="flex justify-end gap-2">
						<button 
							class="btn btn-xs preset-tonal-error font-black uppercase tracking-widest"
							onclick={() => updateStatus(row.id, 'reject')}
							disabled={isProcessing === row.id}
						>
							Tolak
						</button>
						<button 
							class="btn btn-xs bg-teal-600 text-white font-black uppercase tracking-widest flex items-center gap-2 shadow-md shadow-teal-600/20"
							onclick={() => updateStatus(row.id, 'approve')}
							disabled={isProcessing === row.id}
						>
							{#if isProcessing === row.id}
								<RefreshCcw size={12} class="animate-spin" />
							{:else}
								<Check size={12} />
							{/if}
							Setujui
						</button>
					</div>
				{:else}
					<span class="text-slate-400 text-xs italic text-right block font-bold">KHUSUS MANAGER</span>
				{/if}
			{/if}
		{/snippet}
	</DataTable>
</div>
