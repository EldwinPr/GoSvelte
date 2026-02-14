<script lang="ts">
  import Router from 'svelte-spa-router'
  import Dashboard from './pages/Dashboard.svelte'
  import Users from './pages/Users.svelte'
  import NotFound from './pages/NotFound.svelte'

  const routes = {
    '/': Dashboard,
    '/users': Users,
    '*': NotFound,
  }
</script>

<nav class="bg-slate-800 text-white shadow-lg">
  <div class="container mx-auto px-4 py-3 flex justify-between items-center">
    <div class="text-xl font-bold tracking-tight">GoSvelte <span class="text-blue-400">ERP</span></div>
    <ul class="flex space-x-6">
      <li><a href="#/" class="hover:text-blue-400 transition-colors">Dashboard</a></li>
      <li><a href="#/users" class="hover:text-blue-400 transition-colors">Users</a></li>
    </ul>
  </div>
</nav>

<main class="container mx-auto px-4 py-8">
  <Router {routes} />
</main>

<style>
  /* Tailwind is global, no need for local styles here for now */
</style>
