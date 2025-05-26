<script lang="ts">
    import { onMount } from 'svelte';
    import { apiUrl } from '$lib/stores/api';
    
    interface ScanResult {
        clean: boolean;
        status: string;
        filename: string;
    }
    
    interface InfoResult {
        type: 'success' | 'error';
        title: string;
        message: string;
    }
    
    let darkMode = false;
    let selectedFile: File | null = null;
    let fileName = "No file selected";
    let fileSize = "";
    let scanText = "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*";
    
    // Result states
    let infoResults: InfoResult | null = null;
    let textResults: ScanResult[] | null = null;
    let fileResults: ScanResult[] | null = null;
    
    // Loading states
    let pingLoading = false;
    let versionLoading = false;
    let textScanLoading = false;
    let fileScanLoading = false;
    
    // Active tab
    let activeTab = 'server';
    
    // Get the API URL
    let baseApiUrl: string;
    apiUrl.subscribe(value => {
        baseApiUrl = value;
    });
    
    onMount(() => {
        // Check for saved theme preference or use preferred color scheme
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            darkMode = true;
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    });
    
    function toggleDarkMode() {
        darkMode = !darkMode;
        if (darkMode) {
            localStorage.setItem('theme', 'dark');
            document.documentElement.classList.add('dark');
        } else {
            localStorage.setItem('theme', 'light');
            document.documentElement.classList.remove('dark');
        }
    }
    
    function handleFileUpload(event: Event) {
        const target = event.target as HTMLInputElement;
        const file = target.files?.[0];
        if (file) {
            selectedFile = file;
            fileName = file.name;
            fileSize = formatFileSize(file.size);
        }
    }
    
    function handleDrop(event: DragEvent) {
        event.preventDefault();
        const file = event.dataTransfer?.files[0];
        if (file) {
            selectedFile = file;
            fileName = file.name;
            fileSize = formatFileSize(file.size);
        }
    }
    
    function formatFileSize(bytes: number) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
    
    // Helper to construct API URL
    function getApiEndpoint(path: string): string {
        return `${baseApiUrl}${path.startsWith('/') ? path : '/' + path}`;
    }
    
    async function pingServer() {
        pingLoading = true;
        infoResults = null;
        
        try {
            const response = await fetch(getApiEndpoint('/ping'));
            const data = await response.json();
            
            if (data.success) {
                infoResults = {
                    type: 'success',
                    title: 'Connected to ClamAV Server!',
                    message: data.message
                };
            } else {
                infoResults = {
                    type: 'error',
                    title: 'Connection failed',
                    message: data.message
                };
            }
        } catch (error: unknown) {
            infoResults = {
                type: 'error',
                title: 'Error',
                message: error instanceof Error ? error.message : 'Unknown error'
            };
        } finally {
            pingLoading = false;
        }
    }
    
    async function getVersion() {
        versionLoading = true;
        infoResults = null;
        
        try {
            const response = await fetch(getApiEndpoint('/version'));
            const data = await response.json();
            
            infoResults = {
                type: 'success',
                title: 'ClamAV Version',
                message: data.version
            };
        } catch (error: unknown) {
            infoResults = {
                type: 'error',
                title: 'Error',
                message: error instanceof Error ? error.message : 'Unknown error'
            };
        } finally {
            versionLoading = false;
        }
    }
    
    async function scanTextContent() {
        textScanLoading = true;
        textResults = null;
        
        try {
            const response = await fetch(getApiEndpoint('/scan/text'), {
                method: 'POST',
                body: scanText
            });
            
            const data = await response.json();
            textResults = data.results;
        } catch (error: unknown) {
            textResults = [{
                clean: false,
                status: 'Error: ' + (error instanceof Error ? error.message : 'Unknown error'),
                filename: 'text-scan'
            }];
        } finally {
            textScanLoading = false;
        }
    }
    
    async function scanFile() {
        if (!selectedFile) {
            alert('Please select a file to scan');
            return;
        }
        
        fileScanLoading = true;
        fileResults = null;
        
        try {
            const formData = new FormData();
            formData.append('file', selectedFile);
            
            const response = await fetch(getApiEndpoint('/scan/file'), {
                method: 'POST',
                body: formData
            });
            
            const data = await response.json();
            fileResults = data.results;
        } catch (error: unknown) {
            fileResults = [{
                clean: false,
                status: 'Error: ' + (error instanceof Error ? error.message : 'Unknown error'),
                filename: fileName
            }];
        } finally {
            fileScanLoading = false;
        }
    }
</script>

