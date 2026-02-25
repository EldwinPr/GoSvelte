<script lang="ts">
	import { Receipt, Eye, Plus } from 'lucide-svelte';
	import { link } from 'svelte-spa-router';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface Invoice {
		id: string;
		number: string;
		invoice_date: string;
		due_date: string;
		total_amount: number;
		status: string;
		customer?: { customer_name: string };
		payments?: { amount: number }[];
	}

	const columns: Column[] = [
		{ key: 'number', label: 'Nomor', sortable: true, width: '140px' },
		{ key: 'customer', label: 'Pelanggan', sortable: true },
		{ key: 'invoice_date', label: 'Tanggal', sortable: true, align: 'center', width: '110px' },
		{ key: 'paid_pc', label: 'Dibayar', align: 'center', width: '80px' },
		{ key: 'status', label: 'Status', sortable: true, align: 'center', width: '110px' },
		{ key: 'actions', label: 'Aksi', align: 'center', width: '80px' },
		{ key: 'total_amount', label: 'Total', sortable: true, align: 'right', width: '160px' },
	];

	function calculatePaidPercentage(invoice: Invoice) {
		const totalPaid = invoice.payments?.reduce((acc, curr) => acc + curr.amount, 0) || 0;
		if (invoice.total_amount === 0) return 0;
		return Math.round((totalPaid / invoice.total_amount) * 100);
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return "-";
		return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' });
	}

	function getStatusBadge(status: string) {
		const s = (status || '').toLowerCase();
		switch (s) {
			case 'paid': return 'badge preset-filled-success-500';
			case 'finalized': return 'badge bg-teal-600 border-teal-600 text-white';
			case 'overdue': return 'badge preset-filled-error-500';
			case 'sent': return 'badge bg-teal-600 border-teal-600 text-white';
			default: return 'badge preset-tonal-surface';
		}
	}
</script>

<div class="space-y-8 text-left">
	<!-- Unified Header Pattern -->
	<header class="flex flex-col md:flex-row md:items-center justify-between gap-6">
		<div class="flex items-center gap-4">
			<div class="p-3 bg-teal-600 text-white shadow shadow-teal-600/20">
				<Receipt size={24} />
			</div>
			<div>
				<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Billing Invoices</h2>
				<p class="text-slate-500 text-sm font-medium italic opacity-70">Kelola penagihan pelanggan dan status pembayaran.</p>
			</div>
		</div>
		<a use:link href="/invoices/new" class="btn bg-teal-600 text-white font-black uppercase tracking-widest text-[10px] py-3 px-6 shadow-md shadow-teal-600/20 hover:opacity-90">
			<Plus size={16} class="mr-2" /> Buat Invoice
		</a>
	</header>

	<DataTable 
		endpoint="/api/invoices" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari nomor invoice atau pelanggan..."
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as Invoice}
			{#if key === 'number'}
				<span class="font-black text-teal-600 truncate">{row.number}</span>
			{:else if key === 'customer'}
				<p class="font-bold text-slate-700 truncate">{row.customer?.customer_name || 'Pelanggan Umum'}</p>
			{:else if key === 'invoice_date'}
				<span class="text-xs font-bold text-slate-400">{formatDate(row.invoice_date)}</span>
			{:else if key === 'paid_pc'}
				{@const pc = calculatePaidPercentage(row)}
				<div class="flex flex-col items-center gap-1">
					<span class="text-[9px] font-black text-slate-400">{pc}%</span>
					<div class="w-10 bg-slate-100 h-1 overflow-hidden border border-slate-200">
						<div class="bg-teal-600 h-full" style="width: {pc}%"></div>
					</div>
				</div>
			{:else if key === 'status'}
				<span class={getStatusBadge(row.status)}>
					{(row.status || '').toUpperCase()}
				</span>
			{:else if key === 'actions'}
				<div class="flex justify-center">
					<a use:link href="/invoices/detail?id={row.id}" class="text-slate-300 hover:text-teal-600 transition-colors" title="View Detail">
						<Eye size={18} />
					</a>
				</div>
			{:else if key === 'total_amount'}
				<span class="font-black text-slate-900">{formatCurrency(row.total_amount)}</span>
			{:else}
				{row[key as keyof Invoice]}
			{/if}
		{/snippet}
	</DataTable>
</div>
