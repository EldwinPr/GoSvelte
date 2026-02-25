<script lang="ts">
	import { onMount } from 'svelte';
	import { CheckCircle, Wallet, CreditCard, User, HandCoins, Eye, Clock, FileText, Calendar } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { auth } from '../lib/auth.svelte';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';
	import { link } from 'svelte-spa-router';

	interface Requisition {
		id: string;
		name: string;
		category: string;
		description: string;
		amount: number;
		status: 'pending' | 'approved' | 'given' | 'rejected';
		type: 'cash' | 'reimburse';
		user_id: string | null;
		user?: { name: string };
		user_name: string | null;
		approved_by_id: string | null;
		approved_by?: { name: string };
		processed_by_id: string | null;
		processed_by?: { name: string };
		created_at: string;
	}

	const columns: Column[] = [
		{ key: 'type', label: 'Tipe & Kategori', sortable: true, width: '160px' },
		{ key: 'name', label: 'Nama Pengajuan', sortable: true },
		{ key: 'amount', label: 'Jumlah', sortable: true, align: 'right', width: '160px' },
		{ key: 'status', label: 'Status', sortable: true, align: 'center', width: '120px' },
		{ key: 'workflow', label: 'Approval / Given', width: '160px' },
		{ key: 'actions', label: 'Aksi', align: 'center', width: '100px' }
	];

	let activeTab = $state<'personal' | 'general'>('personal');
	let tableRef: any = $state();

	function switchTab(tab: 'personal' | 'general') {
		activeTab = tab;
		setTimeout(() => tableRef.fetchData(), 0);
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return "-";
		return new Date(dateStr).toLocaleDateString('id-ID', {
			day: '2-digit',
			month: 'short',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusBadge(status: string) {
		const s = (status || '').toLowerCase();
		switch (s) {
			case 'given': return 'badge preset-filled-success-500';
			case 'approved': return 'badge preset-filled-primary-500';
			case 'pending': return 'badge preset-filled-warning-500';
			case 'rejected': return 'badge preset-filled-error-500';
			default: return 'badge preset-tonal-surface';
		}
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 text-left">
		<div>
			<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Daftar Pengajuan</h2>
			<p class="text-slate-500 text-sm font-medium italic opacity-70">History and status of requisitions.</p>
		</div>
	</div>

	<!-- Tabs -->
	<div class="flex p-1 bg-slate-100 rounded-lg w-fit">
		<button class="px-6 py-2 rounded-md font-bold text-xs transition-all {activeTab === 'personal' ? 'bg-white text-teal-600 shadow-sm' : 'text-slate-500'}" onclick={() => switchTab('personal')}>
			Personal
		</button>
		<button class="px-6 py-2 rounded-md font-bold text-xs transition-all {activeTab === 'general' ? 'bg-white text-teal-600 shadow-sm' : 'text-slate-500'}" onclick={() => switchTab('general')}>
			General
		</button>
	</div>

	<DataTable 
		bind:this={tableRef}
		endpoint="/api/requisitions" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari pengajuan..."
		extraParams={activeTab === 'personal' ? { user_id: auth.user?.id || '' } : {}}
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as Requisition}
			{#if key === 'type'}
				<div class="flex flex-col text-left">
					<span class="text-xs flex items-center gap-1 font-bold text-slate-800 uppercase tracking-tighter">
						{#if row.type === 'cash'} <Wallet size={12}/> {:else} <CreditCard size={12}/> {/if}
						{row.type === 'cash' ? 'Tunai' : 'Reimburse'}
					</span>
					<span class="text-xs text-teal-600 font-medium">{row.category}</span>
				</div>
			{:else if key === 'name'}
				<div class="flex flex-col text-left">
					<span class="text-sm font-bold text-slate-800 line-clamp-1">{row.name}</span>
					<span class="text-[10px] flex items-center gap-1 text-slate-400 font-black uppercase">
						<User size={10} /> {row.user?.name || row.user_name || 'Anonim'}
					</span>
				</div>
			{:else if key === 'amount'}
				<span class="font-bold text-slate-900">{formatCurrency(row.amount)}</span>
			{:else if key === 'status'}
				<span class={getStatusBadge(row.status)}>{row.status.toUpperCase()}</span>
			{:else if key === 'workflow'}
				<div class="flex flex-col gap-1 text-left">
					<div class="flex items-center gap-1">
						<span class="text-[9px] uppercase font-black text-slate-400 w-12">Approve:</span>
						{#if row.approved_by?.name}
							<span class="text-[10px] font-bold text-slate-700">{row.approved_by.name}</span>
						{:else if row.approved_by_id}
							<span class="text-[10px] font-mono bg-slate-100 px-1 py-0.5 rounded border border-slate-200">{row.approved_by_id.substring(0, 8)}</span>
						{:else}
							<span class="text-[10px] italic text-slate-400">-</span>
						{/if}
					</div>
					<div class="flex items-center gap-1">
						<span class="text-[9px] uppercase font-black text-slate-400 w-12">Given:</span>
						{#if row.processed_by?.name}
							<span class="text-[10px] font-bold text-slate-700">{row.processed_by.name}</span>
						{:else if row.processed_by_id}
							<span class="text-[10px] font-mono bg-slate-100 px-1 py-0.5 rounded border border-slate-200">{row.processed_by_id.substring(0, 8)}</span>
						{:else}
							<span class="text-[10px] italic text-slate-400">-</span>
						{/if}
					</div>
				</div>
			{:else if key === 'actions'}
				<div class="flex justify-center items-center gap-2">
					<a 
						use:link 
						href="/purchasing/detail?id={row.id}" 
						class="btn btn-sm btn-icon bg-slate-50 text-slate-400 hover:text-teal-600 transition-colors" 
						title="View Detail"
					>
						<Eye size={18} />
					</a>
				</div>
			{/if}
		{/snippet}
	</DataTable>
</div>
