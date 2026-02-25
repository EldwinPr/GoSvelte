<script lang="ts">
	import { onMount } from 'svelte';
	import { Wallet, ExternalLink, Calendar } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { link } from 'svelte-spa-router';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface Payment {
		id: string;
		invoice_id: string;
		balance_id: string;
		amount: number;
		payment_date: string;
		method: string;
		created_by?: { name: string };
	}

	const columns: Column[] = [
		{ key: 'payment_date', label: 'Waktu Bayar', sortable: true, width: '160px' },
		{ key: 'created_by', label: 'Petugas', width: '140px' },
		{ key: 'method', label: 'Metode', sortable: true, width: '120px' },
		{ key: 'invoice_id', label: 'Target', width: '140px' },
		{ key: 'amount', label: 'Nominal', sortable: true, align: 'right', width: '160px' },
		{ key: 'actions', label: 'Info', align: 'center', width: '80px' }
	];

	let totalDue = $state(0);

	async function fetchSummary() {
		try {
			const result = await api<{ total_due: number }>('/api/payments/summary');
			totalDue = result.total_due || 0;
		} catch (e) {}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return "-";
		return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
	}

	onMount(fetchSummary);
</script>

<div class="space-y-8 text-left">
	<!-- Unified Header Pattern -->
	<header class="flex flex-col md:flex-row md:items-center justify-between gap-6">
		<div class="flex items-center gap-4">
			<div class="p-3 bg-teal-600 text-white shadow shadow-teal-600/20">
				<Wallet size={24} />
			</div>
			<div>
				<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Recent Payments</h2>
				<p class="text-slate-500 text-sm font-medium italic opacity-70">Monitoring real-time penerimaan dana perusahaan.</p>
			</div>
		</div>
	</header>

	<!-- Summary Card - Flattened -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="card p-8 bg-teal-600 text-white border-none flex flex-col gap-2">
			<div class="flex justify-between items-center opacity-80 uppercase text-[10px] font-black tracking-widest">
				<span>Piutang Outstanding</span>
				<Wallet size={16} />
			</div>
			<h3 class="text-3xl font-black">{formatCurrency(totalDue)}</h3>
			<p class="text-[10px] italic opacity-70">Total tagihan jatuh tempo yang belum terbayar.</p>
		</div>
	</div>

	<DataTable 
		endpoint="/api/payments" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari ID Pembayaran atau Invoice..."
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as Payment}
			{#if key === 'payment_date'}
				<div class="flex items-center gap-2 text-[11px] font-bold text-slate-400">
					<Calendar size={12} />
					{formatDate(row.payment_date)}
				</div>
			{:else if key === 'created_by'}
				<span class="text-xs font-black text-slate-700 uppercase">{row.created_by?.name || 'System'}</span>
			{:else if key === 'method'}
				<span class="badge preset-tonal-surface text-[9px] font-black uppercase px-3 py-1">
					{row.method}
				</span>
			{:else if key === 'invoice_id'}
				<span class="font-mono text-xs text-teal-600 font-black">#{row.invoice_id.substring(0, 8)}</span>
			{:else if key === 'amount'}
				<span class="font-black text-green-600">{formatCurrency(row.amount)}</span>
			{:else if key === 'actions'}
				<a use:link href="/invoices/detail?id={row.invoice_id}" class="text-slate-300 hover:text-teal-600 transition-all mx-auto block w-fit">
					<ExternalLink size={16} />
				</a>
			{:else}
				{row[key as keyof Payment]}
			{/if}
		{/snippet}
	</DataTable>
</div>
