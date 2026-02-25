<script lang="ts">
	import { ExternalLink, Calendar, Receipt } from 'lucide-svelte';
	import { link } from 'svelte-spa-router';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface Transaction {
		id: string;
		date: string;
		description: string;
		category: string;
		amount: number;
		type: 'Debit' | 'Credit';
		reference_id: string;
		reference_type: string;
	}

	const columns: Column[] = [
		{ key: 'date', label: 'Waktu', sortable: true, width: '160px' },
		{ key: 'category', label: 'Kategori', sortable: true, width: '120px' },
		{ key: 'description', label: 'Deskripsi Transaksi', sortable: true },
		{ key: 'reference_type', label: 'Referensi', align: 'center', width: '140px' },
		{ key: 'amount', label: 'Jumlah', sortable: true, align: 'right', width: '180px' },
	];

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

	function getCategoryBadge(category: string) {
		const c = (category || '').toLowerCase();
		if (c === 'income') return 'badge preset-filled-success-500 border-none';
		if (c === 'expense') return 'badge preset-filled-error-500 border-none';
		return 'badge preset-tonal-surface';
	}
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 text-left">
		<div class="flex items-center gap-4">
			<div class="p-3 bg-teal-600 text-white rounded shadow-sm shadow-teal-600/20">
				<Receipt size={24} />
			</div>
			<div>
				<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Company Transactions</h2>
				<p class="text-slate-500 text-sm font-medium italic opacity-70">Log riwayat transaksi keuangan perusahaan secara real-time.</p>
			</div>
		</div>
	</div>

	<DataTable 
		endpoint="/api/transactions" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari deskripsi, kategori, atau tipe..."
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as Transaction}
			{#if key === 'date'}
				<div class="flex items-center gap-2 text-[11px] font-bold text-slate-400">
					<Calendar size={12} />
					{formatDate(row.date)}
				</div>
			{:else if key === 'category'}
				<span class="{getCategoryBadge(row.category)} text-[9px] font-black uppercase px-3 py-1">
					{row.category || 'Other'}
				</span>
			{:else if key === 'description'}
				<p class="font-bold text-slate-700 text-sm truncate max-w-md" title={row.description}>
					{row.description}
				</p>
			{:else if key === 'reference_type'}
				<div class="flex items-center justify-center gap-2">
					<span class="text-[9px] font-black text-slate-400 uppercase tracking-widest bg-slate-50 px-2 py-1 border border-slate-100">
						{row.reference_type}
					</span>
					{#if row.reference_type === 'Requisition'}
						<a use:link href="/purchasing/detail?id={row.reference_id}" class="text-teal-600 hover:scale-110 transition-transform" title="Lihat Detail Pengajuan">
							<ExternalLink size={14} />
						</a>
					{:else if row.reference_type === 'InvoicePayment'}
						<a use:link href="/invoices/detail?id={row.reference_id}" class="text-teal-600 hover:scale-110 transition-transform" title="Lihat Detail Invoice">
							<ExternalLink size={14} />
						</a>
					{/if}
				</div>
			{:else if key === 'amount'}
				<div class="flex flex-col items-end">
					<span class="font-black text-sm {row.type === 'Credit' ? 'text-green-600' : 'text-red-600'}">
						{row.type === 'Credit' ? '+' : '-'}{formatCurrency(row.amount)}
					</span>
					<span class="text-[9px] font-black uppercase tracking-tighter opacity-30">{row.type}</span>
				</div>
			{:else}
				{row[key as keyof Transaction]}
			{/if}
		{/snippet}
	</DataTable>
</div>
