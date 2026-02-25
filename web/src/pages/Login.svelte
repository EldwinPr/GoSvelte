<script lang="ts">
	import { auth } from '../lib/auth.svelte';
	import { CloudUpload, Lock, Mail, ArrowRight } from 'lucide-svelte';
	import { push } from 'svelte-spa-router';

	let email = $state('');
	let password = $state('');
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isLoading = true;
		error = null;

		try {
			await auth.login(email, password);
			push('/');
		} catch (err: any) {
			error = err.message || 'Invalid credentials';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-slate-50 p-4">
	<div class="w-full max-w-md space-y-8">
		<!-- Logo -->
		<div class="text-center">
			<div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-primary-500 text-white mb-4 shadow-xl shadow-primary-500/20">
				<CloudUpload size={40} />
			</div>
			<h1 class="h1 font-extrabold text-slate-900">Sentral Finance</h1>
			<p class="text-slate-500 mt-2">Sign in to your corporate financial dashboard</p>
		</div>

		<!-- Card -->
		<div class="card p-8 bg-white border border-slate-200 shadow-xl">
			{#if error}
				<aside class="alert preset-tonal-error mb-6">
					<div class="alert-message">
						<p class="text-sm font-medium">{error}</p>
					</div>
				</aside>
			{/if}

			<form onsubmit={handleSubmit} class="space-y-6">
				<label class="label">
					<span class="text-sm font-bold text-slate-700 flex items-center gap-2">
						<Mail size={16} /> Email Address
					</span>
					<input 
						class="input" 
						type="email" 
						bind:value={email} 
						placeholder="admin@erp.com" 
						required 
					/>
				</label>

				<label class="label">
					<span class="text-sm font-bold text-slate-700 flex items-center gap-2">
						<Lock size={16} /> Password
					</span>
					<input 
						class="input" 
						type="password" 
						bind:value={password} 
						placeholder="••••••••" 
						required 
					/>
				</label>

				<div class="flex items-center justify-between pt-2">
					<label class="flex items-center gap-2 cursor-pointer">
						<input type="checkbox" class="checkbox" />
						<span class="text-sm text-slate-600">Remember me</span>
					</label>
					<a href="#/forgot" class="text-sm text-primary-600 font-medium hover:underline">Forgot password?</a>
				</div>

				<button 
					type="submit" 
					class="btn preset-filled-primary-500 w-full py-3 font-bold flex items-center justify-center gap-2 shadow-lg shadow-primary-500/30"
					disabled={isLoading}
				>
					{#if isLoading}
						<div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
						Authenticating...
					{:else}
						Sign In <ArrowRight size={18} />
					{/if}
				</button>
			</form>
		</div>

		<!-- Footer -->
		<p class="text-center text-sm text-slate-400">
			Protected by enterprise-grade security. <br/>
			<span class="font-mono text-xs opacity-50 uppercase tracking-tighter">GoSvelte v0.1.0-alpha</span>
		</p>
	</div>
</div>
