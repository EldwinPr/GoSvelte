<script lang="ts">
	import { onMount } from 'svelte';
	import { Landmark, TrendingUp, ArrowUpRight, ArrowDownRight, RefreshCcw } from 'lucide-svelte';
	import { api } from '../lib/api';

	interface CompanyBalance {
		id: string;
		account_name: string;
		balance: number;
		updated_at: string;
	}

	let balances = $state<CompanyBalance[]>([]);
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	let totalBalance = $derived(balances.reduce((acc, curr) => acc + curr.balance, 0));

	async function fetchBalances() {
		isLoading = true;
		try {
			balances = await api('/api/balances');
		} catch (e: any) {
			error = "Akses Ditolak: Anda memerlukan izin Manajer (10+) untuk melihat saldo keuangan.";
		} finally {
			isLoading = false;
		}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
	}

	onMount(fetchBalances);
</script>

<div class="space-y-8">
	<!-- Header -->
	<div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
		<div>
			<h2 class="h2 text-slate-900">Company Balance</h2>
			<p class="text-surface-500 text-sm">Ringkasan real-time dari seluruh aset likuid dan akun perusahaan.</p>
		</div>
		<button onclick={fetchBalances} class="btn preset-tonal-surface flex items-center gap-2">
			<RefreshCcw size={16} class={isLoading ? 'animate-spin' : ''} />
			Refresh
		</button>
	</div>

	{#if error}
		<div class="card p-12 text-center border border-error-500/30 bg-error-500/5">
			<p class="text-error-500 font-bold mb-2 text-lg">Kunci Privasi Keuangan</p>
			<p class="text-surface-500 text-sm">{error}</p>
		</div>
	{:else}
		<!-- Stats Grid -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="card p-6 bg-surface-100-900 border border-surface-200-800 space-y-2">
				<div class="flex justify-between items-start">
					<p class="text-sm font-medium text-surface-500">Total Aset Likuid</p>
					<div class="p-2 bg-primary-500/10 rounded-lg text-primary-500">
						<Landmark size={20} />
					</div>
				</div>
				<h3 class="h3 font-bold">{formatCurrency(totalBalance)}</h3>
				<div class="flex items-center gap-1 text-xs text-success-500">
					<TrendingUp size={12} />
					<span>+2.4% dari bulan lalu</span>
				</div>
			</div>

			<div class="card p-6 bg-surface-100-900 border border-surface-200-800 space-y-2">
				<div class="flex justify-between items-start">
					<p class="text-sm font-medium text-surface-500">Kas Operasional</p>
					<div class="p-2 bg-success-500/10 rounded-lg text-success-500">
						<ArrowUpRight size={20} />
					</div>
				</div>
				<h3 class="h3 font-bold">{formatCurrency(totalBalance * 0.7)}</h3>
				<p class="text-xs text-surface-500">Dialokasikan untuk operasional harian</p>
			</div>

			<div class="card p-6 bg-surface-100-900 border border-surface-200-800 space-y-2">
				<div class="flex justify-between items-start">
					<p class="text-sm font-medium text-surface-500">Dana Cadangan</p>
					<div class="p-2 bg-warning-500/10 rounded-lg text-warning-500">
						<ArrowDownRight size={20} />
					</div>
				</div>
				<h3 class="h3 font-bold">{formatCurrency(totalBalance * 0.3)}</h3>
				<p class="text-xs text-surface-500">Cadangan backup strategis</p>
			</div>
		</div>

		<!-- Accounts Table -->
		<section class="card bg-surface-100-900 border border-surface-200-800 overflow-hidden">
			<header class="p-4 border-b border-surface-200-800 bg-surface-200-800/20">
				<h4 class="h4 font-semibold text-slate-800">Financial Accounts</h4>
			</header>
			
			{#if isLoading && balances.length === 0}
				<div class="p-12 text-center text-surface-500">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-primary-500 border-t-transparent mb-4"></div>
					<p>Loading financial data...</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full text-left">
						<thead class="bg-surface-200-800/50 text-surface-700-300">
							<tr>
								<th class="p-4">Account Name</th>
								<th class="p-4 text-right">Current Balance</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-surface-200-800">
							{#each balances as balance (balance.id)}
								<tr class="hover:bg-surface-200-800/30 transition-colors">
									<td class="p-4">
										<div class="flex items-center gap-3">
											<div class="w-8 h-8 rounded bg-surface-200-800 flex items-center justify-center text-primary-500">
												<Landmark size={16} />
											</div>
											<span class="font-semibold text-slate-700">{balance.account_name}</span>
										</div>
									</td>
									<td class="p-4 text-right font-bold text-primary-500">
										{formatCurrency(balance.balance)}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</section>
	{/if}
</div>
