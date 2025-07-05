import { useEffect } from 'react';
import i18n from 'i18next';

const SUPPORTED_LANGS = ['ru', 'en'] as const;
type SupportedLang = typeof SUPPORTED_LANGS[number];

function getBrowserLang(): SupportedLang | null {
     const browserLang = navigator.language;
     if (!browserLang) return null;
     const lang = browserLang.slice(0, 2);
     if (SUPPORTED_LANGS.includes(lang as SupportedLang)) {
          return lang as SupportedLang;
     }
     return null;
}

export function useAutoLanguage(defaultLang: SupportedLang = 'en') {
     useEffect(() => {
          const lang = getBrowserLang() || defaultLang;
          if (i18n.language !== lang) {
               i18n.changeLanguage(lang);
          }
     }, [defaultLang]);
}
