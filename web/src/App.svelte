<script lang="ts">
  import Router, { push } from 'svelte-spa-router'
  import { onMount } from 'svelte'
  import { auth } from './lib/auth.svelte'
  import Layout from './lib/Layout.svelte'
  
  // Pages
  import Dashboard from './pages/Dashboard.svelte'
  import Users from './pages/Users.svelte'
  import Explorer from './pages/Explorer.svelte'
  import Login from './pages/Login.svelte'
  import Balance from './pages/Balance.svelte'
  import Transactions from './pages/Transactions.svelte'
  import RequisitionNew from './pages/RequisitionNew.svelte'
  import RequisitionList from './pages/RequisitionList.svelte'
  import RequisitionHistory from './pages/RequisitionHistory.svelte'
  import RequisitionDetail from './pages/RequisitionDetail.svelte'
  import InvoiceList from './pages/InvoiceList.svelte'
  import InvoiceNew from './pages/InvoiceNew.svelte'
  import InvoiceDetail from './pages/InvoiceDetail.svelte'
  import Payments from './pages/Payments.svelte'
  import NotFound from './pages/NotFound.svelte'

  const routes = {
    '/': Dashboard,
    '/login': Login,
    '/users': Users,
    '/explorer': Explorer,
    '/balance': Balance,
    '/transactions': Transactions,
    '/purchasing/new': RequisitionNew,
    '/purchasing/list': RequisitionHistory,
    '/purchasing/approvals': RequisitionList,
    '/purchasing/detail': RequisitionDetail,
    '/invoices': InvoiceList,
    '/invoices/new': InvoiceNew,
    '/invoices/detail': InvoiceDetail,
    '/payments': Payments,
    '/credits': Dashboard,
    '*': NotFound,
  }

  let isReady = $state(false);
  let initError = $state<string | null>(null);

  async function initialize() {
    // Forcibly clear any lingering dark mode
    if (typeof document !== 'undefined') {
      document.documentElement.classList.remove('dark');
      localStorage.removeItem('theme-mode');
    }
    try {
      await auth.init();
      
      const currentPath = window.location.hash.slice(1).split('?')[0] || '/';
      if (currentPath !== '/login' && !auth.user) {
        push('/login');
      }
    } catch (e: any) {
      console.error("[App] Init Error:", e);
      initError = e.message;
    } finally {
      isReady = true;
    }
  }

  onMount(initialize);
</script>

{#if !isReady}
  <div class="h-screen w-screen flex items-center justify-center bg-white text-slate-900">
    <div class="flex flex-col items-center gap-4">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-teal-500 border-t-transparent"></div>
      <p class="text-sm font-bold uppercase tracking-widest opacity-50">System Loading...</p>
    </div>
  </div>
{:else if initError}
  <div class="h-screen w-screen flex items-center justify-center bg-red-50 p-4">
    <div class="card p-8 bg-white border border-red-200 shadow-xl max-w-md">
      <h1 class="text-red-600 font-bold mb-2 uppercase tracking-tight">System Boot Failure</h1>
      <p class="text-sm text-slate-600">{initError}</p>
      <button class="btn preset-tonal-surface mt-6 w-full" onclick={() => window.location.reload()}>Retry</button>
    </div>
  </div>
{:else}
  {#if auth.user}
    <Layout>
      <Router {routes} />
    </Layout>
  {:else}
    <Router {routes} />
  {/if}
{/if}

<style>
  :global(body) {
    margin: 0;
    padding: 0;
  }
</style>
