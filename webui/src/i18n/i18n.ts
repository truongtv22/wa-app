import { useSyncExternalStore } from 'react';
import { enMessages } from './messages-en';
import { viMessages } from './messages-vi';
import type { Locale, MessageCatalog } from './types';

const localeStorageKey = 'wa-app-locale';
const localeCookieName = 'wa_lang';
const catalogs: Record<Exclude<Locale, 'zh-CN'>, MessageCatalog> = {
  en: enMessages,
  vi: viMessages,
};

const listeners = new Set<() => void>();
let currentLocale: Locale = detectInitialLocale();

function detectInitialLocale(): Locale {
  const queryLocale = readQueryLocale();
  if (queryLocale) return queryLocale;
  const cookieLocale = normalizeLocale(readCookie(localeCookieName));
  if (cookieLocale) return cookieLocale;
  const storedLocale = normalizeLocale(typeof window !== 'undefined' ? window.localStorage.getItem(localeStorageKey) : '');
  if (storedLocale) return storedLocale;
  const browserLocale = normalizeLocale(typeof navigator !== 'undefined' ? navigator.language : '');
  return browserLocale || 'zh-CN';
}

function readQueryLocale() {
  if (typeof window === 'undefined') return null;
  return normalizeLocale(new URLSearchParams(window.location.search).get('lang'));
}

function readCookie(name: string) {
  if (typeof document === 'undefined') return '';
  const prefix = `${name}=`;
  return document.cookie
    .split(';')
    .map((part) => part.trim())
    .find((part) => part.startsWith(prefix))
    ?.slice(prefix.length) || '';
}

function normalizeLocale(value?: string | null): Locale | null {
  const raw = (value || '').trim().toLowerCase();
  if (!raw) return null;
  if (raw.startsWith('vi')) return 'vi';
  if (raw.startsWith('en')) return 'en';
  if (raw.startsWith('zh')) return 'zh-CN';
  return null;
}

function notify() {
  for (const listener of listeners) listener();
}

function applyLocale(locale: Locale) {
  if (typeof document !== 'undefined') {
    document.documentElement.lang = locale;
    document.cookie = `${localeCookieName}=${encodeURIComponent(locale)}; path=/; max-age=${60 * 60 * 24 * 365}; samesite=lax`;
    document.title = t('app.document_title', 'WA 收件箱');
  }
  if (typeof window !== 'undefined') window.localStorage.setItem(localeStorageKey, locale);
}

export function initI18n() {
  applyLocale(currentLocale);
}

export function getLocale() {
  return currentLocale;
}

export function setLocale(locale: Locale) {
  if (currentLocale === locale) return;
  currentLocale = locale;
  applyLocale(locale);
  notify();
}

export function setLocaleFromInput(value?: string | null) {
  const locale = normalizeLocale(value);
  if (locale) setLocale(locale);
}

export function translateTemplate(template: string, values: Record<string, string | number>) {
  return template.replace(/\{(\w+)\}/g, (_, key: string) => String(values[key] ?? `{${key}}`));
}

export function t(key: string, fallback: string, values?: Record<string, string | number>) {
  const locale = getLocale();
  const translated = locale === 'zh-CN' ? fallback : catalogs[locale]?.[key] || fallback;
  return values ? translateTemplate(translated, values) : translated;
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  return () => listeners.delete(listener);
}

export function useI18n() {
  const locale = useSyncExternalStore(subscribe, getLocale, getLocale);
  return { locale, setLocale, t };
}

export const i18n = { getLocale, setLocale, setLocaleFromInput, t };
