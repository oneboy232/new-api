/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useContext, useEffect } from 'react';
import { StatusContext } from '../../context/Status';
import { useTranslation } from 'react-i18next';
import { useActualTheme } from '../../context/Theme';
import { normalizeLanguage } from '../../i18n/language';

const Docs = () => {
  const { t, i18n } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const actualTheme = useActualTheme();
  const docsLink = statusState?.status?.docs_link || '';

  useEffect(() => {
    const iframe = document.querySelector('iframe[title="docs"]');
    const cw = iframe && iframe.contentWindow;
    if (cw) {
      cw.postMessage(
        { themeMode: actualTheme, lang: normalizeLanguage(i18n.language) },
        '*',
      );
    }
  }, [actualTheme, i18n.language]);

  if (!docsLink) {
    return (
      <div className='flex items-center justify-center h-[60vh] text-semi-color-text-2'>
        {t('管理员暂未设置文档链接')}
      </div>
    );
  }

  return (
    <iframe
      src={docsLink}
      style={{ width: '100%', height: '100vh', border: 'none' }}
      title='docs'
    />
  );
};

export default Docs;
