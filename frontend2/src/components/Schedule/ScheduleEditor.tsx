import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import styles from './ScheduleEditor.module.css';
import ScheduleByDate from './ScheduleByDate';
import ScheduleByWeekday from './ScheduleByWeekday';

const ScheduleEditor: React.FC = () => {
     const { t } = useTranslation();
     const [tab, setTab] = useState<'date' | 'weekday'>('date');

     return (
          <div className={styles.wrapper}>
               <h2 className={styles.heading}>{t('scheduleEditor.heading')}</h2>

               <div className={styles.tabs}>
                    <button
                         className={`${styles.tab} ${tab === 'date' ? styles.activeTab : ''}`}
                         onClick={() => setTab('date')}
                    >
                         {t('scheduleEditor.tabs.byDate')}
                    </button>
                    <button
                         className={`${styles.tab} ${tab === 'weekday' ? styles.activeTab : ''}`}
                         onClick={() => setTab('weekday')}
                    >
                         {t('scheduleEditor.tabs.byWeekday')}
                    </button>
               </div>

               {tab === 'date' && <ScheduleByDate />}
               {tab === 'weekday' && <ScheduleByWeekday />}
          </div>
     );
};

export default ScheduleEditor;
