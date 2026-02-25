console.log("[Main] Entry point reached");
import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

console.log("[Main] Mounting App...");
const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
