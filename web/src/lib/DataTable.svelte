<script lang="ts" generics="T">
	import { onMount, type Snippet } from 'svelte';
	import { RefreshCcw, Search, ArrowUpDown, ArrowUp, ArrowDown, Inbox } from 'lucide-svelte';
	import { api } from './api';
	import type { Column, PaginatedResult } from './types';

	let { 
		endpoint, 
		columns, 
		rowKey,
		searchPlaceholder = "Search...",
		pageSize = 10,
		extraParams = {},
		headerActions,
		cell
	}: {
		endpoint: string;
		columns: Column[];
		rowKey: string;
		searchPlaceholder?: string;
		pageSize?: number;
		extraParams?: Record<string, string | boolean | number>;
		headerActions?: Snippet;
		cell?: Snippet<[T, string]>;
	} = $props();

	// State
	let data = $state<T[]>([]);
	let totalCount = $state(0);
	let currentPage = $state(1);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let searchTerm = $state('');
	let sortKey = $state('');
	let sortOrder = $state<'asc' | 'desc'>('desc');

	let totalPages = $derived(Math.ceil(totalCount / pageSize));

	export async function fetchData() {
		isLoading = true;
		error = null;
		try {
			const query = new URLSearchParams({
				page: currentPage.toString(),
				page_size: pageSize.toString(),
				sort: sortKey,
				order: sortOrder,
				search: searchTerm,
				...Object.fromEntries(Object.entries(extraParams).map(([k, v]) => [k, v.toString()]))
			});
			const result = await api<PaginatedResult<T>>(`${endpoint}?${query.toString()}`);
			data = result.items || [];
			totalCount = result.total_count || 0;
		} catch (e: any) {
			error = e.message || "Gagal memuat data.";
			data = [];
			totalCount = 0;
		} finally {
			isLoading = false;
		}
	}

	// Debounced Search
	let searchTimeout: any;
	function handleSearchInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			currentPage = 1;
			fetchData();
		}, 400);
	}

	// React to parameter changes (e.g. Tab switching)
	$effect(() => {
		void endpoint; 
		void extraParams;
		currentPage = 1;
		fetchData();
	});

	function handleSort(col: Column) {
		if (!col.sortable) return;
		if (sortKey === col.key) {
			sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
		} else {
			sortKey = col.key;
			sortOrder = 'asc';
		}
		currentPage = 1;
		fetchData();
	}

	function changePage(page: number) {
		if (page >= 1 && (page <= totalPages || totalPages === 0)) {
			currentPage = page;
			fetchData();
		}
	}
</script>

<div class="space-y-6 animate-in fade-in duration-500">
	<!-- Search and Actions Row -->
	<div class="flex flex-col md:flex-row justify-between items-center gap-4">
		<div class="relative w-full md:max-w-md group">
			<Search class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 z-10 transition-colors group-focus-within:text-teal-600" size={18} />
			<input 
				type="text" 
				placeholder={searchPlaceholder} 
				class="input pl-10 w-full text-sm font-bold border-slate-200 focus:ring-0 focus:border-teal-600 transition-all shadow-none"
				bind:value={searchTerm}
				oninput={handleSearchInput}
			/>
		</div>
		<div class="flex items-center gap-2">
			<button onclick={fetchData} class="btn btn-sm bg-slate-100 border border-slate-200 text-slate-500 hover:text-teal-600 transition-colors flex items-center gap-2" title="Refresh">
				<RefreshCcw size={14} class={isLoading ? 'animate-spin' : ''} />
				<span class="text-[10px] font-black uppercase tracking-widest">Refresh</span>
			</button>
			{#if headerActions}
				{@render headerActions()}
			{/if}
		</div>
	</div>

	{#if error}
		<div class="card p-12 text-center border border-red-100 bg-red-50/50">
			<p class="text-red-600 font-black uppercase tracking-widest text-[10px] mb-2">System Error</p>
			<p class="text-slate-600 text-sm font-bold">{error}</p>
		</div>
	{:else}
		<section class="card bg-white border border-slate-200 overflow-hidden shadow-none transition-all duration-300">
			{#if isLoading && data.length === 0}
				<div class="p-32 text-center">
					<div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-teal-600 border-t-transparent mb-4"></div>
					<p class="font-black uppercase text-[10px] tracking-[0.2em] text-slate-400">Synchronizing Data...</p>
				</div>
			{:else if data.length === 0}
				<div class="p-32 text-center space-y-4">
					<div class="inline-flex p-4 bg-slate-50 text-slate-200">
						<Inbox size={48} />
					</div>
					<p class="text-slate-400 font-bold italic text-sm">No records found matching your query.</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="table table-hover w-full border-collapse table-fixed">
						<thead>
							<tr class="bg-slate-50/80 text-slate-500 font-black border-b border-slate-200 text-[10px] uppercase tracking-[0.2em]">
								{#each columns as col}
									<th 
										class="p-4 {col.align === 'center' ? 'text-center' : col.align === 'right' ? 'text-right' : 'text-left'} {col.sortable ? 'cursor-pointer hover:bg-slate-100 transition-colors group/header' : ''}" 
										style="width: {col.width || 'auto'}"
										onclick={() => handleSort(col)}
									>
										<div class="flex items-center gap-2 {col.align === 'center' ? 'justify-center' : col.align === 'right' ? 'justify-end' : 'justify-start'}">
											<span class={col.sortable && sortKey === col.key ? 'text-teal-600' : ''}>{col.label}</span>
											{#if col.sortable}
												{#if sortKey === col.key}
													{#if sortOrder === 'asc'} <ArrowUp size={12} class="text-teal-600"/> {:else} <ArrowDown size={12} class="text-teal-600"/> {/if}
												{:else}
													<ArrowUpDown size={12} class="opacity-0 group-hover/header:opacity-30 transition-opacity" />
												{/if}
											{/if}
										</div>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							{#each data as row (row[rowKey as keyof T])}
								<tr class="hover:bg-teal-50/30 transition-colors text-sm text-slate-900 group">
									{#each columns as col}
										<td class="p-4 {col.align === 'center' ? 'text-center' : col.align === 'right' ? 'text-right' : 'text-left'}">
											{#if cell}
												{@render cell(row, col.key)}
											{:else}
												<span class="font-bold opacity-80 group-hover:opacity-100 transition-opacity">{row[col.key as keyof T]}</span>
											{/if}
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<!-- Pagination Footer -->
				<footer class="p-6 bg-slate-50 border-t border-slate-200 flex flex-col md:flex-row justify-between items-center gap-4">
					<p class="text-[10px] font-black uppercase tracking-widest text-slate-400">
						Results: {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, totalCount)} of {totalCount}
					</p>
					<div class="flex gap-1">
						<button 
							class="btn btn-sm bg-white border border-slate-200 text-slate-600 font-black uppercase text-[9px] tracking-widest py-2 px-4 shadow-sm hover:bg-slate-50 disabled:opacity-30 transition-all" 
							onclick={() => changePage(currentPage - 1)}
							disabled={currentPage === 1}
						>
							Prev
						</button>
						<div class="flex items-center gap-1 bg-white border border-slate-200 px-4 py-2 text-[10px] font-black text-slate-600 uppercase tracking-widest">
							Page <span class="text-teal-600 underline underline-offset-4">{currentPage}</span> / {totalPages || 1}
						</div>
						<button 
							class="btn btn-sm bg-white border border-slate-200 text-slate-600 font-black uppercase text-[9px] tracking-widest py-2 px-4 shadow-sm hover:bg-slate-50 disabled:opacity-30 transition-all" 
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
