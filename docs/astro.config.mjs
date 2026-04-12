// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://junara.github.io',
	base: '/encfixture',
	integrations: [
		starlight({
			title: 'encfixture',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/junara/encfixture' }],
			defaultLocale: 'ja',
			locales: {
				ja: { label: '日本語', lang: 'ja' },
				en: { label: 'English', lang: 'en' },
			},
			sidebar: [
				{
					label: 'はじめに',
					translations: { en: 'Getting Started' },
					items: [
						{ label: 'インストール', slug: 'getting-started/installation', translations: { en: 'Installation' } },
						{ label: 'クイックスタート', slug: 'getting-started/quickstart', translations: { en: 'Quick Start' } },
					],
				},
				{
					label: '使い方',
					translations: { en: 'Usage' },
					items: [
						{ label: '画像の生成', slug: 'usage/image', translations: { en: 'Image' } },
						{ label: '動画の生成', slug: 'usage/video', translations: { en: 'Video' } },
						{ label: '音声の生成', slug: 'usage/audio', translations: { en: 'Audio' } },
						{ label: 'オーバーレイ', slug: 'usage/overlay', translations: { en: 'Overlay' } },
					],
				},
				{
					label: 'リファレンス',
					translations: { en: 'Reference' },
					items: [
						{ label: '対応色', slug: 'reference/colors', translations: { en: 'Colors' } },
						{ label: 'JSON 出力', slug: 'reference/json', translations: { en: 'JSON Output' } },
					],
				},
				{
					label: '活用',
					translations: { en: 'Integration' },
					items: [
						{ label: 'Claude Code', slug: 'integration/claude-code' },
					],
				},
			],
		}),
	],
});
