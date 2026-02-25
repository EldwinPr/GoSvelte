<script lang="ts">
	import { onMount } from 'svelte';
	import { FileText, RefreshCcw, ArrowUpRight, ArrowDownRight, Search } from 'lucide-svelte';
	import { api } from '../lib/api';

	interface Transaction {
		id: string;
		date: string;
		description: string;
		amount: number;
		type: 'Debit' | 'Credit';
		reference_id: string;
		reference_type: string;
	}

	interface PaginatedResult {
		items: Transaction[];
		total_count: number;
		page: number;
		page_size: number;
	}

	let transactions = $state<Transaction[]>([]);
	let totalCount = $state(0);
	let currentPage = $state(1);
	let pageSize = $state(10);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let searchTerm = $state('');

	let totalPages = $derived(Math.ceil(totalCount / pageSize));

	async function fetchTransactions() {
		isLoading = true;
		try {
			const result = await api<PaginatedResult>(`/api/transactions?page=${currentPage}&page_size=${pageSize}`);
			transactions = result.items || [];
			totalCount = result.total_count || 0;
		} catch (e: any) {
			error = "Gagal memuat transaksi. Silakan coba lagi nanti.";
		} finally {
			isLoading = false;
		}
	}

	function changePage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
			fetchTransactions();
		}
	}

	let filteredTransactions = $derived(
		transactions.filter(t => 
			t.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
			t.reference_type.toLowerCase().includes(searchTerm.toLowerCase())
		)
	);

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
	}

	function formatDate(dateStr: string) {
		return new Date(dateStr).toLocaleDateString('id-ID', {
			day: '2-digit',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	onMount(fetchTransactions);
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 text-slate-900">Company Transactions</h2>
			<p class="text-surface-500 text-sm font-medium italic">Daftar semua transaksi keuangan perusahaan (Audit Trail).</p>
		</div>
		<button onclick={fetchTransactions} class="btn preset-tonal-surface flex items-center gap-2">
			<RefreshCcw size={16} class={isLoading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	<!-- Search and Filters -->
	<div class="card p-4 bg-white border border-slate-200 shadow-sm flex flex-col md:flex-row gap-4 items-center">
		<div class="relative flex-1 w-full">
			<Search class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
			<input 
				type="text" 
				placeholder="Cari deskripsi atau tipe referensi..." 
				class="input pl-10 bg-slate-50 border-slate-200 w-full"
				bind:value={searchTerm}
			/>
		</div>
	</div>

	{#if error}
		<div class="card p-12 text-center border border-error-500/30 bg-error-500/5">
			<p class="text-error-500 font-bold mb-2 text-lg">Kesalahan</p>
			<p class="text-surface-500 text-sm">{error}</p>
		</div>
	{:else}
		<!-- Transactions Table -->
		<section class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
			{#if isLoading && transactions.length === 0}
				<div class="p-12 text-center text-surface-500">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent mb-4"></div>
					<p>Loading transactions...</p>
				</div>
			{:else if filteredTransactions.length === 0}
				<div class="p-12 text-center text-slate-400 italic bg-slate-50">
					Tidak ada transaksi ditemukan.
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full text-left">
						<thead class="bg-slate-50 text-slate-600 font-bold border-b border-slate-200">
							<tr>
								<th class="p-4">Tanggal & Waktu</th>
								<th class="p-4">Deskripsi</th>
								<th class="p-4">Referensi</th>
								<th class="p-4 text-right">Jumlah</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each filteredTransactions as transaction (transaction.id)}
								<tr class="hover:bg-slate-50 transition-colors">
									<td class="p-4 text-sm whitespace-nowrap">
										{formatDate(transaction.date)}
									</td>
									<td class="p-4">
										<div class="flex items-center gap-3">
											<div class="p-2 rounded-lg {transaction.type === 'Credit' ? 'bg-success-500/10 text-success-500' : 'bg-error-500/10 text-error-500'}">
												{#if transaction.type === 'Credit'}
													<ArrowUpRight size={16} />
												{:else}
													<ArrowDownRight size={16} />
												{/if}
											</div>
											<span class="font-medium text-slate-700">{transaction.description}</span>
										</div>
									</td>
									<td class="p-4">
										<span class="badge preset-tonal-surface text-xs uppercase tracking-wider">
											{transaction.reference_type}
										</span>
									</td>
									<td class="p-4 text-right font-bold {transaction.type === 'Credit' ? 'text-success-600' : 'text-error-600'}">
										{transaction.type === 'Credit' ? '+' : '-'}{formatCurrency(transaction.amount)}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<!-- Pagination Footer -->
				<footer class="p-4 border-t border-slate-200 bg-slate-50 flex flex-col md:flex-row justify-between items-center gap-4">
					<p class="text-xs text-slate-500 font-medium">
						Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, totalCount)} of {totalCount} transactions
					</p>
					<div class="flex gap-2">
						<button 
							class="btn btn-sm preset-tonal-surface font-bold" 
							onclick={() => changePage(currentPage - 1)}
							disabled={currentPage === 1}
						>
							Previous
						</button>
						<div class="flex items-center gap-2 px-4 text-xs font-bold text-slate-600">
							Page {currentPage} of {totalPages}
						</div>
						<button 
							class="btn btn-sm preset-tonal-surface font-bold" 
							onclick={() => changePage(currentPage + 1)}
							disabled={currentPage === totalPages}
						>
							Next
						</button>
					</div>
				</footer>
			{/if}
		</section>
	{/if}
</div>
