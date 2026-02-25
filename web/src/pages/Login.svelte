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
			<div class="inline-flex items-center justify-center w-16 h-16 bg-teal-600 text-white mb-4">
				<CloudUpload size={40} />
			</div>
			<h1 class="h1 font-black text-slate-900 uppercase tracking-tight">Sentral Finance</h1>
			<p class="text-slate-500 mt-2 font-medium">Sign in to your corporate financial dashboard</p>
		</div>

		<!-- Card -->
		<div class="card p-8 bg-white border border-slate-200">
			{#if error}
				<aside class="alert preset-tonal-error mb-6">
					<div class="alert-message">
						<p class="text-sm font-bold">{error}</p>
					</div>
				</aside>
			{/if}

			<form onsubmit={handleSubmit} class="space-y-6 text-left">
				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
						<Mail size={14} /> Email Address
					</span>
					<input 
						class="input" 
						type="email" 
						bind:value={email} 
						placeholder="user@system.com" 
						required 
					/>
				</label>

				<label class="label">
					<span class="text-xs font-black text-slate-400 uppercase tracking-widest flex items-center gap-2">
						<Lock size={14} /> Password
					</span>
					<input 
						class="input" 
						type="password" 
						bind:value={password} 
						placeholder="••••••••" 
						required 
					/>
				</label>

				<div class="flex items-center justify-between pt-2 text-left">
					<label class="flex items-center gap-2 cursor-pointer group">
						<input type="checkbox" class="checkbox" />
						<span class="text-sm text-slate-600 font-bold group-hover:text-teal-600">Remember me</span>
					</label>
					<a href="#/forgot" class="text-sm text-teal-600 font-black hover:underline">Forgot?</a>
				</div>

				<button 
					type="submit" 
					class="btn bg-teal-600 text-white w-full py-4 font-black uppercase tracking-widest flex items-center justify-center gap-2 hover:opacity-90 transition-opacity"
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
		<p class="text-center text-xs text-slate-400 font-bold uppercase tracking-[0.2em] opacity-50">
			GoSvelte v0.1.0-alpha
		</p>
	</div>
</div>
