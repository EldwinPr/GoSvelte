<script lang="ts">
	import type { Snippet } from 'svelte';
	import { link, location } from 'svelte-spa-router';
	import { auth } from './auth.svelte';
	import { 
		LayoutDashboard, Receipt, Wallet, CreditCard, ShoppingCart, 
		CheckCircle, FileText, Landmark, Users, Database, Menu, LogOut, 
		HandCoins 
	} from 'lucide-svelte';

	let { children } = $props<{ children: Snippet }>();

	let isSidebarOpen = $state(true);

	const pageTitles: Record<string, string> = {
		'/': 'Dashboard',
		'/invoices': 'Invoices',
		'/invoices/new': 'Create Invoice',
		'/invoices/detail': 'Invoice Details',
		'/payments': 'Payments',
		'/credits': 'Customer Credits',
		'/purchasing/new': 'New Requisition',
		'/purchasing/list': 'Requisition List',
		'/purchasing/approvals': 'Requisition Approvals',
		'/transactions': 'Transaction History',
		'/balance': 'Company Balance',
		'/users': 'User Management',
		'/explorer': 'Database Explorer'
	};

	const currentPath = $derived($location);
	const currentPageName = $derived(pageTitles[currentPath] || 'System');

	const menuGroups = $derived([
		{
			label: 'MAIN',
			minClearance: 0,
			items: [ { label: 'Dashboard', icon: LayoutDashboard, path: '/', minClearance: 0 } ]
		},
		{
			label: 'FINANCE',
			minClearance: 0,
			items: [
				{ label: 'Invoices', icon: Receipt, path: '/invoices', minClearance: 0 },
				{ label: 'Recent Payments', icon: Wallet, path: '/payments', minClearance: 0 },
				{ label: 'Customer Credits', icon: CreditCard, path: '/credits', minClearance: 0 }
			]
		},
		{
			label: 'PURCHASING',
			minClearance: 0,
			items: [
				{ label: 'Pengajuan Baru', icon: ShoppingCart, path: '/purchasing/new', minClearance: 0 },
				{ label: 'Daftar Pengajuan', icon: FileText, path: '/purchasing/list', minClearance: 0 },
				{ label: 'Persetujuan (Manager)', icon: CheckCircle, path: '/purchasing/approvals', minClearance: 10 }
			]
		},
		{
			label: 'REPORTS',
			minClearance: 0,
			items: [
				{ label: 'Transactions', icon: FileText, path: '/transactions', minClearance: 0 },
				{ label: 'Company Balance', icon: Landmark, path: '/balance', minClearance: 10 }
			]
		},
		{
			label: 'SYSTEM',
			minClearance: 20,
			items: [
				{ label: 'User Management', icon: Users, path: '/users', minClearance: 20 },
				{ label: 'Database Explorer', icon: Database, path: '/explorer', minClearance: 20 }
			]
		}
	]);

	function toggleSidebar() { isSidebarOpen = !isSidebarOpen; }
</script>

<div class="flex h-screen overflow-hidden bg-slate-50 text-slate-900 transition-colors duration-200">
	<!-- Sidebar -->
	<aside class="flex flex-col border-r border-slate-200 transition-all duration-300 bg-white {isSidebarOpen ? 'w-64' : 'w-0 -ml-64'} md:ml-0 md:w-64 shrink-0">
		<div class="p-4 flex items-center gap-3 border-b border-slate-200 h-16 shrink-0 bg-white">
			<HandCoins class="text-teal-600" size={32} />
			<span class="text-xl font-black tracking-tight uppercase text-slate-900">Sentral <span class="text-teal-600">Finance</span></span>
		</div>

		<nav class="flex-1 overflow-y-auto p-4 space-y-6 bg-white">
			{#each menuGroups as group}
				{#if (auth.user?.clearance ?? 0) >= group.minClearance || auth.user?.clearance === 20}
					<div>
						<p class="text-[10px] font-black text-slate-400 mb-3 px-2 uppercase tracking-[0.2em]">{group.label}</p>
						<ul class="space-y-1">
							{#each group.items as item}
								{#if (auth.user?.clearance ?? 0) >= item.minClearance || auth.user?.clearance === 20}
									<li>
										<a 
											use:link 
											href={item.path} 
											class="flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all text-sm font-bold {currentPath === item.path ? 'bg-teal-600 text-white' : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
										>
											<item.icon size={18} />
											<span class="truncate">{item.label}</span>
										</a>
									</li>
								{/if}
							{/each}
						</ul>
					</div>
				{/if}
			{/each}
		</nav>

		<div class="p-4 border-t border-slate-200 bg-slate-50 shrink-0">
			<div class="flex items-center gap-3 px-2 py-1 text-left">
				<div class="w-9 h-9 rounded-full bg-teal-600 flex items-center justify-center text-white font-black shrink-0">
					{auth.user?.name?.charAt(0) || '?'}
				</div>
				<div class="flex-1 min-w-0">
					<p class="text-xs font-black truncate text-slate-900">{auth.user?.name}</p>
					<p class="text-[10px] text-slate-500 font-bold uppercase tracking-tighter">
						{auth.user?.clearance === 20 ? 'Developer' : auth.user?.clearance === 10 ? 'Manager' : 'Finance'}
					</p>
				</div>
				<button onclick={() => auth.logout()} title="Logout" class="p-1.5 hover:bg-red-50 hover:text-red-600 rounded-md transition-colors shrink-0">
					<LogOut size={16} />
				</button>
			</div>
		</div>
	</aside>

	<!-- Main Content Area -->
	<div class="flex-1 flex flex-col min-w-0 overflow-hidden bg-slate-50">
		<!-- Top Bar -->
		<header class="h-16 flex items-center justify-between px-6 border-b border-slate-200 bg-white z-10 shrink-0">
			<div class="flex items-center gap-4">
				<button onclick={toggleSidebar} class="p-2 hover:bg-slate-100 rounded-lg transition-colors md:hidden text-slate-900">
					<Menu size={24} />
				</button>
				<div class="flex items-center gap-2 text-slate-400">
					<span class="text-xs font-black uppercase tracking-[0.2em] hidden sm:block">Page /</span>
					<h1 class="text-sm font-black text-slate-800 uppercase tracking-widest">{currentPageName}</h1>
				</div>
			</div>
		</header>

		<!-- Content -->
		<main class="flex-1 overflow-y-auto p-8 bg-slate-50">
			<div class="max-w-7xl mx-auto">
				{@render children()}
			</div>
		</main>
	</div>
</div>
