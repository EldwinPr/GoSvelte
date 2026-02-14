<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '../lib/api';

  interface User {
    id: number;
    name: string;
    email: string;
  }

  let users: User[] = [];
  let error = '';
  let loading = true;

  onMount(async () => {
    try {
      users = await api<User[]>('/api/users');
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  });
</script>

<div class="flex justify-between items-center mb-6">
  <h1 class="text-3xl font-bold text-gray-800">User Management</h1>
  <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors shadow-sm">
    Add User
  </button>
</div>

{#if loading}
  <div class="flex justify-center py-12">
    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
  </div>
{:else if error}
  <div class="bg-red-50 border-l-4 border-red-500 p-4 mb-6">
    <p class="text-red-700 font-medium">Error: {error}</p>
  </div>
{:else}
  <div class="bg-white shadow-md rounded-lg overflow-hidden border border-gray-200">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
          <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        {#each users as user}
          <tr class="hover:bg-gray-50 transition-colors">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{user.id}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{user.name}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{user.email}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button class="text-blue-600 hover:text-blue-900 mr-3">Edit</button>
              <button class="text-red-600 hover:text-red-900">Delete</button>
            </td>
          </tr>
        {:else}
          <tr>
            <td colspan="4" class="px-6 py-12 text-center text-gray-500 italic">
              No users found in the system.
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}
