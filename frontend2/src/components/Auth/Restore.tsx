import { useCallback, useEffect, useState } from 'react';
import { useAppDispatch } from '../../hooks/hooks';
import { restore, sendCode } from './Auth.thunks';
import styles from './Register.module.css';
import { Spinner } from '../Spinner/Spinner';
import { useNavigate } from 'react-router-dom';
import { PasswordRepeat } from '../RepeatPassword/RepeatPassword';
import { useSelector } from 'react-redux';
import { RootState } from '../../store/store';
import { useTranslation } from 'react-i18next';

export const Restore: React.FC = () => {
     const { t } = useTranslation();
     const dispatch = useAppDispatch();
     const navigate = useNavigate();
     const [sending, setSending] = useState(false);
     const [form, setForm] = useState({
          email: '',
          code: '',
          password: '',
     });
     const isRegistered = useSelector((state: RootState) => state.auth.isRegistered);

     const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
          const { name, value } = e.target;
          setForm(prev => ({ ...prev, [name]: value }));
     };

     const handleSendCode = async () => {
          if (!form.email) return;
          setSending(true);
          await dispatch(sendCode(form.email));
          setSending(false);
     };

     useEffect(() => {
          if (isRegistered) {
               navigate("/login");
          }
     }, [isRegistered, navigate]);

     const handleSubmit = async (e: React.FormEvent) => {
          e.preventDefault();
          await dispatch(restore(form.email, form.code, form.password));
          navigate('/login');
     };

     const handlePasswordChange = useCallback((pwd: string) => {
          setForm(prev => ({ ...prev, password: pwd }));
     }, []);

     return (
          <div className={styles.container}>
               <h1 className={styles.title}>{t('auth.restore_title')}</h1>
               <form onSubmit={handleSubmit} className={styles.form}>
                    <div className={styles.row}>
                         <input
                              type="email"
                              name="email"
                              placeholder="Email"
                              value={form.email}
                              onChange={handleChange}
                              className={styles.input}
                              required
                         />
                         <button
                              type="button"
                              onClick={handleSendCode}
                              className={styles.codeButton}
                              disabled={!form.email || sending}
                         >
                              {sending ? <Spinner /> : t('auth.send_code')}
                         </button>
                    </div>
                    <input
                         type="text"
                         name="code"
                         placeholder={t('auth.confirmation_code')}
                         value={form.code}
                         onChange={handleChange}
                         className="w-full border p-2 rounded"
                         required
                    />
                    <PasswordRepeat onChange={handlePasswordChange} />
                    <button type="submit" className={styles.submit}>
                         {t('auth.change_password')}
                    </button>
               </form>
          </div>
     );
};
