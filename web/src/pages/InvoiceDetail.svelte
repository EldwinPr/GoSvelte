<script lang="ts">
	import { onMount } from 'svelte';
	import { RefreshCcw, Printer, ArrowLeft, CreditCard, Save, CheckCircle, Plus, Trash2, X, Landmark, FileText, User, Clock, HandCoins } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { link } from 'svelte-spa-router';

	let { params } = $props<{ params: { id?: string } }>();
	const id = $derived(params?.id || new URLSearchParams(window.location.hash.split('?')[1]).get('id'));

	interface InvoiceDetail {
		id?: string;
		description: string;
		quantity: number;
		unit_price: number;
		subtotal: number;
	}

	interface Payment {
		id: string;
		amount: number;
		payment_date: string;
		method: string;
	}

	interface Invoice {
		id: string;
		number: string;
		customer_id: string;
		invoice_date: string;
		due_date: string;
		total_amount: number;
		status: string;
		customer?: { customer_name: string };
		details: InvoiceDetail[];
		payments: Payment[];
	}

	interface CompanyBalance {
		id: string;
		account_name: string;
		balance: number;
	}

	let invoice = $state<Invoice | null>(null);
	let isLoading = $state(true);
	let isSaving = $state(false);
	let showPaymentForm = $state(false);
	let error = $state<string | null>(null);

	let balances = $state<CompanyBalance[]>([]);
	let selectedBalanceId = $state('');

	let paymentData = $state({
		amount: 0,
		method: 'Bank Transfer',
		payment_date: new Date().toISOString().split('T')[0]
	});

	let totalPaid = $derived(invoice?.payments?.reduce((acc, curr) => acc + curr.amount, 0) || 0);
	let calculatedTotal = $derived(invoice?.details?.reduce((acc, curr) => acc + (curr.quantity * curr.unit_price), 0) || 0);
	let remainingBalance = $derived(calculatedTotal - totalPaid);

	async function fetchInvoice() {
		if (!id) return;
		isLoading = true;
		error = null;
		try {
			invoice = await api<Invoice>(`/api/invoices/show?id=${id}`);
			paymentData.amount = calculatedTotal - totalPaid;
		} catch (e: any) {
			error = "Gagal memuat detail invoice.";
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		if (id) fetchInvoice();
	});

	async function fetchBalances() {
		try {
			const result = await api<any>('/api/balances?page_size=100');
			balances = result.items || [];
			if (balances.length > 0) selectedBalanceId = balances[0].id;
		} catch (e) {}
	}

	function openPaymentForm() {
		showPaymentForm = true;
		fetchBalances();
	}

	async function processPayment(e: SubmitEvent) {
		e.preventDefault();
		if (!invoice || isSaving || !selectedBalanceId) return;
		isSaving = true;

		try {
			await api('/api/invoices/pay', {
				method: 'POST',
				body: JSON.stringify({
					invoice_id: invoice.id,
					balance_id: selectedBalanceId,
					amount: paymentData.amount,
					method: paymentData.method,
					payment_date: paymentData.payment_date
				})
			});
			showPaymentForm = false;
			await fetchInvoice();
		} catch (e: any) {
			alert("Gagal memproses pembayaran: " + e.message);
		} finally {
			isSaving = false;
		}
	}

	async function saveInvoice(finalize = false) {
		if (!invoice || isSaving) return;
		isSaving = true;
		
		const updatedInvoice = { ...invoice };
		if (finalize) updatedInvoice.status = 'Finalized';

		try {
			await api('/api/invoices', {
				method: 'PUT',
				body: JSON.stringify(updatedInvoice)
			});
			await fetchInvoice();
		} catch (e: any) {
			alert("Gagal menyimpan: " + e.message);
		} finally {
			isSaving = false;
		}
	}

	function addRow() {
		if (invoice?.status === 'Finalized' || invoice?.status === 'Paid') return;
		invoice?.details.push({ description: '', quantity: 1, unit_price: 0, subtotal: 0 });
	}

	function removeRow(index: number) {
		if (invoice?.status === 'Finalized' || invoice?.status === 'Paid') return;
		invoice!.details = invoice!.details.filter((_, i) => i !== index);
		if (invoice!.details.length === 0) addRow();
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return "-";
		try {
			return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' });
		} catch (e) { return dateStr; }
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

<div class="max-w-5xl mx-auto space-y-8 text-left">
	<!-- Unified Detail Header -->
	<header class="flex flex-col md:flex-row md:items-center justify-between gap-6">
		<div class="flex items-center gap-4">
			<a use:link href="/invoices" class="btn btn-sm bg-white border border-slate-200 text-slate-400 hover:text-teal-600 transition-colors">
				<ArrowLeft size={18} />
			</a>
			<div>
				<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Invoice Details</h2>
				{#if invoice}
					<p class="text-xs font-black text-teal-600 uppercase tracking-[0.1em]">#{invoice.number} — {invoice.status}</p>
				{/if}
			</div>
		</div>
		
		<div class="flex gap-2">
			{#if invoice && invoice.status !== 'Finalized' && invoice.status !== 'Paid'}
				<button class="btn btn-sm bg-white border border-teal-100 text-teal-600 font-black uppercase tracking-widest text-[9px] px-4 py-3" onclick={() => saveInvoice(false)} disabled={isSaving}>
					<Save size={14} class="mr-2" /> Simpan Perubahan
				</button>
				<button class="btn btn-sm bg-teal-600 text-white font-black uppercase tracking-widest text-[9px] px-4 py-3 shadow-md shadow-teal-600/20" onclick={() => saveInvoice(true)} disabled={isSaving}>
					<CheckCircle size={14} class="mr-2" /> Finalisasi
				</button>
			{/if}
			
			{#if invoice && invoice.status === 'Finalized'}
				<button class="btn btn-sm bg-green-600 text-white font-black uppercase tracking-widest text-[9px] px-4 py-3 shadow-md" onclick={openPaymentForm}>
					<CreditCard size={14} class="mr-2" /> Log Bayar
				</button>
			{/if}
			
			<button class="btn btn-sm bg-white border border-slate-200 text-slate-400 font-black uppercase tracking-widest text-[9px] px-4 py-3 hover:bg-slate-50" onclick={() => window.print()}>
				<Printer size={14} class="mr-2" /> Cetak
			</button>
		</div>
	</header>

	{#if isLoading}
		<div class="p-32 text-center">
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-teal-600 border-t-transparent mb-4"></div>
			<p class="font-black uppercase text-[10px] tracking-[0.2em] text-slate-400">Fetching Record...</p>
		</div>
	{:else if error || !invoice}
		<div class="card p-12 text-center border border-red-100 bg-red-50">
			<p class="text-red-600 font-black uppercase tracking-widest text-xs mb-2">Error 404</p>
			<p class="text-slate-600 font-bold">{error || 'Invoice record could not be found.'}</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<div class="lg:col-span-2 space-y-6">
				<!-- Details Table -->
				<section class="card bg-white border border-slate-200 overflow-hidden shadow-none">
					<header class="p-4 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
						<h3 class="font-black text-slate-400 text-[9px] uppercase tracking-[0.2em]">Itemized Billing</h3>
						{#if invoice.status !== 'Finalized' && invoice.status !== 'Paid'}
							<button class="btn btn-xs bg-white border border-teal-100 text-teal-600 font-black uppercase text-[9px] tracking-widest py-2 px-3" onclick={addRow}>
								<Plus size={12} class="mr-1" /> Baris Baru
							</button>
						{/if}
					</header>

					<div class="overflow-x-auto">
						<table class="table w-full">
							<thead class="bg-slate-50/50 text-slate-400 font-black text-[9px] uppercase border-b border-slate-100 text-left">
								<tr>
									<th class="p-4">Deskripsi</th>
									<th class="p-4 text-center w-20">Qty</th>
									<th class="p-4 text-right w-32">Harga Satuan</th>
									<th class="p-4 text-right w-32">Subtotal</th>
									<th class="p-4 w-10"></th>
								</tr>
							</thead>
							<tbody class="divide-y divide-slate-100">
								{#each invoice.details as detail, i}
									<tr class="hover:bg-slate-50/50 transition-colors">
										<td class="p-2">
											<input class="input border-transparent bg-transparent focus:bg-white text-sm font-bold w-full" type="text" bind:value={detail.description} disabled={invoice.status === 'Finalized' || invoice.status === 'Paid'} />
										</td>
										<td class="p-2 text-center">
											<input class="input border-transparent bg-transparent focus:bg-white text-center text-sm font-bold w-full" type="number" bind:value={detail.quantity} disabled={invoice.status === 'Finalized' || invoice.status === 'Paid'} />
										</td>
										<td class="p-2">
											<div class="relative">
												<span class="absolute left-3 top-1/2 -translate-y-1/2 text-[10px] font-black text-slate-300 z-10">Rp</span>
												<input class="input border-transparent bg-transparent focus:bg-white text-right text-sm font-bold pl-10 w-full" type="number" bind:value={detail.unit_price} disabled={invoice.status === 'Finalized' || invoice.status === 'Paid'} />
											</div>
										</td>
										<td class="p-4 text-right text-sm font-black text-slate-900">{formatCurrency(detail.quantity * detail.unit_price)}</td>
										<td class="p-2 text-center">
											{#if invoice.status !== 'Finalized' && invoice.status !== 'Paid'}
												<button class="text-slate-300 hover:text-red-600 transition-all" onclick={() => removeRow(i)}><Trash2 size={14} /></button>
											{/if}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>

					<footer class="p-6 bg-slate-50 border-t border-slate-200 flex justify-end items-center gap-8">
						<span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total Tagihan</span>
						<span class="text-xl font-black text-teal-600">{formatCurrency(calculatedTotal)}</span>
					</footer>
				</section>

				<!-- Payments -->
				<section class="card bg-white border border-slate-200 overflow-hidden shadow-none">
					<header class="p-4 bg-slate-50 border-b border-slate-200">
						<h4 class="font-black text-slate-400 text-[9px] uppercase tracking-[0.2em]">Riwayat Pembayaran</h4>
					</header>
					{#if !invoice.payments || invoice.payments.length === 0}
						<div class="p-12 text-center text-slate-300 italic text-xs font-bold uppercase tracking-widest">Belum ada dana masuk.</div>
					{:else}
						<div class="overflow-x-auto">
							<table class="table w-full text-left">
								<thead class="bg-slate-50/50 text-slate-400 font-black text-[9px] uppercase border-b border-slate-100">
									<tr>
										<th class="p-4">Tanggal</th>
										<th class="p-4">Metode</th>
										<th class="p-4 text-right">Jumlah</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-slate-100">
									{#each invoice.payments as payment (payment.id)}
										<tr class="hover:bg-slate-50/50 transition-colors">
											<td class="p-4 text-xs font-bold text-slate-500">{formatDate(payment.payment_date)}</td>
											<td class="p-4 text-xs font-black text-slate-700 uppercase">{payment.method}</td>
											<td class="p-4 text-right text-sm font-black text-green-600">{formatCurrency(payment.amount)}</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</section>
			</div>

			<aside class="space-y-6">
				<!-- Customer -->
				<section class="card p-6 bg-white border border-slate-200 space-y-4">
					<h4 class="font-black text-[9px] uppercase text-slate-300 tracking-widest border-b border-slate-100 pb-2">Informasi Client</h4>
					<div>
						<p class="font-black text-slate-800 text-sm">{invoice.customer?.customer_name || 'General Customer'}</p>
						<p class="text-[10px] font-mono text-slate-400">REF: {invoice.customer_id || '-'}</p>
					</div>
					<div class="grid grid-cols-2 gap-4 pt-4 border-t border-slate-50">
						<div>
							<p class="text-[8px] font-black text-slate-400 uppercase tracking-widest">Terbit</p>
							<p class="text-xs font-bold text-slate-700 opacity-80">{formatDate(invoice.invoice_date)}</p>
						</div>
						<div>
							<p class="text-[8px] font-black text-slate-400 uppercase tracking-widest">Tempo</p>
							<p class="text-xs font-bold text-slate-700 opacity-80">{formatDate(invoice.due_date)}</p>
						</div>
					</div>
				</section>

				<!-- Summary Card -->
				<section class="card p-8 bg-teal-600 text-white space-y-4 shadow shadow-teal-600/20">
					<h4 class="font-black text-[9px] uppercase tracking-widest opacity-70">Sisa Tagihan</h4>
					<p class="text-3xl font-black">{formatCurrency(calculatedTotal - totalPaid)}</p>
					<div class="space-y-2 border-t border-white/20 pt-4">
						<div class="flex justify-between text-[9px] font-black uppercase">
							<span>Terkumpul</span>
							<span>{formatCurrency(totalPaid)}</span>
						</div>
						<div class="w-full bg-white/20 h-1 overflow-hidden border border-white/10">
							<div class="bg-white h-full transition-all duration-1000" style="width: {calculatedTotal > 0 ? (totalPaid / calculatedTotal) * 100 : 0}%"></div>
						</div>
					</div>
				</section>
			</aside>
		</div>
	{/if}

	<!-- Payment Form Modal -->
	{#if showPaymentForm}
		<div class="fixed inset-0 bg-slate-900/60 z-50 flex items-center justify-center p-4">
			<div class="card p-10 bg-white border border-slate-200 max-w-md w-full space-y-8 text-left shadow-2xl">
				<div class="flex justify-between items-center border-b border-slate-100 pb-6">
					<h3 class="h3 font-black text-slate-900 flex items-center gap-3 uppercase tracking-tight">
						<CreditCard class="text-teal-600" /> Log Dana Masuk
					</h3>
					<button class="text-slate-400 hover:text-slate-600 p-2" onclick={() => showPaymentForm = false}><X size={28} /></button>
				</div>

				<form onsubmit={processPayment} class="space-y-6" autocomplete="off">
					<label class="label">
						<span class="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2 block">Akun Tujuan</span>
						<select class="select font-bold" bind:value={selectedBalanceId} required>
							{#each balances as balance (balance.id)}
								<option value={balance.id}>{balance.account_name} — {formatCurrency(balance.balance)}</option>
							{/each}
						</select>
					</label>

					<label class="label">
						<span class="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2 block">Jumlah Bayar</span>
						<div class="relative">
							<span class="absolute left-4 top-1/2 -translate-y-1/2 font-black text-slate-200 text-sm z-10">Rp</span>
							<input class="input pl-12 font-black text-slate-900" type="number" bind:value={paymentData.amount} max={remainingBalance} required />
						</div>
						<p class="text-[9px] text-slate-400 mt-2 italic font-black uppercase tracking-widest">SISA TAGIHAN: {formatCurrency(remainingBalance)}</p>
					</label>

					<label class="label">
						<span class="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2 block">Metode Pembayaran</span>
						<select class="select font-bold" bind:value={paymentData.method}>
							<option value="Bank Transfer">Transfer Bank</option>
							<option value="Cash">Tunai (Cash)</option>
							<option value="Credit Card">Kartu Kredit</option>
						</select>
					</label>

					<label class="label">
						<span class="text-[9px] font-black uppercase tracking-widest text-slate-400 mb-2 block">Tanggal Pembayaran</span>
						<input class="input font-bold" type="date" bind:value={paymentData.payment_date} required />
					</label>

					<div class="flex gap-4 pt-10 border-t border-slate-100">
						<button type="button" class="btn bg-slate-100 text-slate-600 border border-slate-200 flex-1 font-black uppercase tracking-widest text-xs py-4" onclick={() => showPaymentForm = false}>Batal</button>
						<button type="submit" class="btn bg-teal-600 text-white flex-1 font-black uppercase tracking-widest text-xs py-4 shadow-teal-600/30 hover:opacity-90 transition-opacity" disabled={isSaving || !selectedBalanceId}>
							{#if isSaving} <RefreshCcw size={18} class="animate-spin mr-2" /> {:else} Konfirmasi {/if}
						</button>
					</div>
				</form>
			</div>
		</div>
	{/if}
</div>
