<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';
  import { auth } from '../lib/auth.svelte';

  interface DashboardStats {
    total_balance: number;
    transaction_count: number;
    recent_transactions: any[];
    pending_req_count: number;
    recent_requisitions: any[];
    recent_invoices: any[];
    unpaid_invoice_total: number;
    mtd_omset: number;
    mtd_expenses: number;
    monthly_trends: any[];
  }

  let stats = $state<DashboardStats | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);

  // Derived stats
  let netProfit = $derived(stats ? stats.mtd_omset - stats.mtd_expenses : 0);
  let profitMargin = $derived(stats && stats.mtd_omset > 0 ? (netProfit / stats.mtd_omset) * 100 : 0);
  let maxTrendValue = $derived(stats ? Math.max(...stats.monthly_trends.map(t => Math.max(t.income, t.out)), 1) : 1);

  async function fetchStats() {
    try {
      loading = true;
      stats = await api<DashboardStats>('/api/dashboard/stats');
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  onMount(fetchStats);

  function formatCurrency(val: number) {
    return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(val);
  }

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('en-GB', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="space-y-8 animate-in fade-in duration-500">
  <div class="flex justify-between items-center">
    <div>
      <h1 class="h1 text-slate-900 font-black uppercase tracking-tight leading-none">Command Center</h1>
      <p class="text-slate-500 font-medium mt-2">Welcome back, <span class="text-teal-600 font-bold">{auth.user?.name}</span></p>
    </div>
    <button class="btn preset-filled-primary-500" onclick={fetchStats} disabled={loading}>
      {loading ? 'Syncing...' : 'Sync Now'}
    </button>
  </div>

  {#if loading && !stats}
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      {#each Array(4) as _}
        <div class="card p-6 bg-white border border-slate-200 animate-pulse">
          <div class="h-4 bg-slate-100 rounded w-1/2 mb-4"></div>
          <div class="h-8 bg-slate-100 rounded w-3/4"></div>
        </div>
      {/each}
    </div>
  {:else if error}
    <div class="card p-6 bg-red-50 border border-red-200 text-red-700">
      <p class="font-bold uppercase tracking-widest text-xs mb-2">Sync Error</p>
      {error}
    </div>
  {:else if stats}
    <!-- Top Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="card p-6 bg-white border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
        <p class="text-slate-500 font-bold mb-1 uppercase text-[10px] tracking-[0.2em]">Total Liquidity</p>
        <p class="text-2xl text-slate-900 font-black tracking-tighter">{formatCurrency(stats.total_balance)}</p>
        <div class="mt-2 flex items-center gap-1 text-[10px] font-bold text-teal-600">
          <span class="w-2 h-2 rounded-full bg-teal-500 animate-pulse"></span> ACTIVE ACCOUNTS
        </div>
      </div>

      <div class="card p-6 bg-white border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
        <p class="text-slate-500 font-bold mb-1 uppercase text-[10px] tracking-[0.2em]">MTD Omset (Revenue)</p>
        <p class="text-2xl text-teal-600 font-black tracking-tighter">{formatCurrency(stats.mtd_omset)}</p>
        <p class="mt-2 text-[10px] font-bold text-slate-400 uppercase">Gross Income This Month</p>
      </div>

      <div class="card p-6 bg-white border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
        <p class="text-slate-500 font-bold mb-1 uppercase text-[10px] tracking-[0.2em]">MTD Expenses</p>
        <p class="text-2xl text-rose-600 font-black tracking-tighter">{formatCurrency(stats.mtd_expenses)}</p>
        <p class="mt-2 text-[10px] font-bold text-slate-400 uppercase">Operational Burn</p>
      </div>

      <div class="card p-6 bg-slate-900 border border-slate-800 shadow-sm">
        <p class="text-slate-400 font-bold mb-1 uppercase text-[10px] tracking-[0.2em]">Net Performance</p>
        <p class="text-2xl {netProfit >= 0 ? 'text-teal-400' : 'text-rose-400'} font-black tracking-tighter">
          {netProfit >= 0 ? '+' : ''}{formatCurrency(netProfit)}
        </p>
        <p class="mt-2 text-[10px] font-bold text-slate-400 uppercase">MTD Net Profit ({profitMargin.toFixed(1)}%)</p>
      </div>
    </div>

    <!-- Analytics Section -->
    <div class="card p-6 bg-white border border-slate-200 shadow-sm">
      <div class="flex justify-between items-end mb-8">
        <div>
          <h3 class="font-black uppercase tracking-tighter text-slate-900">Financial Performance</h3>
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">6-Month Trend: Omset vs Outflow</p>
        </div>
        <div class="flex gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-teal-500 rounded-sm"></div>
            <span class="text-[10px] font-bold text-slate-500 uppercase">Income</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-rose-500 rounded-sm"></div>
            <span class="text-[10px] font-bold text-slate-500 uppercase">Expense</span>
          </div>
        </div>
      </div>

      <div class="h-48 flex items-end gap-2 md:gap-8 px-8">
        {#each stats.monthly_trends as month}
          <div class="flex-1 flex flex-col items-center gap-1 group relative">
            <!-- Tooltip -->
            <div class="absolute -top-12 bg-slate-900 text-white text-[9px] p-2 rounded opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-10 whitespace-nowrap font-bold shadow-xl">
              IN: {formatCurrency(month.income)}<br/>
              OUT: {formatCurrency(month.out)}
            </div>
            
            <div class="w-full flex items-end justify-center gap-1 h-full">
              <div 
                class="w-8 bg-teal-500/80 group-hover:bg-teal-500 transition-all rounded-t-sm" 
                style="height: {(month.income / maxTrendValue) * 100}%"
              ></div>
              <div 
                class="w-8 bg-rose-500/80 group-hover:bg-rose-500 transition-all rounded-t-sm" 
                style="height: {(month.out / maxTrendValue) * 100}%"
              ></div>
            </div>
            <p class="text-[10px] font-bold text-slate-400 uppercase mt-2">
              {new Date(month.period + '-01').toLocaleDateString('en-GB', { month: 'short', year: '2-digit' })}
            </p>
          </div>
        {/each}
      </div>
    </div>

    <!-- Health Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="card p-4 bg-white border border-slate-200 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-slate-500 font-bold uppercase text-[9px] tracking-[0.2em]">Receivables</p>
          <p class="text-xl text-slate-900 font-black tracking-tighter">{formatCurrency(stats.unpaid_invoice_total)}</p>
        </div>
        <div class="text-right">
          <p class="text-[9px] font-bold text-amber-600 uppercase">Awaiting Collection</p>
          <a href="#/invoices" class="text-[9px] font-black uppercase text-teal-600 hover:underline">View Invoices →</a>
        </div>
      </div>

      <div class="card p-4 bg-white border border-slate-200 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-slate-500 font-bold uppercase text-[9px] tracking-[0.2em]">Pending Requests</p>
          <p class="text-xl text-slate-900 font-black tracking-tighter">{stats.pending_req_count}</p>
        </div>
        <div class="text-right">
          <p class="text-[9px] font-bold text-blue-600 uppercase">Awaiting Approval</p>
          <a href="#/purchasing/approvals" class="text-[9px] font-black uppercase text-teal-600 hover:underline">Review Now →</a>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Recent Activity -->
      <div class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
        <div class="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
          <h3 class="font-black uppercase tracking-tighter text-slate-900">Recent Transactions</h3>
          <a href="#/transactions" class="text-[10px] font-black uppercase text-teal-600 hover:underline">View All</a>
        </div>
        <div class="divide-y divide-slate-100">
          {#each stats.recent_transactions as tx}
            <div class="p-4 flex justify-between items-center hover:bg-slate-50 transition-colors">
              <div>
                <p class="text-sm font-bold text-slate-900">{tx.description}</p>
                <p class="text-[10px] text-slate-500 font-medium uppercase">{formatDate(tx.date)} • {tx.category}</p>
              </div>
              <div class="text-right">
                <p class="text-sm font-black {tx.type === 'Credit' ? 'text-teal-600' : 'text-rose-600'}">
                  {tx.type === 'Credit' ? '+' : '-'}{formatCurrency(tx.amount)}
                </p>
                <p class="text-[10px] font-bold uppercase opacity-50">{tx.type}</p>
              </div>
            </div>
          {/each}
          {#if stats.recent_transactions.length === 0}
            <div class="p-8 text-center text-slate-400 italic text-sm">No recent activity</div>
          {/if}
        </div>
      </div>

      <!-- Quick Info -->
      <div class="space-y-8">
        <!-- Recent Requisitions -->
        <div class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
          <div class="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
            <h3 class="font-black uppercase tracking-tighter text-slate-900">Pending Approvals</h3>
            <a href="#/purchasing/approvals" class="text-[10px] font-black uppercase text-teal-600 hover:underline">Manage</a>
          </div>
          <div class="divide-y divide-slate-100">
            {#each stats.recent_requisitions as req}
              <div class="p-4 flex justify-between items-center hover:bg-slate-50">
                <div>
                  <p class="text-sm font-bold text-slate-900">{req.name}</p>
                  <p class="text-[10px] text-slate-500 font-medium uppercase">By {req.user?.name || 'Unknown'} • {req.category}</p>
                </div>
                <div class="text-right">
                  <p class="text-sm font-black text-slate-900">{formatCurrency(req.amount)}</p>
                  <span class="badge preset-filled-warning-500 text-[8px] font-black uppercase">Pending</span>
                </div>
              </div>
            {/each}
            {#if stats.recent_requisitions.length === 0}
              <div class="p-8 text-center text-slate-400 italic text-sm">Clear for now!</div>
            {/if}
          </div>
        </div>

        <!-- Recent Invoices -->
        <div class="card bg-white border border-slate-200 overflow-hidden shadow-sm">
          <div class="p-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
            <h3 class="font-black uppercase tracking-tighter text-slate-900">Recent Invoices</h3>
            <a href="#/invoices" class="text-[10px] font-black uppercase text-teal-600 hover:underline">View All</a>
          </div>
          <div class="divide-y divide-slate-100">
            {#each stats.recent_invoices as inv}
              <div class="p-4 flex justify-between items-center hover:bg-slate-50">
                <div>
                  <p class="text-sm font-bold text-slate-900">{inv.number}</p>
                  <p class="text-[10px] text-slate-500 font-medium uppercase">{inv.customer?.customer_name}</p>
                </div>
                <div class="text-right">
                  <p class="text-sm font-black text-slate-900">{formatCurrency(inv.total_amount)}</p>
                  <span class="badge {inv.status === 'Paid' ? 'preset-filled-success-500' : 'preset-filled-surface-200'} text-[8px] font-black uppercase">{inv.status}</span>
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
