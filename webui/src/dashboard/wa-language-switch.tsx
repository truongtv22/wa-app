import { Button } from '@/components/ui/button';
import { useI18n } from '@/i18n/i18n';
import type { Locale } from '@/i18n/types';

const locales: Array<{ value: Locale; shortLabel: string; textKey: string; defaultText: string }> = [
  { value: 'zh-CN', shortLabel: '中', textKey: 'app.language_zh', defaultText: '中文' },
  { value: 'en', shortLabel: 'EN', textKey: 'app.language_en', defaultText: 'English' },
  { value: 'vi', shortLabel: 'VI', textKey: 'app.language_vi', defaultText: 'Tiếng Việt' },
];

export function WaLanguageSwitch() {
  const { locale, setLocale, t } = useI18n();
  return (
    <div className="inline-flex items-center gap-1 rounded-xl border border-border bg-card/90 p-1 shadow-sm backdrop-blur">
      <span className="px-2 text-xs text-muted-foreground">{t('app.language_switcher', '语言')}</span>
      {locales.map((item) => (
        <Button
          key={item.value}
          type="button"
          size="sm"
          variant={locale === item.value ? 'default' : 'ghost'}
          aria-label={t(item.textKey, item.defaultText)}
          title={t(item.textKey, item.defaultText)}
          onClick={() => setLocale(item.value)}
        >
          {item.shortLabel}
        </Button>
      ))}
    </div>
  );
}
