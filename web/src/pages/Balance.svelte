<script lang="ts">
	import { onMount } from 'svelte';
	import { Landmark, TrendingUp, ArrowUpRight, ArrowDownRight, RefreshCcw } from 'lucide-svelte';
	import { api } from '../lib/api';
	import DataTable from '../lib/DataTable.svelte';
	import type { Column } from '../lib/types';

	interface CompanyBalance {
		id: string;
		account_name: string;
		balance: number;
		updated_at: string;
	}

	let totalBalance = $state(0);
	let isLoadingSummary = $state(true);
	let tableRef: any = $state();

	const columns: Column[] = [
		{ key: 'account_name', label: 'Account Name', sortable: true },
		{ key: 'balance', label: 'Current Balance', sortable: true, align: 'right', width: '200px' }
	];

	async function fetchSummary() {
		isLoadingSummary = true;
		try {
			const result = await api<any>('/api/balances?page_size=100');
			const items = result.items || [];
			totalBalance = items.reduce((acc: number, curr: any) => acc + curr.balance, 0);
		} catch (e) {
		} finally {
			isLoadingSummary = false;
		}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	onMount(fetchSummary);
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4 text-left">
		<div>
			<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Company Balance</h2>
			<p class="text-slate-500 text-sm font-medium italic opacity-70">Ringkasan real-time dari seluruh aset likuid dan akun perusahaan.</p>
		</div>
		<button onclick={() => { fetchSummary(); tableRef.fetchData(); }} class="btn bg-slate-100 border border-slate-200 text-slate-600 flex items-center gap-2 font-black uppercase text-xs">
			<RefreshCcw size={16} class={isLoadingSummary ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	<!-- Stats Grid -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<div class="card p-8 bg-white border border-slate-200 space-y-3">
			<div class="flex justify-between items-start">
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Aset Likuid</p>
				<div class="p-2 bg-teal-50 rounded text-teal-600">
					<Landmark size={24} />
				</div>
			</div>
			<h3 class="h2 font-black text-slate-900">{formatCurrency(totalBalance)}</h3>
			<div class="flex items-center gap-1 text-xs text-green-600 font-bold">
				<TrendingUp size={14} />
				<span>+2.4% vs last month</span>
			</div>
		</div>

		<div class="card p-8 bg-white border border-slate-200 space-y-3">
			<div class="flex justify-between items-start">
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Kas Operasional</p>
				<div class="p-2 bg-green-50 rounded text-green-600">
					<ArrowUpRight size={24} />
				</div>
			</div>
			<h3 class="h2 font-black text-slate-900">{formatCurrency(totalBalance * 0.7)}</h3>
			<p class="text-[10px] text-slate-400 font-bold uppercase tracking-tight italic">Estimated daily budget</p>
		</div>

		<div class="card p-8 bg-white border border-slate-200 space-y-3">
			<div class="flex justify-between items-start">
				<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Dana Cadangan</p>
				<div class="p-2 bg-amber-50 rounded text-amber-600">
					<ArrowDownRight size={24} />
				</div>
			</div>
			<h2 class="h2 font-black text-slate-900">{formatCurrency(totalBalance * 0.3)}</h2>
			<p class="text-[10px] text-slate-400 font-bold uppercase tracking-tight italic">Strategic backup funds</p>
		</div>
	</div>

	<DataTable 
		bind:this={tableRef}
		endpoint="/api/balances" 
		{columns} 
		rowKey="id" 
		searchPlaceholder="Cari nama akun..."
	>
		{#snippet cell(rowData, key)}
			{@const row = rowData as CompanyBalance}
			{#if key === 'account_name'}
				<div class="flex items-center gap-3">
					<div class="w-8 h-8 rounded bg-slate-50 flex items-center justify-center text-teal-600 border border-slate-100">
						<Landmark size={16} />
					</div>
					<span class="font-bold text-slate-800 text-left">{row.account_name}</span>
				</div>
			{:else if key === 'balance'}
				<span class="font-black text-teal-600">{formatCurrency(row.balance)}</span>
			{:else}
				{row[key as keyof CompanyBalance]}
			{/if}
		{/snippet}
	</DataTable>
</div>
