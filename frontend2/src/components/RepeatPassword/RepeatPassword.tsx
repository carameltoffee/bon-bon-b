import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

interface PasswordRepeatProps {
     onChange: (password: string) => void;
}

export const PasswordRepeat: React.FC<PasswordRepeatProps> = ({ onChange }) => {
     const { t } = useTranslation();
     const [password, setPassword] = useState('');
     const [repeat, setRepeat] = useState('');
     const [error, setError] = useState('');

     useEffect(() => {
          if (password && repeat && password !== repeat) {
               setError(t('passwordRepeat.errorMismatch'));
               onChange('');
          } else {
               setError('');
               if (password && repeat) {
                    onChange(password);
               }
          }
     }, [password, repeat, onChange, t]);

     return (
          <div style={{ width: '100%' }}>
               <input
                    type="password"
                    placeholder={t('passwordRepeat.placeholderPassword')}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    style={{ width: '100%', display: 'block', marginBottom: '0.5rem' }}
                    required
               />
               <input
                    type="password"
                    placeholder={t('passwordRepeat.placeholderRepeat')}
                    value={repeat}
                    onChange={(e) => setRepeat(e.target.value)}
                    style={{ width: '100%', display: 'block' }}
                    required
               />
               {error && <div style={{ color: 'red', marginTop: '0.5rem' }}>{error}</div>}
          </div>
     );
};
