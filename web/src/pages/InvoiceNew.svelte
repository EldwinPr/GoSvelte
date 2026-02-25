<script lang="ts">
	import { onMount } from 'svelte';
	import { Receipt, Plus, Trash2, Save, RefreshCcw, AlertCircle, Calendar } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { push } from 'svelte-spa-router';

	interface Customer {
		id: string;
		customer_name: string;
	}

	interface InvoiceDetail {
		description: string;
		quantity: number;
		unit_price: number;
		subtotal: number;
	}

	let formData = $state({
		number: `INV-${Date.now().toString().slice(-6)}`,
		customer_id: '',
		invoice_date: new Date().toISOString().split('T')[0],
		due_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
		details: [] as InvoiceDetail[]
	});

	let isSubmitting = $state(false);
	let error = $state<string | null>(null);

	onMount(async () => {
		addRow();
	});

	function addRow() {
		formData.details.push({
			description: '',
			quantity: 1,
			unit_price: 0,
			subtotal: 0
		});
	}

	function removeRow(index: number) {
		formData.details = formData.details.filter((_, i) => i !== index);
		if (formData.details.length === 0) addRow();
	}

	let totalAmount = $derived(
		formData.details.reduce((acc, curr) => acc + (curr.quantity * curr.unit_price), 0)
	);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (isSubmitting) return;
		isSubmitting = true;
		error = null;

		try {
			await api('/api/invoices', {
				method: 'POST',
				body: JSON.stringify(formData)
			});
			push('/invoices');
		} catch (e: any) {
			error = "Gagal membuat invoice: " + e.message;
		} finally {
			isSubmitting = false;
		}
	}

	function formatCurrency(amount: number) {
		return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount || 0);
	}
</script>

<div class="max-w-4xl mx-auto space-y-8">
	<div class="flex flex-col gap-2 text-left">
		<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Buat Invoice Baru</h2>
		<p class="text-slate-500 text-sm font-medium italic">Lengkapi rincian penagihan untuk dikirim ke pelanggan.</p>
	</div>

	{#if error}
		<aside class="alert preset-tonal-error shadow-none">
			<AlertCircle size={24} />
			<div class="alert-message">
				<h3 class="font-bold text-slate-900 text-left">Kesalahan</h3>
				<p class="text-slate-800 text-left">{error}</p>
			</div>
		</aside>
	{/if}

	<form onsubmit={handleSubmit} class="space-y-6">
		<div class="card p-8 bg-white border border-slate-200 shadow-sm grid grid-cols-1 md:grid-cols-2 gap-6">
			<div class="flex items-center gap-2 mb-2 md:col-span-2 border-b border-slate-100 pb-4">
				<div class="p-2 bg-teal-50 rounded-lg text-teal-600">
					<Receipt size={20} />
				</div>
				<h3 class="h3 font-black text-slate-800 uppercase tracking-tight">Informasi Dasar</h3>
			</div>

			<label class="label text-left">
				<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Nomor Invoice</span>
				<input class="input font-mono font-bold" type="text" bind:value={formData.number} required />
			</label>

			<label class="label text-left">
				<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Pelanggan</span>
				<input class="input font-bold" type="text" placeholder="Nama Pelanggan / ID..." required />
			</label>

			<label class="label text-left">
				<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-2 mb-2">
					<Calendar size={14} /> Tanggal Invoice
				</span>
				<input class="input font-bold" type="date" bind:value={formData.invoice_date} required />
			</label>

			<label class="label text-left">
				<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-2 mb-2">
					<Calendar size={14} /> Jatuh Tempo
				</span>
				<input class="input font-bold" type="date" bind:value={formData.due_date} required />
			</label>
		</div>

		<div class="card p-0 bg-white border border-slate-200 shadow-sm overflow-hidden">
			<header class="p-6 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
				<h3 class="font-black text-slate-400 uppercase text-[10px] tracking-[0.2em]">Rincian Barang / Jasa</h3>
				<button type="button" class="btn btn-sm bg-teal-50 text-teal-600 font-black uppercase text-[10px] tracking-widest px-4 border border-teal-100" onclick={addRow}>
					<Plus size={14} class="mr-1" /> Tambah Baris
				</button>
			</header>

			<div class="overflow-x-auto">
				<table class="table w-full">
					<thead class="bg-slate-50 text-slate-400 font-black text-[10px] uppercase border-b border-slate-100">
						<tr>
							<th class="p-4 w-1/2 text-left">Deskripsi</th>
							<th class="p-4 w-20 text-center">Qty</th>
							<th class="p-4 text-right">Harga Satuan</th>
							<th class="p-4 text-right">Subtotal</th>
							<th class="p-4 w-10"></th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each formData.details as row, i}
							<tr class="hover:bg-slate-50/50 transition-colors">
								<td class="p-2 text-left">
									<input class="input border-transparent bg-transparent focus:bg-slate-50 text-sm font-bold w-full" type="text" bind:value={row.description} required />
								</td>
								<td class="p-2 text-center">
									<input class="input border-transparent bg-transparent focus:bg-slate-50 text-center text-sm font-bold w-full" type="number" bind:value={row.quantity} min="1" required />
								</td>
								<td class="p-2">
									<div class="relative">
										<span class="absolute left-4 top-1/2 -translate-y-1/2 text-xs font-black text-slate-300 z-10 pointer-events-none">Rp</span>
										<input class="input border-transparent bg-transparent focus:bg-slate-50 text-right text-sm font-bold pl-12 w-full" type="number" bind:value={row.unit_price} min="0" required />
									</div>
								</td>
								<td class="p-4 text-right font-black text-sm text-slate-900">
									{formatCurrency(row.quantity * row.unit_price)}
								</td>
								<td class="p-2 text-center">
									<button type="button" class="text-slate-300 hover:text-red-600 transition-colors" onclick={() => removeRow(i)}>
										<Trash2 size={16} />
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<footer class="p-8 bg-slate-50 flex flex-col items-end gap-2 border-t border-slate-200">
				<div class="flex items-center gap-12">
					<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Invoice</span>
					<span class="text-3xl font-black text-teal-600">{formatCurrency(totalAmount)}</span>
				</div>
			</footer>
		</div>

		<div class="flex gap-4 pt-4">
			<button type="button" class="btn bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200 flex-1 font-black uppercase tracking-widest text-xs py-4" onclick={() => push('/invoices')}>Cancel</button>
			<button type="submit" class="btn bg-teal-600 text-white flex-1 font-black uppercase tracking-widest text-xs py-4 shadow-md shadow-teal-600/20 hover:opacity-90" disabled={isSubmitting}>
				{#if isSubmitting}
					<RefreshCcw size={18} class="animate-spin mr-2" />
					Processing...
				{:else}
					<Save size={18} class="mr-2" />
					Simpan Invoice
				{/if}
			</button>
		</div>
	</form>
</div>