<svelte:head>
    <!-- Add the dark mode class to html tag when darkMode is active -->
    {#if darkMode}
        <script>
            document.documentElement.classList.add('dark');
        </script>
    {:else}
        <script>
            document.documentElement.classList.remove('dark');
        </script>
    {/if}
</svelte:head>

<div class="min-h-screen transition-colors duration-200 bg-gray-50 {darkMode ? 'dark bg-gray-900' : ''}">
    <button 
        class="fixed top-4 right-4 w-10 h-10 rounded-full flex items-center justify-center cursor-pointer z-50 text-lg
        {darkMode ? 'bg-gray-800 text-yellow-400 border border-gray-700' : 'bg-white text-gray-800 border border-gray-200 shadow-sm'}"
        on:click={toggleDarkMode}
        aria-label={darkMode ? 'Switch to light mode' : 'Switch to dark mode'}
    >
        {#if darkMode}
            <i class="fas fa-sun"></i>
        {:else}
            <i class="fas fa-moon"></i>
        {/if}
    </button>

    <div class="max-w-6xl mx-auto px-4 pt-6 pb-10 flex flex-col min-h-screen">
        <header class="text-center mb-8">
            <h1 class="text-2xl md:text-3xl font-bold mb-2 {darkMode ? 'text-white' : 'text-[#222]'}">
                <i class="fas fa-shield-virus mr-2 text-blue-600"></i>
                ClamAV Virus Scanner
            </h1>
            <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                Secure file and text scanning for your protection
            </p>
        </header>
        
        <!-- Tab Navigation -->
        <div class="flex justify-center mb-8 border-b {darkMode ? 'border-gray-700' : 'border-gray-200'} overflow-x-auto">
            <button 
                class={"px-6 py-3 text-sm font-medium transition-colors " + (activeTab === 'server' 
                    ? (darkMode ? "text-blue-400 border-b-2 border-blue-400" : "text-blue-600 border-b-2 border-blue-600") 
                    : (darkMode ? "text-gray-400 hover:text-gray-300" : "text-gray-500 hover:text-gray-700"))}
                on:click={() => activeTab = 'server'}
            >
                <i class="fas fa-server mr-2"></i> Server Status
            </button>
            <button 
                class={"px-6 py-3 text-sm font-medium transition-colors " + (activeTab === 'text' 
                    ? (darkMode ? "text-blue-400 border-b-2 border-blue-400" : "text-blue-600 border-b-2 border-blue-600") 
                    : (darkMode ? "text-gray-400 hover:text-gray-300" : "text-gray-500 hover:text-gray-700"))}
                on:click={() => activeTab = 'text'}
            >
                <i class="fas fa-file-alt mr-2"></i> Text Scan
            </button>
            <button 
                class={"px-6 py-3 text-sm font-medium transition-colors " + (activeTab === 'file' 
                    ? (darkMode ? "text-blue-400 border-b-2 border-blue-400" : "text-blue-600 border-b-2 border-blue-600") 
                    : (darkMode ? "text-gray-400 hover:text-gray-300" : "text-gray-500 hover:text-gray-700"))}
                on:click={() => activeTab = 'file'}
            >
                <i class="fas fa-file-upload mr-2"></i> File Scan
            </button>
        </div>
        
        <div class="flex-grow overflow-auto px-4 md:px-12">
            <!-- Server Status Tab -->
            {#if activeTab === 'server'}
                <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-8 border {darkMode ? 'border-gray-700' : 'border-gray-100'}">
                    <div class="flex flex-col md:flex-row md:items-center mb-6">
                        <div class="flex items-center justify-center w-12 h-12 rounded-full 
                        bg-blue-50 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400 mr-4 text-xl">
                            <i class="fas fa-server"></i>
                        </div>
                        <div>
                            <h2 class="text-xl font-medium {darkMode ? 'text-white' : 'text-[#222]'}">Server Status</h2>
                            <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                                Check the connection to the ClamAV daemon and get its version
                            </p>
                        </div>
                    </div>
                    
                    <div class="flex flex-wrap gap-4 mb-6">
                        <button 
                            class="inline-flex items-center justify-center px-6 py-3 rounded-lg text-sm font-medium 
                            bg-blue-600 text-white hover:bg-blue-700 transition-colors
                            disabled:opacity-50 disabled:cursor-not-allowed"
                            on:click={pingServer}
                            disabled={pingLoading}
                        >
                            {#if pingLoading}
                                <i class="fas fa-spinner fa-spin mr-2"></i>
                            {:else}
                                <i class="fas fa-network-wired mr-2"></i>
                            {/if}
                            Ping Server
                        </button>
                        
                        <button 
                            class="inline-flex items-center justify-center px-6 py-3 rounded-lg text-sm font-medium 
                            {darkMode ? 'bg-gray-700 text-white hover:bg-gray-600' : 'bg-gray-100 text-gray-800 hover:bg-gray-200'} transition-colors
                            disabled:opacity-50 disabled:cursor-not-allowed"
                            on:click={getVersion}
                            disabled={versionLoading}
                        >
                            {#if versionLoading}
                                <i class="fas fa-spinner fa-spin mr-2"></i>
                            {:else}
                                <i class="fas fa-code-branch mr-2"></i>
                            {/if}
                            Get Version
                        </button>
                    </div>
                    
                    {#if infoResults}
                        <div class="p-5 rounded-lg {infoResults.type === 'success' 
                            ? (darkMode ? 'bg-green-900/20 text-green-200 border border-green-700/30' : 'bg-green-50 text-green-800 border border-green-100') 
                            : (darkMode ? 'bg-red-900/20 text-red-200 border border-red-700/30' : 'bg-red-50 text-red-800 border border-red-100')
                        } flex items-start">
                            <i class="fas {infoResults.type === 'success' ? 'fa-check-circle' : 'fa-exclamation-triangle'} mr-4 text-lg mt-0.5"></i>
                            <div>
                                <div class="font-medium text-base">{infoResults.title}</div>
                                <div class="mt-1">{infoResults.message}</div>
                            </div>
                        </div>
                    {/if}
                </div>
            {/if}
            
            <!-- Text Scan Tab -->
            {#if activeTab === 'text'}
                <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-8 border {darkMode ? 'border-gray-700' : 'border-gray-100'}">
                    <div class="flex flex-col md:flex-row md:items-center mb-6">
                        <div class="flex items-center justify-center w-12 h-12 rounded-full 
                        bg-blue-50 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400 mr-4 text-xl">
                            <i class="fas fa-file-alt"></i>
                        </div>
                        <div>
                            <h2 class="text-xl font-medium {darkMode ? 'text-white' : 'text-[#222]'}">Scan Text</h2>
                            <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                                Enter text to scan for potential virus signatures
                            </p>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <textarea 
                            bind:value={scanText}
                            class="w-full min-h-[140px] p-4 border rounded-lg text-sm font-mono resize-y outline-none transition-colors
                            {darkMode 
                                ? 'bg-gray-700 border-gray-600 text-gray-100 focus:border-blue-500' 
                                : 'bg-white border-gray-200 text-gray-900 focus:border-blue-600'}"
                            placeholder="Enter text to scan..."
                        ></textarea>
                        <p class="mt-2 text-xs {darkMode ? 'text-gray-400' : 'text-gray-500'} italic">
                            <i class="fas fa-info-circle mr-1"></i>
                            The default text is the EICAR test virus string.
                        </p>
                    </div>
                    
                    <div>
                        <button 
                            class="inline-flex items-center justify-center px-6 py-3 rounded-lg text-sm font-medium 
                            bg-blue-600 text-white hover:bg-blue-700 transition-colors
                            disabled:opacity-50 disabled:cursor-not-allowed"
                            on:click={scanTextContent}
                            disabled={textScanLoading}
                        >
                            {#if textScanLoading}
                                <i class="fas fa-spinner fa-spin mr-2"></i>
                            {:else}
                                <i class="fas fa-search mr-2"></i>
                            {/if}
                            Scan Text
                        </button>
                    </div>
                    
                    {#if textResults}
                        <div class="mt-8 {darkMode ? 'bg-gray-700/30' : 'bg-gray-50'} rounded-lg p-6 border-l-4 border-blue-600">
                            <h3 class="text-lg font-medium mb-4 {darkMode ? 'text-white' : 'text-[#222]'}">Scan Results</h3>
                            {#each textResults as result}
                                <div class="flex items-center mb-4 p-4 {darkMode ? 'bg-gray-800' : 'bg-white'} rounded-lg shadow-sm">
                                    <div class="mr-4 text-xl {result.clean 
                                        ? (darkMode ? 'text-green-400' : 'text-green-500') 
                                        : (darkMode ? 'text-red-400' : 'text-red-600')}">
                                        <i class="fas {result.clean ? 'fa-check-circle' : 'fa-virus'}"></i>
                                    </div>
                                    <div>
                                        <div class="font-medium {darkMode ? 'text-white' : 'text-[#222]'}">{result.filename}</div>
                                        <div class="{darkMode ? 'text-gray-300' : 'text-gray-700'}">
                                            Status: <span class="font-medium {result.clean 
                                                ? (darkMode ? 'text-green-400' : 'text-green-600') 
                                                : (darkMode ? 'text-red-400' : 'text-red-600')}">{result.status}</span>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}
            
            <!-- File Scan Tab -->
            {#if activeTab === 'file'}
                <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-8 border {darkMode ? 'border-gray-700' : 'border-gray-100'}">
                    <div class="flex flex-col md:flex-row md:items-center mb-6">
                        <div class="flex items-center justify-center w-12 h-12 rounded-full 
                        bg-blue-50 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400 mr-4 text-xl">
                            <i class="fas fa-file-upload"></i>
                        </div>
                        <div>
                            <h2 class="text-xl font-medium {darkMode ? 'text-white' : 'text-[#222]'}">Scan File</h2>
                            <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                                Upload a file to check it for viruses
                            </p>
                        </div>
                    </div>
                    
                    <!-- File upload area -->
                    <div 
                        class="flex flex-col items-center justify-center border-2 border-dashed 
                        {darkMode ? 'border-gray-600 bg-gray-700/20' : 'border-gray-200 bg-gray-50'} 
                        rounded-lg p-10 text-center my-6 
                        cursor-pointer transition-all duration-300 
                        hover:border-blue-600"
                        on:dragover|preventDefault
                        on:dragleave|preventDefault
                        on:drop|preventDefault={handleDrop}
                        on:click={() => {
                            const fileInput = document.getElementById('file-input');
                            if (fileInput) fileInput.click();
                        }}
                        on:keydown={(e) => {
                            if (e.key === 'Enter' || e.key === ' ') {
                                e.preventDefault();
                                const fileInput = document.getElementById('file-input');
                                if (fileInput) fileInput.click();
                            }
                        }}
                        role="button"
                        tabindex="0"
                        aria-label="Upload a file by clicking or dropping here"
                    >
                        <i class="fas fa-cloud-upload-alt text-4xl {darkMode ? 'text-blue-400' : 'text-blue-600'} mb-4"></i>
                        <p class="text-lg font-medium {darkMode ? 'text-white' : 'text-[#222]'} mb-2">
                            Drag & drop file here
                        </p>
                        <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                            or click to browse your files
                        </p>
                        <input 
                            type="file" 
                            id="file-input" 
                            class="hidden" 
                            on:change={handleFileUpload}
                            aria-hidden="true"
                        />
                    </div>
                    
                    {#if selectedFile}
                        <div class="{darkMode ? 'bg-blue-900/20 border border-blue-800/30' : 'bg-blue-50 border border-blue-100'} p-4 rounded-lg my-6">
                            <div class="flex items-center">
                                <i class="fas fa-file text-blue-600 dark:text-blue-400 mr-4 text-lg"></i>
                                <div>
                                    <div class="font-medium {darkMode ? 'text-white' : 'text-[#222]'}">{fileName}</div>
                                    <div class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">{fileSize}</div>
                                </div>
                            </div>
                        </div>
                    {/if}
                    
                    <div>
                        <button 
                            class="inline-flex items-center justify-center px-6 py-3 rounded-lg text-sm font-medium 
                            bg-blue-600 text-white hover:bg-blue-700 transition-colors
                            disabled:opacity-50 disabled:cursor-not-allowed"
                            on:click={scanFile}
                            disabled={fileScanLoading || !selectedFile}
                        >
                            {#if fileScanLoading}
                                <i class="fas fa-spinner fa-spin mr-2"></i>
                            {:else}
                                <i class="fas fa-search mr-2"></i>
                            {/if}
                            Scan File
                        </button>
                    </div>
                    
                    {#if fileResults}
                        <div class="mt-8 {darkMode ? 'bg-gray-700/30' : 'bg-gray-50'} rounded-lg p-6 border-l-4 border-blue-600">
                            <h3 class="text-lg font-medium mb-4 {darkMode ? 'text-white' : 'text-[#222]'}">Scan Results for {fileName}</h3>
                            {#each fileResults as result}
                                <div class="flex items-center mb-4 p-4 {darkMode ? 'bg-gray-800' : 'bg-white'} rounded-lg shadow-sm">
                                    <div class="mr-4 text-xl {result.clean 
                                        ? (darkMode ? 'text-green-400' : 'text-green-500') 
                                        : (darkMode ? 'text-red-400' : 'text-red-600')}">
                                        <i class="fas {result.clean ? 'fa-check-circle' : 'fa-virus'}"></i>
                                    </div>
                                    <div>
                                        <div class="font-medium {darkMode ? 'text-white' : 'text-[#222]'}">{result.filename}</div>
                                        <div class="{darkMode ? 'text-gray-300' : 'text-gray-700'}">
                                            Status: <span class="font-medium {result.clean 
                                                ? (darkMode ? 'text-green-400' : 'text-green-600') 
                                                : (darkMode ? 'text-red-400' : 'text-red-600')}">{result.status}</span>
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            {/if}
        </div>
        
        <footer class="text-center mt-10 py-4 border-t {darkMode ? 'border-gray-800' : 'border-gray-100'}">
            <p class="text-sm {darkMode ? 'text-gray-400' : 'text-gray-600'}">
                <i class="fas fa-shield-virus {darkMode ? 'text-blue-400' : 'text-blue-600'} mr-2"></i>
                ClamAV Virus Scanner | Secure file scanning for your protection
            </p>
        </footer>
    </div>
</div>
