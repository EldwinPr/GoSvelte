<script lang="ts">
	import type { Snippet } from 'svelte';
	import { link } from 'svelte-spa-router';
	import { auth } from './auth.svelte';
	import { HandCoins, LayoutDashboard, Receipt, Wallet, CreditCard, ShoppingCart, CheckCircle, FileText, Landmark, Users, Database, Menu, LogOut } from 'lucide-svelte';

	interface User {
		name: string;
		clearance: number;
	}

	interface Props {
		children: Snippet;
		user: User;
	}

	let { children, user }: Props = $props();

	let isSidebarOpen = $state(true); 

	// Dynamic Page Title mapping
	const pageTitles: Record<string, string> = {
		'/': 'Dashboard',
		'/invoices': 'Invoices',
		'/payments': 'Payments',
		'/credits': 'Customer Credits',
		'/purchasing/new': 'New Requisition',
		'/purchasing/approvals': 'Requisition Approvals',
		'/transactions': 'Transaction History',
		'/balance': 'Company Balance',
		'/users': 'User Management',
		'/explorer': 'Database Explorer'
	};

	let currentPath = $state(window.location.hash.replace('#', '') || '/');
	
	// Sync path on hash change
	window.addEventListener('hashchange', () => {
		currentPath = window.location.hash.replace('#', '') || '/';
	});

	const currentPageName = $derived(pageTitles[currentPath] || 'ERP System');

	const menuGroups = $derived([
		{
			label: 'MAIN',
			minClearance: 0,
			items: [
				{ label: 'Dashboard', icon: LayoutDashboard, path: '/', minClearance: 0 }
			]
		},
		{
			label: 'FINANCE',
			minClearance: 0,
			items: [
				{ label: 'Invoices', icon: Receipt, path: '/invoices', minClearance: 0 },
				{ label: 'Payments', icon: Wallet, path: '/payments', minClearance: 0 },
				{ label: 'Customer Credits', icon: CreditCard, path: '/credits', minClearance: 0 }
			]
		},
		{
			label: 'PURCHASING',
			minClearance: 0,
			items: [
				{ label: 'New Requisition', icon: ShoppingCart, path: '/purchasing/new', minClearance: 0 },
				{ label: 'Approvals', icon: CheckCircle, path: '/purchasing/approvals', minClearance: 10 }
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

<div class="flex h-screen overflow-hidden bg-white">
	<!-- Sidebar -->
	<aside class="flex flex-col border-r border-slate-200 transition-all duration-300 bg-slate-50 {isSidebarOpen ? 'w-64' : 'w-0 -ml-64'} md:ml-0 md:w-64">
		<div class="p-4 flex items-center gap-3 border-b border-slate-200 h-16 shrink-0">
			<HandCoins class="text-primary-500" size={32} />
			<span class="text-xl font-bold tracking-tight text-slate-900">Sentral <span class="text-primary-500">Finance</span></span>
		</div>

		<nav class="flex-1 overflow-y-auto p-4 space-y-6">
			{#each menuGroups as group}
				{#if user.clearance >= group.minClearance || user.clearance === 20}
					<div>
						<p class="text-xs font-semibold text-slate-500 mb-2 px-2 uppercase tracking-wider">{group.label}</p>
						<ul class="space-y-1">
							{#each group.items as item}
								{#if user.clearance >= item.minClearance || user.clearance === 20}
									<li>
										<a use:link href={item.path} class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-primary-500/10 hover:text-primary-500 text-slate-700 transition-colors">
											<item.icon size={18} />
											<span class="truncate font-medium">{item.label}</span>
										</a>
									</li>
								{/if}
							{/each}
						</ul>
					</div>
				{/if}
			{/each}
		</nav>

		<div class="p-4 border-t border-slate-200 bg-slate-100 shrink-0">
			<div class="flex items-center gap-3 px-2 py-1">
				<div class="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-bold shrink-0">
					{user.name.charAt(0)}
				</div>
				<div class="flex-1 min-w-0">
					<p class="text-sm font-semibold truncate text-slate-900">{user.name}</p>
					<p class="text-xs text-slate-500 truncate">
						{user.clearance === 20 ? 'Developer' : user.clearance === 10 ? 'Manager' : 'Finance'}
					</p>
				</div>
				<button 
					onclick={() => auth.logout()} 
					title="Logout" 
					class="p-1 text-slate-400 hover:text-error-500 transition-colors shrink-0"
				>
					<LogOut size={18} />
				</button>
			</div>
		</div>
	</aside>

	<!-- Main Content Area -->
	<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
		<!-- Top Bar -->
		<header class="h-16 flex items-center justify-between px-4 border-b border-slate-200 bg-white z-10 shrink-0">
			<div class="flex items-center gap-4">
				<button onclick={toggleSidebar} class="p-2 hover:bg-slate-100 rounded-lg transition-colors md:hidden">
					<Menu size={24} />
				</button>
				<h1 class="text-lg font-bold text-slate-900 hidden sm:block">{currentPageName}</h1>
			</div>
			
			<div class="flex items-center gap-4">
				<div class="text-sm text-slate-600 px-3 py-1 bg-slate-100 rounded-full font-medium">
					System Status: <span class="text-success-600 font-bold">Online</span>
				</div>
			</div>
		</header>

		<!-- Content -->
		<main class="flex-1 overflow-y-auto p-6">
			{@render children()}
		</main>
	</div>
</div>
