<script lang="ts">
  import Router, { push } from 'svelte-spa-router'
  import { onMount } from 'svelte'
  import { auth } from './lib/auth.svelte'
  import Layout from './lib/Layout.svelte'
  import Dashboard from './pages/Dashboard.svelte'
  import Users from './pages/Users.svelte'
  import Explorer from './pages/Explorer.svelte'
  import Login from './pages/Login.svelte'
  import NotFound from './pages/NotFound.svelte'

  const routes = {
    '/': Dashboard,
    '/login': Login,
    '/users': Users,
    '/explorer': Explorer,
    // Add placeholders for the other routes to prevent 404s
    '/invoices': Dashboard,
    '/payments': Dashboard,
    '/credits': Dashboard,
    '/purchasing/new': Dashboard,
    '/purchasing/approvals': Dashboard,
    '/transactions': Dashboard,
    '/balance': Dashboard,
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

  let isInitializing = $state(true);

  onMount(async () => {
    await auth.init();
    isInitializing = false;
  });
</script>

{#if isInitializing}
  <div class="h-screen flex items-center justify-center bg-slate-50">
    <div class="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
  </div>
{:else if auth.user}
  <Layout user={auth.user}>
    <Router {routes} on:routeLoaded={handleRoute} />
  </Layout>
{:else}
  <Router {routes} on:routeLoaded={handleRoute} />
{/if}

<style>
  /* Base styles are handled in app.css and Skeleton UI */
</style>
