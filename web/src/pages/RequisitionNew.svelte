<script lang="ts">
	import { onMount } from 'svelte';
	import { ShoppingCart, RefreshCcw, User, AlertCircle, Save } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { auth } from '../lib/auth.svelte';
	import { push } from 'svelte-spa-router';

	let formData = $state({
		name: '',
		category: 'Kebutuhan Kantor',
		type: 'cash',
		description: '',
		amount: 0,
		user_id: null as string | null,
		user_name: null as string | null
	});

	const categories = ['Kebutuhan Kantor', 'Perjalanan Dinas', 'Makan & Minum', 'Pemeliharaan', 'Lainnya'];
	let isManualInput = $state(false);
	let isSubmitting = $state(false);
	let error = $state<string | null>(null);
	let amountDisplay = $state('');

	function handleAmountInput(e: Event) {
		const input = e.target as HTMLInputElement;
		let value = input.value.replace(/\D/g, '');
		
		if (value === '') {
			amountDisplay = '';
			formData.amount = 0;
			return;
		}

		formData.amount = parseInt(value);
		amountDisplay = new Intl.NumberFormat('id-ID').format(formData.amount);
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (isSubmitting) return;
		isSubmitting = true;
		error = null;

		const dataToSend = { ...formData };
		
		if (isManualInput) {
			dataToSend.user_id = null;
		} else {
			dataToSend.user_id = auth.user?.id || null;
			dataToSend.user_name = auth.user?.name || null;
		}

		try {
			await api('/api/requisitions', {
				method: 'POST',
				body: JSON.stringify(dataToSend)
			});
			push('/purchasing/list');
		} catch (e: any) {
			error = "Gagal membuat pengajuan: " + e.message;
		} finally {
			isSubmitting = false;
		}
	}

	function handleManualToggle() {
		if (isManualInput) {
			formData.user_name = null;
		}
	}
</script>

<div class="max-w-2xl mx-auto space-y-8">
	<!-- Header -->
	<div class="flex flex-col gap-2">
		<h2 class="h2 text-slate-900 font-bold uppercase tracking-tight">Pengajuan</h2>
		<p class="text-surface-500 text-sm italic font-medium">Form pengajuan pembelian barang atau pengembalian biaya (Reimburse).</p>
	</div>

	{#if error}
		<aside class="alert preset-tonal-error">
			<AlertCircle size={24} />
			<div class="alert-message">
				<h3 class="font-bold">Kesalahan</h3>
				<p>{error}</p>
			</div>
		</aside>
	{/if}

	<div class="card p-8 bg-white border border-slate-200 shadow-lg space-y-6">
		<div class="flex items-center gap-2 mb-4">
			<div class="p-3 bg-primary-500/10 rounded-full text-primary-500">
				<ShoppingCart size={24} />
			</div>
			<h3 class="h3 font-bold text-slate-800">Formulir Pengajuan</h3>
		</div>

		<form onsubmit={handleSubmit} class="space-y-6">
			<!-- User Identification Section -->
			<div class="space-y-4 bg-slate-50 p-6 rounded-xl border border-slate-200">
				<div class="flex flex-col gap-4">
					<label class="flex items-center gap-3 cursor-pointer group">
						<input 
							type="checkbox" 
							class="checkbox checkbox-primary w-5 h-5" 
							bind:checked={isManualInput} 
							onchange={handleManualToggle}
						/>
						<span class="text-sm font-bold text-slate-700 group-hover:text-primary-600 transition-colors">
							Bukan Staff Finance (Input Nama Manual)
						</span>
					</label>

					{#if isManualInput}
						<label class="label animate-in fade-in slide-in-from-top-2 duration-300">
							<span class="text-sm font-semibold text-slate-700">Nama Pengaju</span>
							<input 
								class="input bg-white border-slate-200" 
								type="text" 
								bind:value={formData.user_name} 
								placeholder="Masukkan nama lengkap pengaju..." 
								required 
							/>
						</label>
					{:else}
						<div class="flex items-center gap-3 p-3 bg-white border border-slate-200 rounded-lg">
							<div class="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-bold text-xs">
								{auth.user?.name?.charAt(0) || '?'}
							</div>
							<div class="flex-1">
								<p class="text-xs font-bold text-slate-900">{auth.user?.name}</p>
								<p class="text-[10px] text-slate-500 uppercase font-bold tracking-wider">Staff Finance (Otomatis)</p>
							</div>
						</div>
					{/if}
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Nama Pengajuan</span>
					<input class="input bg-white border-slate-200" type="text" bind:value={formData.name} placeholder="Contoh: ATK Kantor" required />
				</label>
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Kategori</span>
					<input class="input bg-white border-slate-200" type="text" bind:value={formData.category} placeholder="Contoh: Kantor, Perjalanan, Makan" required />
				</label>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="space-y-2">
					<span class="text-sm font-semibold text-slate-700">Tipe Pengajuan</span>
					<div class="flex gap-4">
						<label class="flex items-center gap-2 cursor-pointer">
							<input type="radio" name="type" value="cash" bind:group={formData.type} class="radio radio-primary" />
							<span class="text-sm text-slate-700">Cash</span>
						</label>
						<label class="flex items-center gap-2 cursor-pointer">
							<input type="radio" name="type" value="reimburse" bind:group={formData.type} class="radio radio-primary" />
							<span class="text-sm text-slate-700">Reimburse</span>
						</label>
					</div>
				</div>
				<label class="label">
					<span class="text-sm font-semibold text-slate-700">Jumlah Dana (Rupiah)</span>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 font-bold text-slate-400">Rp</span>
						<input 
							class="input pl-12 bg-white border-slate-200" 
							type="text" 
							value={amountDisplay}
							oninput={handleAmountInput}
							placeholder="0" 
							required 
						/>
					</div>
				</label>
			</div>

			<label class="label">
				<span class="text-sm font-semibold text-slate-700">Deskripsi Detail</span>
				<textarea class="textarea bg-white border-slate-200" rows="3" bind:value={formData.description} placeholder="Keterangan tambahan mengenai pengajuan..." required></textarea>
			</label>

			<div class="flex gap-4 pt-4">
				<button type="button" class="btn preset-tonal-surface flex-1 font-bold" onclick={() => push('/')}>Cancel</button>
				<button type="submit" class="btn preset-filled-primary-500 flex-1 font-bold shadow-md shadow-primary-500/20" disabled={isSubmitting}>
					{#if isSubmitting}
						<RefreshCcw size={18} class="animate-spin mr-2" />
						Processing...
					{:else}
						<Save size={18} class="mr-2" />
						Kirim Pengajuan
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
