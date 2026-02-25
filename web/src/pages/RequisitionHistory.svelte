<script lang="ts">
	import { onMount } from 'svelte';
	import { ShoppingBag, RefreshCcw, CheckCircle, Clock, Search, Wallet, CreditCard, User, HandCoins } from 'lucide-svelte';
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
	let activeTab = $state<'personal' | 'general'>('personal');
	let searchTerm = $state('');
	let isProcessing = $state<string | null>(null);

	let totalPages = $derived(Math.ceil(totalCount / pageSize));

	async function fetchRequisitions() {
		isLoading = true;
		const userIdParam = activeTab === 'personal' ? `&user_id=${auth.user?.id}` : '';
		try {
			const result = await api<PaginatedResult>(`/api/requisitions?page=${currentPage}&page_size=${pageSize}${userIdParam}`);
			requisitions = result.items || [];
			totalCount = result.total_count || 0;
		} catch (e: any) {
			error = "Gagal memuat data pengajuan.";
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

	function switchTab(tab: 'personal' | 'general') {
		activeTab = tab;
		currentPage = 1;
		fetchRequisitions();
	}

	let filteredRequisitions = $derived(
		requisitions.filter(r => 
			r.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
			r.category.toLowerCase().includes(searchTerm.toLowerCase())
		)
	);

	async function markAsGiven(id: string) {
		if (isProcessing) return;
		isProcessing = id;
		try {
			await api('/api/requisitions/give', {
				method: 'POST',
				body: JSON.stringify({ id })
			});
			await fetchRequisitions();
		} catch (e: any) {
			alert(`Gagal mencairkan dana: ` + e.message);
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
			<h2 class="h2 text-slate-900 font-bold uppercase tracking-tight">Daftar Pengajuan</h2>
			<p class="text-surface-500 text-sm italic font-medium">Lihat riwayat dan status seluruh pengajuan Anda maupun perusahaan.</p>
		</div>
		<button onclick={fetchRequisitions} class="btn preset-tonal-surface flex items-center gap-2">
			<RefreshCcw size={16} class={isLoading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	<!-- Search and Tabs -->
	<div class="flex flex-col md:flex-row justify-between items-center gap-4 bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
		<div class="flex p-1 bg-slate-100 rounded-lg w-full md:w-auto">
			<button 
				class="flex-1 md:flex-none px-6 py-2 rounded-md font-bold text-xs transition-all {activeTab === 'personal' ? 'bg-white text-primary-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}"
				onclick={() => switchTab('personal')}
			>
				Pribadi (My Requests)
			</button>
			<button 
				class="flex-1 md:flex-none px-6 py-2 rounded-md font-bold text-xs transition-all {activeTab === 'general' ? 'bg-white text-primary-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}"
				onclick={() => switchTab('general')}
			>
				Semua Pengajuan
			</button>
		</div>

		<div class="relative w-full md:w-64">
			<Search class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
			<input 
				type="text" 
				placeholder="Cari pengajuan..." 
				class="input pl-10 bg-slate-50 border-slate-200 w-full text-xs py-2"
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
		<!-- Requisitions Table -->
		<section class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
			{#if isLoading && requisitions.length === 0}
				<div class="p-12 text-center text-surface-500">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent mb-4"></div>
					<p>Loading...</p>
				</div>
			{:else if filteredRequisitions.length === 0}
				<div class="p-12 text-center text-slate-400 italic bg-slate-50">
					Tidak ada pengajuan ditemukan untuk kriteria ini.
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full text-left">
						<thead class="bg-slate-50 text-slate-600 font-bold border-b border-slate-200 text-xs">
							<tr>
								<th class="p-4 uppercase tracking-wider">Tipe & Kategori</th>
								<th class="p-4 uppercase tracking-wider">Nama Pengajuan</th>
								<th class="p-4 uppercase tracking-wider">Jumlah</th>
								<th class="p-4 uppercase tracking-wider">Status</th>
								<th class="p-4 uppercase tracking-wider">Approval / Given</th>
								<th class="p-4 uppercase tracking-wider text-right">Aksi</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each filteredRequisitions as req (req.id)}
								<tr class="hover:bg-slate-50 transition-colors">
									<td class="p-4">
										<div class="flex flex-col">
											<span class="text-xs flex items-center gap-1 font-bold text-slate-800 uppercase tracking-tighter">
												{#if req.type === 'cash'} <Wallet size={12}/> {:else} <CreditCard size={12}/> {/if}
												{req.type === 'cash' ? 'Tunai' : 'Reimburse'}
											</span>
											<span class="text-xs text-primary-500 font-medium">{req.category}</span>
										</div>
									</td>
									<td class="p-4">
										<div class="flex flex-col">
											<span class="text-sm font-bold text-slate-800 line-clamp-1">{req.name}</span>
											<span class="text-[10px] flex items-center gap-1 text-slate-400">
												<User size={10} /> {req.user_name || req.user_id || 'Anonim'}
											</span>
										</div>
									</td>
									<td class="p-4 font-bold text-slate-900 text-sm whitespace-nowrap">
										{formatCurrency(req.amount)}
									</td>
									<td class="p-4">
										<span class={getStatusBadge(req.status)}>
											{(req.status || '').toUpperCase()}
										</span>
									</td>
									<td class="p-4">
										<div class="flex flex-col gap-1">
											<div class="flex items-center gap-1">
												<span class="text-[9px] uppercase font-bold text-slate-400 w-12">Approve:</span>
												{#if req.approved_by_id}
													<span class="text-[10px] font-mono bg-slate-100 px-1 py-0.5 rounded border border-slate-200" title={req.approved_by_id}>
														{req.approved_by_id.substring(0, 8)}
													</span>
												{:else}
													<span class="text-[10px] italic text-slate-400">-</span>
												{/if}
											</div>
											<div class="flex items-center gap-1">
												<span class="text-[9px] uppercase font-bold text-slate-400 w-12">Given:</span>
												{#if req.processed_by_id}
													<span class="text-[10px] font-mono bg-slate-100 px-1 py-0.5 rounded border border-slate-200" title={req.processed_by_id}>
														{req.processed_by_id.substring(0, 8)}
													</span>
												{:else}
													<span class="text-[10px] italic text-slate-400">-</span>
												{/if}
											</div>
										</div>
									</td>
									<td class="p-4 text-right">
										{#if req.status === 'approved' && auth.user?.clearance >= 0}
											<button 
												class="btn btn-sm preset-filled-success-500 font-bold flex items-center gap-2 ml-auto"
												onclick={() => markAsGiven(req.id)}
												disabled={isProcessing === req.id}
											>
												{#if isProcessing === req.id}
													<RefreshCcw size={14} class="animate-spin" />
												{:else}
													<HandCoins size={14} />
												{/if}
												Cairkan
											</button>
										{:else if req.status === 'given'}
											<CheckCircle size={18} class="text-success-500 ml-auto" />
										{:else}
											<span class="text-slate-300 italic text-xs">Menunggu</span>
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
