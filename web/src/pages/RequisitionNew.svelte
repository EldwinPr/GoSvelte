<script lang="ts">
	import { ShoppingCart, RefreshCcw, User, AlertCircle, Save } from 'lucide-svelte';
	import { api } from '../lib/api';
	import { auth } from '../lib/auth.svelte';
	import { push } from 'svelte-spa-router';

	let formData = $state({
		name: '',
		category: '',
		type: 'cash',
		description: '',
		amount: 0,
		user_id: null as string | null,
		user_name: null as string | null
	});

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
	<div class="flex flex-col gap-2 text-left">
		<h2 class="h2 text-slate-900 font-black uppercase tracking-tight">Pengajuan</h2>
		<p class="text-slate-500 text-sm font-medium italic opacity-70">Form pengajuan pembelian barang atau pengembalian biaya (Reimburse).</p>
	</div>

	{#if error}
		<aside class="alert preset-tonal-error shadow-none">
			<AlertCircle size={24} />
			<div class="alert-message text-left">
				<h3 class="font-black uppercase tracking-widest text-xs">Kesalahan</h3>
				<p class="text-sm font-bold">{error}</p>
			</div>
		</aside>
	{/if}

	<div class="card p-10 bg-white border border-slate-200 space-y-8">
		<div class="flex items-center gap-3 border-b border-slate-100 pb-6 text-left">
			<div class="p-3 bg-teal-50 text-teal-600 border border-teal-100">
				<ShoppingCart size={28} />
			</div>
			<h3 class="h3 font-black text-slate-800 uppercase tracking-tight">Formulir Pengajuan</h3>
		</div>

		<form onsubmit={handleSubmit} class="space-y-8 text-left" autocomplete="off">
			<!-- User Identification Section -->
			<div class="space-y-4 bg-slate-50 p-6 border border-slate-200">
				<div class="flex flex-col gap-4">
					<label class="flex items-center gap-3 cursor-pointer group">
						<input 
							type="checkbox" 
							class="checkbox" 
							bind:checked={isManualInput} 
							onchange={handleManualToggle}
						/>
						<span class="text-sm font-black text-slate-700 uppercase tracking-tight">
							Bukan Staff Finance (Input Manual)
						</span>
					</label>

					{#if isManualInput}
						<label class="label animate-in fade-in slide-in-from-top-2 duration-300">
							<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Nama Pengaju</span>
							<input 
								class="input font-bold" 
								type="text" 
								bind:value={formData.user_name} 
								required 
							/>
						</label>
					{:else}
						<div class="flex items-center gap-4 p-4 bg-white border border-slate-200 shadow-none">
							<div class="w-10 h-10 bg-teal-600 flex items-center justify-center text-white font-black text-sm">
								{auth.user?.name?.charAt(0) || '?'}
							</div>
							<div class="flex-1">
								<p class="text-sm font-black text-slate-900">{auth.user?.name}</p>
								<p class="text-[9px] text-teal-600 uppercase font-black tracking-[0.2em]">Identitas Otomatis</p>
							</div>
						</div>
					{/if}
				</div>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<label class="label text-left">
					<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Nama Pengajuan</span>
					<input class="input font-bold" type="text" bind:value={formData.name} required />
				</label>
				<label class="label text-left">
					<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Kategori</span>
					<input class="input font-bold" type="text" bind:value={formData.category} required />
				</label>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div class="space-y-3">
					<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block">Tipe Pengajuan</span>
					<div class="flex gap-6 p-1 bg-slate-50 border border-slate-200 w-fit">
						<label class="flex items-center gap-2 cursor-pointer px-4 py-2">
							<input type="radio" name="type" value="cash" bind:group={formData.type} class="radio" />
							<span class="text-xs font-black text-slate-700 uppercase">Cash</span>
						</label>
						<label class="flex items-center gap-2 cursor-pointer px-4 py-2 border-l border-slate-200">
							<input type="radio" name="type" value="reimburse" bind:group={formData.type} class="radio" />
							<span class="text-xs font-black text-slate-700 uppercase">Reimburse</span>
						</label>
					</div>
				</div>
				<label class="label text-left">
					<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Jumlah Dana</span>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 font-black text-slate-300 text-sm z-10 pointer-events-none">Rp</span>
						<input 
							class="input pl-12 font-black text-slate-900 shadow-none" 
							type="text" 
							value={amountDisplay}
							oninput={handleAmountInput}
							required 
						/>
					</div>
				</label>
			</div>

			<label class="label text-left">
				<span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-2">Deskripsi Detail</span>
				<textarea class="textarea font-bold" rows="4" bind:value={formData.description} required></textarea>
			</label>

			<div class="flex gap-4 pt-6 text-left">
				<button type="button" class="btn bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200 flex-1 font-black uppercase tracking-widest text-xs py-4" onclick={() => push('/')}>Batal</button>
				<button type="submit" class="btn bg-teal-600 text-white flex-1 font-black uppercase tracking-widest text-xs py-4 hover:opacity-90 transition-opacity" disabled={isSubmitting}>
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
