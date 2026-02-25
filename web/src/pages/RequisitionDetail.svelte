<script lang="ts">
	import { onMount } from 'svelte';
	import { RefreshCcw, ArrowLeft, CheckCircle, FileText, User, Clock, HandCoins, Landmark, X, Wallet, CreditCard } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { link, push } from 'svelte-spa-router';
	import { auth } from '../lib/auth.svelte';

	let { params } = $props<{ params: { id?: string } }>();
	const id = $derived(params?.id || new URLSearchParams(window.location.hash.split('?')[1]).get('id'));

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

	interface CompanyBalance {
		id: string;
		account_name: string;
		balance: number;
	}

	let requisition = $state<Requisition | null>(null);
	let isLoading = $state(true);
	let isProcessing = $state(false);
	let isDisbursing = $state(false);
	let error = $state<string | null>(null);

	let balances = $state<CompanyBalance[]>([]);
	let selectedBalanceId = $state('');

	async function fetchRequisition() {
		if (!id) return;
		console.log("[ReqDetail] Fetching:", id);
		isLoading = true;
		error = null;
		try {
			requisition = await api<Requisition>(`/api/requisitions/show?id=${id}`);
		} catch (e: any) {
			console.error("[ReqDetail] Fetch Error:", e);
			error = "Pengajuan tidak ditemukan atau gagal dimuat.";
		} finally {
			isLoading = false;
		}
	}

	// REACTIVITY: Fetch data whenever the ID changes
	$effect(() => {
		if (id) fetchRequisition();
	});

	async function fetchBalances() {
		try {
			const result = await api<any>('/api/balances?page_size=100');
			balances = result.items || [];
			if (balances.length > 0) selectedBalanceId = balances[0].id;
		} catch (e) {}
	}

	async function markAsGiven() {
		if (!requisition || !selectedBalanceId || isProcessing) return;
		isProcessing = true;
		try {
			await api('/api/requisitions/give', {
				method: 'POST',
				body: JSON.stringify({ id: requisition.id, balance_id: selectedBalanceId })
			});
			await fetchRequisition();
			isDisbursing = false;
			alert("Dana berhasil dicairkan.");
		} catch (e: any) {
			alert(`Gagal mencairkan dana: ` + e.message);
		} finally {
			isProcessing = false;
		}
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

<div class="max-w-4xl mx-auto space-y-8">
	<div class="flex items-center gap-4 text-left">
		<a use:link href="/purchasing/list" class="btn btn-sm btn-icon bg-white border border-slate-200 text-slate-600 hover:bg-slate-100">
			<ArrowLeft size={18} />
		</a>
		<div>
			<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Detail Pengajuan</h2>
			{#if requisition}
				<p class="text-xs font-mono font-bold text-teal-600 uppercase">#{requisition.id.substring(0, 12)}... — {requisition.status}</p>
			{/if}
		</div>
	</div>

	{#if isLoading}
		<div class="p-20 text-center">
			<RefreshCcw size={32} class="animate-spin mx-auto mb-4 text-slate-200" />
			<p class="font-black uppercase text-xs tracking-widest text-slate-400">Loading Data...</p>
		</div>
	{:else if error || !requisition}
		<div class="card p-12 text-center border border-red-200 bg-red-50">
			<p class="text-red-600 font-black uppercase tracking-widest text-xs mb-2">Error</p>
			<p class="text-slate-600 font-medium">{error || 'Requisition not found'}</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<div class="lg:col-span-2 space-y-6 text-left animate-in fade-in duration-300">
				<!-- Core Info Card -->
				<section class="card bg-white border border-slate-200 p-8 space-y-6">
					<div>
						<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Nama Pengajuan</p>
						<h2 class="text-2xl font-black text-slate-800">{requisition.name}</h2>
					</div>

					<div class="grid grid-cols-2 gap-6">
						<div>
							<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Kategori</p>
							<p class="font-bold text-teal-600 uppercase tracking-tighter text-sm">{requisition.category}</p>
						</div>
						<div>
							<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Metode</p>
							<div class="flex items-center gap-2 font-bold text-slate-700 uppercase text-xs">
								{#if requisition.type === 'cash'} <Wallet size={14}/> {:else} <CreditCard size={14}/> {/if}
								{requisition.type}
							</div>
						</div>
					</div>

					<div class="pt-6 border-t border-slate-100">
						<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Deskripsi Detail</p>
						<p class="text-slate-600 leading-relaxed font-medium bg-slate-50 p-4 border border-slate-100 italic">
							"{requisition.description || 'Tidak ada deskripsi tambahan.'}"
						</p>
					</div>
				</section>

				<!-- Workflow Info -->
				<section class="card bg-white border border-slate-200 overflow-hidden">
					<header class="p-4 bg-slate-50 border-b border-slate-200">
						<h4 class="font-black text-slate-400 text-[10px] uppercase tracking-[0.2em]">Audit & Workflow Trail</h4>
					</header>
					<div class="p-6 space-y-6">
						<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-slate-100 flex items-center justify-center text-slate-400">
									<User size={20} />
								</div>
								<div>
									<p class="text-[9px] font-black text-slate-400 uppercase">Diajukan Oleh</p>
									<p class="text-sm font-bold text-slate-700">{requisition.user?.name || requisition.user_name || 'Anonim'}</p>
								</div>
							</div>
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-slate-100 flex items-center justify-center text-slate-400">
									<Clock size={20} />
								</div>
								<div>
									<p class="text-[9px] font-black text-slate-400 uppercase">Waktu Submit</p>
									<p class="text-sm font-bold text-slate-700">{formatDate(requisition.created_at)}</p>
								</div>
							</div>
						</div>

						<div class="grid grid-cols-1 md:grid-cols-2 gap-6 pt-6 border-t border-slate-50">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-teal-50 flex items-center justify-center text-teal-600">
									<CheckCircle size={20} />
								</div>
								<div>
									<p class="text-[9px] font-black text-slate-400 uppercase">Disetujui Oleh</p>
									<p class="text-sm font-bold text-slate-700">{requisition.approved_by?.name || '-'}</p>
								</div>
							</div>
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 bg-green-50 flex items-center justify-center text-green-600">
									<HandCoins size={20} />
								</div>
								<div>
									<p class="text-[9px] font-black text-slate-400 uppercase">Dicairkan Oleh</p>
									<p class="text-sm font-bold text-slate-700">{requisition.processed_by?.name || '-'}</p>
								</div>
							</div>
						</div>
					</div>
				</section>
			</div>

			<div class="space-y-6 text-left">
				<!-- Status & Amount Card -->
				<section class="card p-8 bg-white border border-slate-200 space-y-6 shadow-sm">
					<div>
						<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-2">Status Current</p>
						<span class={getStatusBadge(requisition.status)}>{requisition.status.toUpperCase()}</span>
					</div>
					
					<div class="pt-4 border-t border-slate-100">
						<p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Dana</p>
						<p class="text-3xl font-black text-teal-600">{formatCurrency(requisition.amount)}</p>
					</div>

					{#if requisition.status === 'approved' && (auth.user?.clearance ?? 0) >= 0 && !isDisbursing}
						<button class="btn bg-teal-600 text-white w-full font-black uppercase tracking-widest text-xs py-4 shadow-md shadow-teal-600/20" onclick={() => { isDisbursing = true; fetchBalances(); }}>
							Cairkan Sekarang
						</button>
					{/if}
				</section>

				{#if isDisbursing}
					<section class="card p-8 bg-slate-900 text-white space-y-6 animate-in slide-in-from-top-4">
						<div class="flex justify-between items-center">
							<h4 class="font-black text-xs uppercase tracking-widest">Pencairan Dana</h4>
							<button onclick={() => isDisbursing = false}><X size={18}/></button>
						</div>
						
						<label class="label text-left">
							<span class="text-[10px] font-black uppercase text-slate-400 mb-2 block">Pilih Akun Sumber</span>
							<select class="select bg-slate-800 border-slate-700 text-white font-bold" bind:value={selectedBalanceId}>
								{#each balances as balance (balance.id)}
									<option value={balance.id}>{balance.account_name} — {formatCurrency(balance.balance)}</option>
								{/each}
							</select>
						</label>

						<button class="btn bg-green-600 text-white w-full font-black uppercase tracking-widest text-xs py-4 shadow-lg" onclick={markAsGiven} disabled={isProcessing}>
							{isProcessing ? 'Processing...' : 'Konfirmasi Cairkan'}
						</button>
					</section>
				{/if}
			</div>
		</div>
	{/if}
</div>
