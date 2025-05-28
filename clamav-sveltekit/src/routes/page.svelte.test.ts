import { describe, test, expect } from 'vitest';
import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/svelte';
import Page from './+page.svelte';

describe('/+page.svelte', () => {
	test('should render h1', () => {
		render(Page);
		expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
	});

	test('should render file upload form', () => {
		render(Page);
		expect(screen.getByLabelText(/upload file/i)).toBeInTheDocument();
	});

	test('should render text scanning area', () => {
		render(Page);
		expect(screen.getByLabelText(/text to scan/i)).toBeInTheDocument();
	});
});
