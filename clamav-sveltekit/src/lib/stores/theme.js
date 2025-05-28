import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Initialize the dark mode state from local storage or system preference
const initialValue = browser 
    ? localStorage.getItem('theme') === 'dark' || 
      (!localStorage.getItem('theme') && window.matchMedia('(prefers-color-scheme: dark)').matches)
    : false;

export const darkMode = writable(initialValue);

// Subscribe to changes and update localStorage and body class
if (browser) {
    darkMode.subscribe(value => {
        localStorage.setItem('theme', value ? 'dark' : 'light');
        if (value) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    });
} 