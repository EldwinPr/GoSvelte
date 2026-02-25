<script lang="ts">
  import Router, { push } from 'svelte-spa-router'
  import { onMount } from 'svelte'
  import { auth } from './lib/auth.svelte'
  import Layout from './lib/Layout.svelte'
  import Dashboard from './pages/Dashboard.svelte'
  import Users from './pages/Users.svelte'
  import Explorer from './pages/Explorer.svelte'
  import Login from './pages/Login.svelte'
  import Balance from './pages/Balance.svelte'
  import Transactions from './pages/Transactions.svelte'
  import RequisitionNew from './pages/RequisitionNew.svelte'
  import RequisitionList from './pages/RequisitionList.svelte'
  import RequisitionHistory from './pages/RequisitionHistory.svelte'
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
    // Add placeholders for the other routes to prevent 404s
    '/invoices': Dashboard,
    '/payments': Dashboard,
    '/credits': Dashboard,
    '*': NotFound,
  }

  // Simple route guard logic
  function handleRoute(event: any) {
    if (isInitializing) return;
    const detail = event.detail;
    if (detail.location !== '/login' && !auth.user) {
      push('/login');
    }
  }

  // Start Auth Check immediately (Eager)
  const authPromise = auth.init();

  let isInitializing = $state(true);

  onMount(async () => {
    try {
      await authPromise;
    } finally {
      isInitializing = false;
    }
  });
</script>

{#if isInitializing}
  <div class="h-screen flex items-center justify-center bg-slate-50">
    <div class="flex flex-col items-center gap-4">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
      <p class="text-sm font-bold text-slate-400 uppercase tracking-widest">Initializing ERP...</p>
    </div>
  </div>
{:else}
  {#if auth.user}
    <Layout user={auth.user}>
      <Router {routes} on:routeLoaded={handleRoute} />
    </Layout>
  {:else}
    <div class="min-h-screen bg-slate-50">
      <Router {routes} on:routeLoaded={handleRoute} />
    </div>
  {/if}
{/if}

<style>
  /* Base styles are handled in app.css and Skeleton UI */
</style>
