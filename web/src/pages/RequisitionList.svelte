<script lang="ts">
	import { onMount } from 'svelte';
	import { ShoppingCart, RefreshCcw, CheckCircle, Clock, Check, AlertCircle, User, CreditCard, Wallet, HandCoins } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { auth } from '../lib/auth.svelte';

	interface Requisition {
		id: string;
		name: string;
		category: string;
		description: string;
		amount: number;
		status: 'pending' | 'approved' | 'given';
		type: 'cash' | 'reimburse';
		user_id: string | null;
		user_name: string | null;
		approved_by_id: string | null;
		processed_by_id: string | null;
		created_at: string;
	}

	interface PaginatedResult {
		items: Requisition[];
		total_count: number;
		page: number;
		page_size: number;
	}

	let requisitions = $state<Requisition[]>([]);
	let totalCount = $state(0);
	let currentPage = $state(1);
	let pageSize = $state(10);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let isProcessing = $state<string | null>(null);

	let totalPages = $derived(Math.ceil(totalCount / pageSize));

	async function fetchRequisitions() {
		isLoading = true;
		try {
			const result = await api<PaginatedResult>(`/api/requisitions?page=${currentPage}&page_size=${pageSize}&pending=true`);
			requisitions = result.items || [];
			totalCount = result.total_count || 0;
		} catch (e: any) {
			error = "Gagal memuat data pengajuan. Pastikan Anda memiliki izin yang cukup.";
		} finally {
			isLoading = false;
		}
	}

	function changePage(page: number) {
		if (page >= 1 && (page <= totalPages || totalPages === 0)) {
			currentPage = page;
			fetchRequisitions();
		}
	}

	async function updateStatus(id: string) {
		if (isProcessing) return;
		isProcessing = id;
		try {
			await api('/api/requisitions/approve', {
				method: 'POST',
				body: JSON.stringify({ id })
			});
			await fetchRequisitions();
		} catch (e: any) {
			alert(`Gagal menyetujui: ` + e.message);
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

	onMount(fetchRequisitions);
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 text-slate-900">Persetujuan Pengajuan</h2>
			<p class="text-surface-500 text-sm italic font-medium">Khusus Manager: Periksa dan berikan persetujuan untuk pengajuan staf.</p>
		</div>
		<button onclick={fetchRequisitions} class="btn preset-tonal-surface flex items-center gap-2">
			<RefreshCcw size={16} class={isLoading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	{#if error}
		<div class="card p-12 text-center border border-error-500/30 bg-error-500/5">
			<p class="text-error-500 font-bold mb-2 text-lg">Kesalahan</p>
			<p class="text-surface-500 text-sm">{error}</p>
		</div>
	{:else}
		<!-- Requisitions Table -->
		<section class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
			{#if isLoading && requisitions.length === 0}
				<div class="p-12 text-center text-surface-500">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent mb-4"></div>
					<p>Loading requisitions...</p>
				</div>
			{:else if requisitions.length === 0}
				<div class="p-12 text-center text-slate-400 italic bg-slate-50">
					Tidak ada pengajuan yang menunggu persetujuan.
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full text-left">
						<thead class="bg-slate-50 text-slate-600 font-bold border-b border-slate-200 text-xs">
							<tr>
								<th class="p-4 uppercase tracking-wider">Pengaju / Tipe</th>
								<th class="p-4 uppercase tracking-wider">Kategori & Nama</th>
								<th class="p-4 uppercase tracking-wider">Jumlah</th>
								<th class="p-4 uppercase tracking-wider">Status</th>
								<th class="p-4 uppercase tracking-wider text-right">Tindakan</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each requisitions as req (req.id)}
								<tr class="hover:bg-slate-50 transition-colors">
									<td class="p-4">
										<div class="flex flex-col">
											<span class="text-sm font-bold text-slate-800">{req.user_name || req.user_id || 'Anonim'}</span>
											<span class="text-xs flex items-center gap-1 text-slate-500">
												{#if req.type === 'cash'} <Wallet size={12}/> {:else} <CreditCard size={12}/> {/if}
												{req.type === 'cash' ? 'Tunai' : 'Reimburse'}
											</span>
										</div>
									</td>
									<td class="p-4">
										<div class="flex flex-col">
											<span class="text-xs font-bold text-primary-600 uppercase tracking-tighter">{req.category}</span>
											<span class="text-sm font-medium text-slate-700 line-clamp-1">{req.name}</span>
										</div>
									</td>
									<td class="p-4 font-bold text-slate-900">
										{formatCurrency(req.amount)}
									</td>
									<td class="p-4">
										<span class={getStatusBadge(req.status)}>
											{(req.status || '').toUpperCase()}
										</span>
									</td>
									<td class="p-4 text-right">
										{#if auth.user?.clearance >= 10}
											<button 
												class="btn btn-sm preset-filled-primary-500 font-bold flex items-center gap-2 ml-auto"
												onclick={() => updateStatus(req.id)}
												disabled={isProcessing === req.id}
											>
												{#if isProcessing === req.id}
													<RefreshCcw size={14} class="animate-spin" />
												{:else}
													<Check size={14} />
												{/if}
												Setujui
											</button>
										{:else}
											<span class="text-slate-400 text-xs italic text-right block">Khusus Manager</span>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<!-- Pagination Footer -->
				<footer class="p-4 border-t border-slate-200 bg-slate-50 flex flex-col md:flex-row justify-between items-center gap-4">
					<p class="text-xs text-slate-500 font-medium">
						Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, totalCount)} of {totalCount} requisitions
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
							Page {currentPage} of {totalPages || 1}
						</div>
						<button 
							class="btn btn-sm preset-tonal-surface font-bold" 
							onclick={() => changePage(currentPage + 1)}
							disabled={currentPage === totalPages || totalPages === 0}
						>
							Next
						</button>
					</div>
				</footer>
			{/if}
		</section>
	{/if}
</div>
