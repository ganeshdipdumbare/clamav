<script>
    import { darkMode } from '$lib/stores/theme';
    
    function toggleDarkMode() {
        darkMode.update(value => !value);
    }
</script>

<button class="dark-mode-toggle" on:click={toggleDarkMode}>
    {#if $darkMode}
        <i class="fas fa-sun"></i>
    {:else}
        <i class="fas fa-moon"></i>
    {/if}
</button> 